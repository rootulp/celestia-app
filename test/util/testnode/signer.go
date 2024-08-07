package testnode

import (
	"github.com/celestiaorg/celestia-app/v3/app"
	"github.com/celestiaorg/celestia-app/v3/app/encoding"
	"github.com/celestiaorg/celestia-app/v3/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/v3/pkg/user"
	"github.com/celestiaorg/celestia-app/v3/test/util/testfactory"
)

func NewOfflineSigner() (*user.Signer, error) {
	encCfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	kr, _ := NewKeyring(testfactory.TestAccName)
	return user.NewSigner(kr, encCfg.TxConfig, testfactory.ChainID, appconsts.LatestVersion, user.NewAccount(testfactory.TestAccName, 0, 0))
}

func NewTxClientFromContext(ctx Context) (*user.TxClient, error) {
	encCfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	return user.SetupTxClient(ctx.GoContext(), ctx.Keyring, ctx.GRPCClient, encCfg)
}
