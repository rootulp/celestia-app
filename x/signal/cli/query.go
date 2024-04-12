package cli

import (
	"fmt"
	"strconv"

	"github.com/celestiaorg/celestia-app/v2/x/signal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the CLI query commands for this module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryTally())
	return cmd
}

func CmdQueryTally() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tally version",
		Short:   "Query for the tally of voting power that has signalled for a particular version",
		Args:    cobra.ExactArgs(1),
		Example: "tally 3",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			version, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			upgradeQueryClient := types.NewQueryClient(clientCtx)
			resp, err := upgradeQueryClient.VersionTally(cmd.Context(), &types.QueryVersionTallyRequest{Version: version})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
