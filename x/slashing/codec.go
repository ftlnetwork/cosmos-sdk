package slashing

import (
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgUnjail{}, "ftlnetwork-sdk/MsgUnjail", nil)
}

var cdcEmpty = codec.New()
