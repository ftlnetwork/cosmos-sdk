package ibc

import (
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(IBCTransferMsg{}, "ftlnetwork-sdk/IBCTransferMsg", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "ftlnetwork-sdk/IBCReceiveMsg", nil)
}
