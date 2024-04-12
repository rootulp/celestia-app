package module_test

import (
	"testing"

	"github.com/celestiaorg/celestia-app/v2/app"
	"github.com/celestiaorg/celestia-app/v2/app/encoding"
	"github.com/celestiaorg/celestia-app/v2/app/module"
	"github.com/celestiaorg/celestia-app/v2/x/signal"
	signaltypes "github.com/celestiaorg/celestia-app/v2/x/signal/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/tests/mocks"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func TestModuleManagerMigration(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockAppModule1 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule2 := mocks.NewMockAppModule(mockCtrl)
	mockAppModule3 := mocks.NewMockAppModule(mockCtrl)

	mockAppModule1.EXPECT().Name().Return("testModule").AnyTimes()
	mockAppModule2.EXPECT().Name().Return("testModule").AnyTimes()
	mockAppModule3.EXPECT().Name().Return("differentModule").AnyTimes()
	mockAppModule1.EXPECT().ConsensusVersion().Return(uint64(1)).AnyTimes()
	mockAppModule2.EXPECT().ConsensusVersion().Return(uint64(2)).AnyTimes()
	mockAppModule3.EXPECT().ConsensusVersion().Return(uint64(5)).AnyTimes()
	mockAppModule3.EXPECT().InitGenesis(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)
	mockAppModule3.EXPECT().DefaultGenesis(gomock.Any()).Return(nil)

	mm, err := module.NewManager([]module.VersionedModule{
		// this is an existing module that gets updated in v2
		{Module: mockAppModule1, FromVersion: 1, ToVersion: 1},
		{Module: mockAppModule2, FromVersion: 2, ToVersion: 3},
		// This is a new module that gets added in v2
		{Module: mockAppModule3, FromVersion: 2, ToVersion: 2},
	})
	require.NoError(t, err)
	require.NotNil(t, mm)

	mockServer := mocks.NewMockServer(mockCtrl)
	cdc := encoding.MakeConfig(app.ModuleEncodingRegisters...)

	isCalled := false
	cfg := module.NewConfigurator(cdc.Codec, mockServer, mockServer)
	err = cfg.RegisterMigration("testModule", 1, func(_ sdk.Context) error {
		isCalled = true
		return nil
	})
	require.NoError(t, err)

	err = mm.RunMigrations(sdk.NewContext(nil, tmproto.Header{}, false, log.NewNopLogger()), cfg, 1, 2)
	require.NoError(t, err)
	require.True(t, isCalled)

	supportedVersions := mm.SupportedVersions()
	require.Len(t, supportedVersions, 3)
	require.Contains(t, supportedVersions, uint64(1))
	require.Contains(t, supportedVersions, uint64(2))
	require.Contains(t, supportedVersions, uint64(3))
}

func TestConfiguratorRegistersAllMessageTypes(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockServer := mocks.NewMockServer(mockCtrl)
	mockServer.EXPECT().RegisterService(gomock.Any(), gomock.Any()).Times(2).Return()
	cdc := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	configurator := module.NewConfigurator(cdc.Codec, mockServer, mockServer)

	storeKey := sdk.NewKVStoreKey(signaltypes.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	keeper := signal.NewKeeper(storeKey, nil)
	upgradeModule := signal.NewAppModule(keeper)
	mm, err := module.NewManager([]module.VersionedModule{
		{Module: upgradeModule, FromVersion: 2, ToVersion: 2},
	})
	require.NoError(t, err)
	require.NotNil(t, mm)

	mm.RegisterServices(configurator)
	acceptedMessages := configurator.GetAcceptedMessages()
	require.Equal(t, map[uint64]map[string]struct{}{
		2: {"/celestia.signal.v1.MsgSignalVersion": {}, "/celestia.signal.v1.MsgTryUpgrade": {}},
	}, acceptedMessages)

	require.NotNil(t, keeper)
}
