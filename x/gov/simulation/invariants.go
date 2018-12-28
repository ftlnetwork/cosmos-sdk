package simulation

import (
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
	"github.com/ftlnetwork/ftlnetwork-sdk/x/mock/simulation"
)

// AllInvariants tests all governance invariants
func AllInvariants() simulation.Invariant {
	return func(ctx sdk.Context) error {
		// TODO Add some invariants!
		// Checking proposal queues, no passed-but-unexecuted proposals, etc.
		return nil
	}
}
