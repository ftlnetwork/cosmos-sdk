package simulation

import (
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/mock/simulation"
)

// TODO Any invariants to check here?
// AllInvariants tests all slashing invariants
func AllInvariants() simulation.Invariant {
	return func(_ sdk.Context) error {
		return nil
	}
}
