package app_test

import (
	"fmt"
	"testing"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func TestApp(t *testing.T) {
	encCfg := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	celestiaApp := app.New(log.NewNopLogger(), db.NewMemDB(), nil, true, 0, encCfg, 0, util.EmptyAppOptions{})

	t.Run("LoadHeight", func(t *testing.T) {
		type testCase struct {
			height int64
			want   error
		}
		testCases := []testCase{
			{height: -1, want: nil},
			{height: 0, want: nil},
			{height: 1, want: nil},
		}

		for _, tc := range testCases {
			name := fmt.Sprintf("height=%d", tc.height)
			t.Run(name, func(t *testing.T) {
				got := celestiaApp.LoadHeight(tc.height)
				assert.Equal(t, tc.want, got)
			})
		}
	})
}
