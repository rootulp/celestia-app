package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/test/txsim"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// A set of environment variables that can be used instead of flags
const (
	TxsimGRPC     = "TXSIM_GRPC"
	TxsimRPC      = "TXSIM_RPC"
	TxsimSeed     = "TXSIM_SEED"
	TxsimPoll     = "TXSIM_POLL"
	TxsimKeypath  = "TXSIM_KEYPATH"
	TxsimMnemonic = "TXSIM_MNEMONIC"
)

// Values for all flags
var (
	keyPath, keyMnemonic, rpcEndpoints, grpcEndpoints string
	blobSizes, blobAmounts                            string
	seed                                              int64
	pollTime                                          time.Duration
	send, sendIterations, sendAmount                  int
	stake, stakeValue, blob                           int
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	if err := command().ExecuteContext(ctx); err != nil {
		fmt.Print(err)
	}
}

// command returns the cobra command which wraps the txsim.Run() function using flags and/or
// environment variables to instruct the client.
func command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txsim",
		Short: "Celestia Transaction Simulator",
		Long: `
Txsim is a tool for randomized transaction generation on celestia networks. The tool relies on 
defined sequences; recursive patterns between one or more accounts which will continually submit 
transactions. You can use flags or environment variables (TXSIM_RPC, TXSIM_GRPC, TXSIM_SEED, 
TXSIM_POLL, TXSIM_KEYPATH) to configure the client. The keyring provided should have at least one
well funded account that can act as the master account. The command runs until all sequences error.`,
		Example: "txsim --key-path /path/to/keyring --rpc-endpoints localhost:26657 --grpc-endpoints localhost:9090 --seed 1234 --poll-time 1s --blob 5",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				keys keyring.Keyring
				err  error
				cdc  = encoding.MakeConfig(app.ModuleEncodingRegisters...).Codec
			)

			// setup the keyring
			switch {
			case keyPath != "":
				keys, err = keyring.New(app.Name, keyring.BackendTest, keyPath, nil, cdc)
			case keyPath == "" && keyMnemonic != "":
				keys = keyring.NewInMemory(cdc)
				_, err = keys.NewAccount("master", keyMnemonic, keyring.DefaultBIP39Passphrase, "", hd.Secp256k1)
			case os.Getenv(TxsimKeypath) != "":
				keys, err = keyring.New(app.Name, keyring.BackendTest, os.Getenv(TxsimKeypath), nil, cdc)
			case os.Getenv(TxsimMnemonic) != "":
				keys = keyring.NewInMemory(cdc)
				_, err = keys.NewAccount("master", os.Getenv(TxsimMnemonic), keyring.DefaultBIP39Passphrase, "", hd.Secp256k1)
			default:
				return errors.New("keyring not specified. Use --key-path, --key-mnemonic or TXSIM_KEYPATH env var")
			}
			if err != nil {
				return err
			}

			// get the rpc and grpc endpoints
			if rpcEndpoints == "" {
				rpcEndpoints = os.Getenv(TxsimRPC)
				if rpcEndpoints == "" {
					return errors.New("rpc endpoints not specified. Use --rpc-endpoints or TXSIM_RPC env var")
				}
			}
			if grpcEndpoints == "" {
				grpcEndpoints = os.Getenv(TxsimGRPC)
				if grpcEndpoints == "" {
					return errors.New("grpc endpoints not specified. Use --grpc-endpoints or TXSIM_GRPC env var")
				}
			}

			if stake == 0 && send == 0 && blob == 0 {
				return errors.New("no sequences specified. Use --stake, --send or --blob")
			}

			// setup the sequences
			sequences := []txsim.Sequence{}

			if stake > 0 {
				sequences = append(sequences, txsim.NewStakeSequence(stakeValue).Clone(stake)...)
			}

			if send > 0 {
				sequences = append(sequences, txsim.NewSendSequence(2, sendAmount, sendIterations).Clone(send)...)
			}

			if blob > 0 {
				sizes, err := readRange(blobSizes)
				if err != nil {
					return fmt.Errorf("invalid blob sizes: %w", err)
				}

				blobsPerPFB, err := readRange(blobAmounts)
				if err != nil {
					return fmt.Errorf("invalid blob amounts: %w", err)
				}

				sequences = append(sequences, txsim.NewBlobSequence(sizes, blobsPerPFB).Clone(blob)...)
			}

			if seed == 0 {
				if os.Getenv(TxsimSeed) != "" {
					seed, err = strconv.ParseInt(os.Getenv(TxsimSeed), 10, 64)
					if err != nil {
						return fmt.Errorf("parsing seed: %w", err)
					}
				} else {
					// use a random seed if none is set
					seed = rand.Int63()
				}
			}

			if os.Getenv(TxsimPoll) != "" && pollTime != txsim.DefaultPollTime {
				pollTime, err = time.ParseDuration(os.Getenv(TxsimPoll))
				if err != nil {
					return fmt.Errorf("parsing poll time: %w", err)
				}
			}

			err = txsim.Run(
				cmd.Context(),
				strings.Split(rpcEndpoints, ","),
				strings.Split(grpcEndpoints, ","),
				keys,
				seed,
				pollTime,
				sequences...,
			)
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		},
	}
	cmd.Flags().AddFlagSet(flags())

	return cmd
}

func flags() *flag.FlagSet {
	flags := &flag.FlagSet{}
	flags.StringVar(&keyPath, "key-path", "", "path to the keyring")
	flags.StringVar(&keyMnemonic, "key-mnemonic", "", "space separated mnemonic for the keyring. The hdpath used is an empty string")
	flags.StringVar(&rpcEndpoints, "rpc-endpoints", "", "comma separated list of rpc endpoints")
	flags.StringVar(&grpcEndpoints, "grpc-endpoints", "", "comma separated list of grpc endpoints")
	flags.Int64Var(&seed, "seed", 0, "seed for the random number generator")
	flags.DurationVar(&pollTime, "poll-time", txsim.DefaultPollTime, "poll time for the transaction client")
	flags.IntVar(&send, "send", 0, "number of send sequences to run")
	flags.IntVar(&sendIterations, "send-iterations", 1000, "number of send iterations to run per sequence")
	flags.IntVar(&sendAmount, "send-amount", 1000, "amount to send from one account to another")
	flags.IntVar(&stake, "stake", 0, "number of stake sequences to run")
	flags.IntVar(&stakeValue, "stake-value", 1000, "amount of initial stake per sequence")
	flags.IntVar(&blob, "blob", 0, "number of blob sequences to run")
	flags.StringVar(&blobSizes, "blob-sizes", "100-1000", "range of blob sizes to send")
	flags.StringVar(&blobAmounts, "blob-amounts", "1", "range of blobs to send per PFB in a sequence")
	return flags
}

// readRange takes a string expected to be of the form "1-10" and returns the corresponding Range.
// If only one number is set i.e. "5", the range returned is {5, 5}.
func readRange(r string) (txsim.Range, error) {
	if r == "" {
		return txsim.Range{}, errors.New("range is empty")
	}

	res := strings.Split(r, "-")
	n, err := strconv.Atoi(res[0])
	if err != nil {
		return txsim.Range{}, err
	}
	if len(res) == 1 {
		return txsim.NewRange(n, n), nil
	}
	m, err := strconv.Atoi(res[1])
	if err != nil {
		return txsim.Range{}, err
	}

	return txsim.NewRange(n, m), nil
}
