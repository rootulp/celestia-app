package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/celestiaorg/go-square/blob"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc"

	"github.com/celestiaorg/celestia-app/v2/app"
	"github.com/celestiaorg/celestia-app/v2/app/encoding"
	apperrors "github.com/celestiaorg/celestia-app/v2/app/errors"
	"github.com/celestiaorg/celestia-app/v2/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v2/x/blob/types"
	"github.com/celestiaorg/celestia-app/v2/x/minfee"
)

const (
	DefaultPollTime              = 3 * time.Second
	DefaultGasMultiplier float64 = 1.1
)

type Option func(client *TxClient)

// WithGasMultiplier is a functional option allows to configure the gas multiplier.
func WithGasMultiplier(multiplier float64) Option {
	return func(c *TxClient) {
		c.gasMultiplier = multiplier
	}
}

// WithDefaultGasPrice sets the gas price.
func WithDefaultGasPrice(price float64) Option {
	return func(c *TxClient) {
		c.defaultGasPrice = price
	}
}

func WithPollTime(time time.Duration) Option {
	return func(c *TxClient) {
		c.pollTime = time
	}
}

func WithDefaultAddress(address sdktypes.AccAddress) Option {
	return func(c *TxClient) {
		record, err := c.signer.keys.KeyByAddress(address)
		if err != nil {
			panic(err)
		}
		c.defaultAccount = record.Name
		c.defaultAddress = address
	}
}

func WithDefaultAccount(name string) Option {
	return func(c *TxClient) {
		rec, err := c.signer.keys.Key(name)
		if err != nil {
			panic(err)
		}
		addr, err := rec.GetAddress()
		if err != nil {
			panic(err)
		}
		c.defaultAccount = name
		c.defaultAddress = addr
	}
}

// TxClient is an abstraction for building, signing, and broadcasting Celestia transactions
// It supports multiple accounts. If none is specified, it will
// try use the default account.
// TxClient is thread-safe.
type TxClient struct {
	mtx      sync.Mutex
	signer   *Signer
	registry codectypes.InterfaceRegistry
	grpc     *grpc.ClientConn
	// how often to poll the network for confirmation of a transaction
	pollTime time.Duration
	// gasMultiplier is used to increase gas limit as it is sometimes underestimated
	gasMultiplier float64
	// defaultGasPrice is the price used if no price is provided
	defaultGasPrice float64
	defaultAccount  string
	defaultAddress  sdktypes.AccAddress
}

