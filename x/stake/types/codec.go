package types

import (
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "ftlnetwork-sdk/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "ftlnetwork-sdk/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "ftlnetwork-sdk/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgBeginUnbonding{}, "ftlnetwork-sdk/BeginUnbonding", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "ftlnetwork-sdk/BeginRedelegate", nil)
}

// generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
