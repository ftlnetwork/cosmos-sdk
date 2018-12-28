package cli

import (
	"github.com/spf13/cobra"

	"github.com/ftlnetwork/ftlnetwork-sdk/client/context"
	"github.com/ftlnetwork/ftlnetwork-sdk/client/utils"
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	"github.com/ftlnetwork/ftlnetwork-sdk/docs/examples/democoin/x/cool"
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	authtxb "github.com/ftlnetwork/ftlnetwork-sdk/x/auth/client/txbuilder"
)

// QuizTxCmd invokes the coolness quiz transaction.
func QuizTxCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cool [answer]",
		Short: "What's cooler than being cool?",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := cool.NewMsgQuiz(from, args[0])

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

// SetTrendTxCmd sends a new cool trend transaction.
func SetTrendTxCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "setcool [answer]",
		Short: "You're so cool, tell us what is cool!",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := cool.NewMsgSetTrend(from, args[0])

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
