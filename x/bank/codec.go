package bank

import (
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "ftlnetwork-sdk/Send", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
