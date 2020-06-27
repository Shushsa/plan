package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/Shushsa/plan/x/paramining/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	ouroTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Paramining transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ouroTxCmd.AddCommand(client.PostCommands(
		GetCmdReinvest(cdc),
	)...)

	return ouroTxCmd
}

func GetCmdReinvest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reinvest",
		Short: "Gets paramining by sending a reinvest transaction",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			msg := types.NewMsgReinvest(cliCtx.GetFromAddress())
			err := msg.ValidateBasic()

			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
