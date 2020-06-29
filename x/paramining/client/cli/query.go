package cli

import (
	"fmt"

	"github.com/Shushsa/plan/x/paramining/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	planQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the paramining module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	planQueryCmd.AddCommand(flags.GetCommands(
		GetCmdGetParamining(cdc),
	)...)

	return planQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdGetParamining(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [address]",
		Short: "get address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			_, err := sdk.AccAddressFromBech32(address)

			if err != nil {
				fmt.Printf("Wrong address %s \n", address)
				return nil
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/paramining/get/%s", address), nil)

			if err != nil {
				fmt.Printf("Cannot get paramining of %s \n", address)
				return nil
			}

			var out types.ParaminingResolve

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}