// NewTxClient returns a new signer using the provided keyring
func NewTxClient(
	signer *Signer,
	conn *grpc.ClientConn,
	registry codectypes.InterfaceRegistry,
	options ...Option,
) (*TxClient, error) {
	records, err := signer.keys.List()
	if err != nil {
		return nil, fmt.Errorf("retrieving keys: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("signer must have at least one key")
	}

	addr, err := records[0].GetAddress()
	if err != nil {
		return nil, err
	}

	txClient := &TxClient{
		mtx:             sync.Mutex{},
		signer:          signer,
		registry:        registry,
		grpc:            conn,
		pollTime:        DefaultPollTime,
		gasMultiplier:   DefaultGasMultiplier,
		defaultGasPrice: appconsts.DefaultMinGasPrice,
		defaultAccount:  records[0].Name,
		defaultAddress:  addr,
	}

	for _, opt := range options {
		opt(txClient)
	}

	return txClient, nil
}

// SetupTxClient uses the underlying grpc connection to populate the chainID, accountNumber and sequence number of all
// the accounts in the keyring.
func SetupTxClient(
	ctx context.Context,
	keys keyring.Keyring,
	conn *grpc.ClientConn,
	encCfg encoding.Config,
	options ...Option,
) (*TxClient, error) {
	resp, err := tmservice.NewServiceClient(conn).GetLatestBlock(
		ctx,
		&tmservice.GetLatestBlockRequest{},
	)
	if err != nil {
		return nil, err
	}

	chainID := resp.SdkBlock.Header.ChainID
	appVersion := resp.SdkBlock.Header.Version.App

	records, err := keys.List()
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, 0, len(records))
	for _, record := range records {
		addr, err := record.GetAddress()
		if err != nil {
			return nil, err
		}
		accNum, seqNum, err := QueryAccount(ctx, conn, encCfg.InterfaceRegistry, addr)
		if err != nil {
			// skip over the accounts that don't exist in state
			continue
		}

		accounts = append(accounts, NewAccount(record.Name, accNum, seqNum))
	}

	// query the min gas price from the chain
	minPrice, err := QueryMinimumGasPrice(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("querying minimum gas price: %w", err)
	}
	options = append([]Option{WithDefaultGasPrice(minPrice)}, options...)

	signer, err := NewSigner(keys, encCfg.TxConfig, chainID, appVersion, accounts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	return NewTxClient(signer, conn, encCfg.InterfaceRegistry, options...)
}

// SubmitPayForBlob forms a transaction from the provided blobs, signs it, and submits it to the chain.
// TxOptions may be provided to set the fee and gas limit.
func (client *TxClient) SubmitPayForBlob(ctx context.Context, blobs []*blob.Blob, opts ...TxOption) (*sdktypes.TxResponse, error) {
	resp, err := client.BroadcastPayForBlob(ctx, blobs, opts...)
	if err != nil {
		return resp, err
	}

	return client.ConfirmTx(ctx, resp.TxHash)
}

func (client *TxClient) SubmitPayForBlobWithAccount(ctx context.Context, account string, blobs []*blob.Blob, opts ...TxOption) (*sdktypes.TxResponse, error) {
	resp, err := client.BroadcastPayForBlobWithAccount(ctx, account, blobs, opts...)
	if err != nil {
		return resp, err
	}

	return client.ConfirmTx(ctx, resp.TxHash)
}

// BroadcastPayForBlob signs and broadcasts a transaction to pay for blobs.
// It does not confirm that the transaction has been committed on chain.
// If no gas or gas price is set, it will estimate the gas and use
// the max effective gas price: max(localMinGasPrice, networkMinGasPrice).
func (client *TxClient) BroadcastPayForBlob(ctx context.Context, blobs []*blob.Blob, opts ...TxOption) (*sdktypes.TxResponse, error) {
	return client.BroadcastPayForBlobWithAccount(ctx, client.defaultAccount, blobs, opts...)
}

func (client *TxClient) BroadcastPayForBlobWithAccount(ctx context.Context, account string, blobs []*blob.Blob, opts ...TxOption) (*sdktypes.TxResponse, error) {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	if err := client.checkAccountLoaded(ctx, account); err != nil {
		return nil, err
	}

	blobSizes := make([]uint32, len(blobs))
	for i, blob := range blobs {
		blobSizes[i] = uint32(len(blob.Data))
	}

	gasLimit := uint64(float64(types.DefaultEstimateGas(blobSizes)) * client.gasMultiplier)
	fee := uint64(math.Ceil(appconsts.DefaultMinGasPrice * float64(gasLimit)))
	// prepend calculated params, so it can be overwritten in case the user has specified it.
	opts = append([]TxOption{SetGasLimit(gasLimit), SetFee(fee)}, opts...)

	txBytes, _, err := client.signer.CreatePayForBlobs(account, blobs, opts...)
	if err != nil {
		return nil, err
	}

	return client.broadcastTx(ctx, txBytes, account)
}

// SubmitTx forms a transaction from the provided messages, signs it, and submits it to the chain. TxOptions
// may be provided to set the fee and gas limit.
func (client *TxClient) SubmitTx(ctx context.Context, msgs []sdktypes.Msg, opts ...TxOption) (*sdktypes.TxResponse, error) {
	resp, err := client.BroadcastTx(ctx, msgs, opts...)
	if err != nil {
		return resp, err
	}

	return client.ConfirmTx(ctx, resp.TxHash)
}

func (client *TxClient) BroadcastTx(ctx context.Context, msgs []sdktypes.Msg, opts ...TxOption) (*sdktypes.TxResponse, error) {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	account, err := client.getAccountNameFromMsgs(msgs)
	if err != nil {
		return nil, err
	}

	if err := client.checkAccountLoaded(ctx, account); err != nil {
		return nil, err
	}

	txBuilder, err := client.signer.txBuilder(msgs, opts...)
	if err != nil {
		return nil, err
	}

	hasUserSetFee := false
	for _, coin := range txBuilder.GetTx().GetFee() {
		if coin.Denom == appconsts.BondDenom {
			hasUserSetFee = true
			break
		}
	}

	gasLimit := txBuilder.GetTx().GetGas()
	if gasLimit == 0 {
		if !hasUserSetFee {
			// add at least 1utia as fee to builder as it affects gas calculation.
			txBuilder.SetFeeAmount(sdktypes.NewCoins(sdktypes.NewCoin(appconsts.BondDenom, sdktypes.NewInt(1))))
		}
		gasLimit, err = client.estimateGas(ctx, txBuilder)
		if err != nil {
			return nil, err
		}
		txBuilder.SetGasLimit(gasLimit)
	}

	if !hasUserSetFee {
		fee := int64(math.Ceil(appconsts.DefaultMinGasPrice * float64(gasLimit)))
		txBuilder.SetFeeAmount(sdktypes.NewCoins(sdktypes.NewCoin(appconsts.BondDenom, sdktypes.NewInt(fee))))
	}

	account, _, err = client.signer.signTransaction(txBuilder)
	if err != nil {
		return nil, err
	}

	txBytes, err := client.signer.EncodeTx(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return client.broadcastTx(ctx, txBytes, account)
}

func (client *TxClient) broadcastTx(ctx context.Context, txBytes []byte, signer string) (*sdktypes.TxResponse, error) {
	txClient := sdktx.NewServiceClient(client.grpc)
	resp, err := txClient.BroadcastTx(
		ctx,
		&sdktx.BroadcastTxRequest{
			Mode:    sdktx.BroadcastMode_BROADCAST_MODE_SYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.TxResponse.Code != abci.CodeTypeOK {
		if apperrors.IsNonceMismatchCode(resp.TxResponse.Code) {
			// query the account to update the sequence number on-chain for the account
			_, seqNum, err := QueryAccount(ctx, client.grpc, client.registry, client.signer.accounts[signer].address)
			if err != nil {
				return nil, fmt.Errorf("querying account for new sequence number: %w\noriginal tx response: %s", err, resp.TxResponse.RawLog)
			}
			if err := client.signer.SetSequence(signer, seqNum); err != nil {
				return nil, fmt.Errorf("setting sequence: %w", err)
			}
			return client.retryBroadcastingTx(ctx, txBytes)
		}
		return resp.TxResponse, fmt.Errorf("tx failed with code %d: %s", resp.TxResponse.Code, resp.TxResponse.RawLog)
	}

	// after the transaction has been submitted, we can increment the
	// sequence of the signer
	if err := client.signer.IncrementSequence(signer); err != nil {
		return nil, fmt.Errorf("increment sequencing: %w", err)
	}
	return resp.TxResponse, nil
}

// retryBroadcastingTx creates a new transaction by copying over an existing transaction but creates a new signature with the
// new sequence number. It then calls `broadcastTx` and attempts to submit the transaction
func (client *TxClient) retryBroadcastingTx(ctx context.Context, txBytes []byte) (*sdktypes.TxResponse, error) {
	blobTx, isBlobTx := blob.UnmarshalBlobTx(txBytes)
	if isBlobTx {
		txBytes = blobTx.Tx
	}
	tx, err := client.signer.DecodeTx(txBytes)
	if err != nil {
		return nil, err
	}

	opts := make([]TxOption, 0)
	if granter := tx.FeeGranter(); granter != nil {
		opts = append(opts, SetFeeGranter(granter))
	}
	if payer := tx.FeePayer(); payer != nil {
		opts = append(opts, SetFeePayer(payer))
	}
	if memo := tx.GetMemo(); memo != "" {
		opts = append(opts, SetMemo(memo))
	}
	if fee := tx.GetFee(); fee != nil {
		opts = append(opts, SetFee(fee.AmountOf(appconsts.BondDenom).Uint64()))
	}
	if gas := tx.GetGas(); gas > 0 {
		opts = append(opts, SetGasLimit(gas))
	}

	txBuilder, err := client.signer.txBuilder(tx.GetMsgs(), opts...)
	if err != nil {
		return nil, err
	}
	signer, _, err := client.signer.signTransaction(txBuilder)
	if err != nil {
		return nil, fmt.Errorf("resigning transaction: %w", err)
	}

	newTxBytes, err := client.signer.EncodeTx(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	// rewrap the blob tx if it was originally a blob tx
	if isBlobTx {
		newTxBytes, err = blob.MarshalBlobTx(newTxBytes, blobTx.Blobs...)
		if err != nil {
			return nil, err
		}
	}

	return client.broadcastTx(ctx, newTxBytes, signer)
}

// ConfirmTx periodically pings the provided node for the commitment of a transaction by its
// hash. It will continually loop until the context is cancelled, the tx is found or an error
// is encountered.
func (client *TxClient) ConfirmTx(ctx context.Context, txHash string) (*sdktypes.TxResponse, error) {
	txClient := sdktx.NewServiceClient(client.grpc)

	pollTicker := time.NewTicker(client.pollTime)
	defer pollTicker.Stop()

	for {
		resp, err := txClient.GetTx(ctx, &sdktx.GetTxRequest{Hash: txHash})
		if err == nil {
			if resp.TxResponse.Code != 0 {
				return resp.TxResponse, fmt.Errorf("tx was included but failed with code %d: %s", resp.TxResponse.Code, resp.TxResponse.RawLog)
			}
			return resp.TxResponse, nil
		}
		// FIXME: this is a relatively brittle of working out whether to retry or not. The tx might be not found for other
		// reasons. It may have been removed from the mempool at a later point. We should build an endpoint that gives the
		// signer more information on the status of their transaction and then update the logic here
		if !strings.Contains(err.Error(), "not found") {
			return &sdktypes.TxResponse{}, err
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return &sdktypes.TxResponse{}, ctx.Err()
		case <-pollTicker.C:
		}
	}
}

// EstimateGas simulates the transaction, calculating the amount of gas that was consumed during execution. The final
// result will be multiplied by gasMultiplier(that is set in TxClient)
func (client *TxClient) EstimateGas(ctx context.Context, msgs []sdktypes.Msg, opts ...TxOption) (uint64, error) {
	client.mtx.Lock()
	defer client.mtx.Unlock()

	txBuilder, err := client.signer.txBuilder(msgs, opts...)
	if err != nil {
		return 0, err
	}

	return client.estimateGas(ctx, txBuilder)
}

func (client *TxClient) estimateGas(ctx context.Context, txBuilder client.TxBuilder) (uint64, error) {
	// add at least 1utia as fee to builder as it affects gas calculation.
	txBuilder.SetFeeAmount(sdktypes.NewCoins(sdktypes.NewCoin(appconsts.BondDenom, sdktypes.NewInt(1))))

	_, _, err := client.signer.signTransaction(txBuilder)
	if err != nil {
		return 0, err
	}
	txBytes, err := client.signer.EncodeTx(txBuilder.GetTx())
	if err != nil {
		return 0, err
	}
	resp, err := sdktx.NewServiceClient(client.grpc).Simulate(ctx, &sdktx.SimulateRequest{
		TxBytes: txBytes,
	})
	if err != nil {
		return 0, err
	}

	gasLimit := uint64(float64(resp.GasInfo.GasUsed) * client.gasMultiplier)
	return gasLimit, nil
}

// Account returns an account of the signer from the key name. Also returns a bool if the
// account exists.
// Thread-safe
func (client *TxClient) Account(name string) *Account {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	acc, exists := client.signer.accounts[name]
	if !exists {
		return nil
	}
	return acc.Copy()
}

func (client *TxClient) AccountByAddress(address sdktypes.AccAddress) *Account {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	return client.signer.AccountByAddress(address)
}

func (client *TxClient) DefaultAddress() sdktypes.AccAddress {
	return client.defaultAddress
}

func (client *TxClient) DefaultAccountName() string { return client.defaultAccount }

func (client *TxClient) checkAccountLoaded(ctx context.Context, account string) error {
	if _, exists := client.signer.accounts[account]; exists {
		return nil
	}
	record, err := client.signer.keys.Key(account)
	if err != nil {
		return fmt.Errorf("trying to find account %s on keyring: %w", account, err)
	}
	addr, err := record.GetAddress()
	if err != nil {
		return fmt.Errorf("retrieving address from keyring: %w", err)
	}
	accNum, sequence, err := QueryAccount(ctx, client.grpc, client.registry, addr)
	if err != nil {
		return fmt.Errorf("querying account %s: %w", account, err)
	}
	return client.signer.AddAccount(NewAccount(account, accNum, sequence))
}

func (client *TxClient) getAccountNameFromMsgs(msgs []sdktypes.Msg) (string, error) {
	var addr sdktypes.AccAddress
	for _, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return "", fmt.Errorf("only one signer per transaction supported, got %d", len(signers))
		}
		if addr == nil {
			addr = signers[0]
		}
		if !bytes.Equal(addr, signers[0]) {
			return "", errors.New("not supported: got two different signers across multiple messages")
		}
	}
	record, err := client.signer.keys.KeyByAddress(addr)
	if err != nil {
		return "", err
	}
	return record.Name, nil
}

// Signer exposes the tx clients underlying signer
func (client *TxClient) Signer() *Signer {
	return client.signer
}

func (client *TxClient) SetDefaultGasPrice(price float64) {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	client.defaultGasPrice = price
}

func (client *TxClient) SetGasMultiplier(multiplier float64) {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	client.gasMultiplier = multiplier
}

// QueryMinimumGasPrice queries both the nodes local and network wide
// minimum gas prices, returning the maximum of the two.
func QueryMinimumGasPrice(ctx context.Context, grpcConn *grpc.ClientConn) (float64, error) {
	cfgRsp, err := nodeservice.NewServiceClient(grpcConn).Config(ctx, &nodeservice.ConfigRequest{})
	if err != nil {
		return 0, err
	}

	localMinCoins, err := sdktypes.ParseDecCoins(cfgRsp.MinimumGasPrice)
	if err != nil {
		return 0, err
	}
	localMinPrice := localMinCoins.AmountOf(app.BondDenom).MustFloat64()

	networkMinPrice, err := QueryNetworkMinGasPrice(ctx, grpcConn)
	if err != nil {
		// check if the network version supports a global min gas
		// price using a regex check. If not (i.e. v1) use the
		// local price only
		if strings.Contains(err.Error(), "unknown subspace: minfee") {
			return localMinPrice, nil
		}
		return 0, err
	}

	// return the highest value of the two
	if networkMinPrice > localMinPrice {
		return networkMinPrice, nil
	}
	return localMinPrice, nil
}

func QueryNetworkMinGasPrice(ctx context.Context, grpcConn *grpc.ClientConn) (float64, error) {
	paramsClient := paramtypes.NewQueryClient(grpcConn)
	// NOTE: that we don't prove that this is the correct value
	paramResponse, err := paramsClient.Params(ctx, &paramtypes.QueryParamsRequest{Subspace: minfee.ModuleName, Key: string(minfee.KeyNetworkMinGasPrice)})
	if err != nil {
		return 0, fmt.Errorf("querying params module: %w", err)
	}

	var networkMinPrice float64
	// Value is empty if network min gas price is not supported i.e. v1 state machine.
	if paramResponse.Param.Value != "" {
		networkMinPrice, err = strconv.ParseFloat(strings.Trim(paramResponse.Param.Value, `"`), 64)
		if err != nil {
			return 0, fmt.Errorf("parsing network min gas price: %w", err)
		}
	}
	return networkMinPrice, nil
}
