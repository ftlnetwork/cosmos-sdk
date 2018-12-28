package cli

import (
	"os"

	"github.com/ftlnetwork/ftlnetwork-sdk/client/context"
	"github.com/ftlnetwork/ftlnetwork-sdk/client/utils"
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	authtxb "github.com/ftlnetwork/ftlnetwork-sdk/x/auth/client/txbuilder"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/slashing"

	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail",
		Args:  cobra.NoArgs,
		Short: "unjail validator previously jailed for downtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			valAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := slashing.NewMsgUnjail(sdk.ValAddress(valAddr))
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(os.Stdout, txBldr, cliCtx, []sdk.Msg{msg}, false)
			}
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
