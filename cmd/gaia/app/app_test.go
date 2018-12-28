package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/ftlnetwork/ftlnetwork-sdk/codec"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/auth"
	distr "github.com/ftlnetwork/ftlnetwork-sdk/x/distribution"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/gov"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/mint"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/slashing"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/stake"

	abci "github.com/tendermint/tendermint/abci/types"
)

func setGenesis(gapp *GaiaApp, accs ...*auth.BaseAccount) error {
	genaccs := make([]GenesisAccount, len(accs))
	for i, acc := range accs {
		genaccs[i] = NewGenesisAccount(acc)
	}

	genesisState := NewGenesisState(
		genaccs,
		auth.DefaultGenesisState(),
		stake.DefaultGenesisState(),
		mint.DefaultGenesisState(),
		distr.DefaultGenesisState(),
		gov.DefaultGenesisState(),
		slashing.DefaultGenesisState(),
	)

	stateBytes, err := codec.MarshalJSONIndent(gapp.cdc, genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	vals := []abci.ValidatorUpdate{}
	gapp.InitChain(abci.RequestInitChain{Validators: vals, AppStateBytes: stateBytes})
	gapp.Commit()

	return nil
}

func TestGaiadExport(t *testing.T) {
	db := db.NewMemDB()
	gapp := NewGaiaApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true)
	setGenesis(gapp)

	// Making a new app object with the db, so that initchain hasn't been called
	newGapp := NewGaiaApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true)
	_, _, err := newGapp.ExportAppStateAndValidators(false)
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
