package client

import (
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	bank "github.com/ftlnetwork/ftlnetwork-sdk/x/bank"
)

// create the sendTx msg
func CreateMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}
