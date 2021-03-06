package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	"github.com/ftlnetwork/ftlnetwork-sdk/store"
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/auth"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/bank"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/params"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/stake"

	"github.com/ftlnetwork/ftlnetwork-sdk/x/distribution/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delPk2   = ed25519.GenPrivKey().PubKey()
	delPk3   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
	delAddr2 = sdk.AccAddress(delPk2.Address())
	delAddr3 = sdk.AccAddress(delPk3.Address())

	valOpPk1    = ed25519.GenPrivKey().PubKey()
	valOpPk2    = ed25519.GenPrivKey().PubKey()
	valOpPk3    = ed25519.GenPrivKey().PubKey()
	valOpAddr1  = sdk.ValAddress(valOpPk1.Address())
	valOpAddr2  = sdk.ValAddress(valOpPk2.Address())
	valOpAddr3  = sdk.ValAddress(valOpPk3.Address())
	valAccAddr1 = sdk.AccAddress(valOpPk1.Address()) // generate acc addresses for these validator keys too
	valAccAddr2 = sdk.AccAddress(valOpPk2.Address())
	valAccAddr3 = sdk.AccAddress(valOpPk3.Address())

	valConsPk1   = ed25519.GenPrivKey().PubKey()
	valConsPk2   = ed25519.GenPrivKey().PubKey()
	valConsPk3   = ed25519.GenPrivKey().PubKey()
	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())
	valConsAddr3 = sdk.ConsAddress(valConsPk3.Address())

	addrs = []sdk.AccAddress{
		delAddr1, delAddr2, delAddr3,
		valAccAddr1, valAccAddr2, valAccAddr3,
	}

	emptyDelAddr sdk.AccAddress
	emptyValAddr sdk.ValAddress
	emptyPubkey  crypto.PubKey
)

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // distr
	return cdc
}

// test input with default values
func CreateTestInputDefault(t *testing.T, isCheckTx bool, initCoins int64) (
	sdk.Context, auth.AccountKeeper, Keeper, stake.Keeper, DummyFeeCollectionKeeper) {

	communityTax := sdk.NewDecWithPrec(2, 2)
	return CreateTestInputAdvanced(t, isCheckTx, initCoins, communityTax)
}

// hogpodge of all sorts of input required for testing
func CreateTestInputAdvanced(t *testing.T, isCheckTx bool, initCoins int64,
	communityTax sdk.Dec) (
	sdk.Context, auth.AccountKeeper, Keeper, stake.Keeper, DummyFeeCollectionKeeper) {

	keyDistr := sdk.NewKVStoreKey(types.StoreKey)
	keyStake := sdk.NewKVStoreKey(stake.StoreKey)
	tkeyStake := sdk.NewTransientStoreKey(stake.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyFeeCollection := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStake, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStake, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyFeeCollection, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	cdc := MakeTestCodec()
	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	ck := bank.NewBaseKeeper(accountKeeper)
	sk := stake.NewKeeper(cdc, keyStake, tkeyStake, ck, pk.Subspace(stake.DefaultParamspace), stake.DefaultCodespace)
	sk.SetPool(ctx, stake.InitialPool())
	sk.SetParams(ctx, stake.DefaultParams())

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range addrs {
		pool := sk.GetPool(ctx)
		_, _, err := ck.AddCoins(ctx, addr, sdk.Coins{
			{sk.GetParams(ctx).BondDenom, sdk.NewInt(initCoins)},
		})
		require.Nil(t, err)
		pool.LooseTokens = pool.LooseTokens.Add(sdk.NewDec(initCoins))
		sk.SetPool(ctx, pool)
	}

	fck := DummyFeeCollectionKeeper{}
	keeper := NewKeeper(cdc, keyDistr, pk.Subspace(DefaultParamspace), ck, sk, fck, types.DefaultCodespace)

	// set the distribution hooks on staking
	sk.SetHooks(keeper.Hooks())

	// set genesis items required for distribution
	keeper.SetFeePool(ctx, types.InitialFeePool())
	keeper.SetCommunityTax(ctx, communityTax)
	keeper.SetBaseProposerReward(ctx, sdk.NewDecWithPrec(1, 2))
	keeper.SetBonusProposerReward(ctx, sdk.NewDecWithPrec(4, 2))

	return ctx, accountKeeper, keeper, sk, fck
}

//__________________________________________________________________________________
// fee collection keeper used only for testing
type DummyFeeCollectionKeeper struct{}

var heldFees sdk.Coins
var _ types.FeeCollectionKeeper = DummyFeeCollectionKeeper{}

// nolint
func (fck DummyFeeCollectionKeeper) GetCollectedFees(_ sdk.Context) sdk.Coins {
	return heldFees
}
func (fck DummyFeeCollectionKeeper) SetCollectedFees(in sdk.Coins) {
	heldFees = in
}
func (fck DummyFeeCollectionKeeper) ClearCollectedFees(_ sdk.Context) {
	heldFees = sdk.Coins{}
}
