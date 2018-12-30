package types

import (
	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgWithdrawDelegatorRewardsAll{}, "ftlnetwork-sdk/MsgWithdrawDelegationRewardsAll", nil)
	cdc.RegisterConcrete(MsgWithdrawDelegatorReward{}, "ftlnetwork-sdk/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(MsgWithdrawValidatorRewardsAll{}, "ftlnetwork-sdk/MsgWithdrawValidatorRewardsAll", nil)
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "ftlnetwork-sdk/MsgModifyWithdrawAddress", nil)
}

// generic sealed codec to be used throughout module
var MsgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	MsgCdc = cdc.Seal()
}
