package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ouroboros-crypto/node/x/structure/types"
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"

)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	ouroborosQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the structure module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	ouroborosQueryCmd.AddCommand(client.GetCommands(
		GetCmdGetStructure(cdc),
		GetGemdGetUpperStructure(cdc),
	)...)

	return ouroborosQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdGetStructure(cdc *codec.Codec) *cobra.Command {
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

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/structure/get/%s", address), nil)

			if err != nil {
				fmt.Printf("Cannot fetch the structure %s \n", address)
				return nil
			}

			var out types.Structure

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}
// GetCmdResolveName queries information about a name
func GetGemdGetUpperStructure(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "upper [address]",
		Short: "upper address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			address := args[0]

			_, err := sdk.AccAddressFromBech32(address)

			if err != nil {
				fmt.Printf("Wrong address %s \n", address)
				return nil
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/structure/upper/%s", address), nil)

			if err != nil {
				fmt.Printf("Cannot fetch the structure %s \n", address)
				return nil
			}

			var out types.UpperStructure

			cdc.MustUnmarshalJSON(res, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}