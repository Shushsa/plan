package cli

import (
	"fmt"

	"github.com/Shushsa/plan/x/emission/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	nameserviceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the paramining module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	nameserviceQueryCmd.AddCommand(flags.GetCommands(
		GetCmdGetParamining(cdc),
	)...)

	return nameserviceQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdGetParamining(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "get",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData("custom/emission/get", nil)

			if err != nil {
				fmt.Println("Cannot get emission\n")

				return nil
			}

			var out types.Emission

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}
