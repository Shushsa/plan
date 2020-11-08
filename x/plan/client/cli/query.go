package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/plan-crypto/node/x/plan/types"
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"

)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	planQueryCmd := &cobra.Command{
		Use:                        "plan",
		Short:                      "Querying commands for the plan module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	planQueryCmd.AddCommand(client.GetCommands(
		GetCmdProfile(cdc),
	)...)

	return planQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdProfile(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "profile [address]",
		Short: "profile address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			_, err := sdk.AccAddressFromBech32(address)

			if err != nil {
				fmt.Printf("Wrong address %s \n", address)
				return nil
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/plan/profile/%s", address), nil)

			if err != nil {
				fmt.Printf("Cannot get profile %s \n", address)
				return nil
			}

			var out types.ProfileResolve

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}