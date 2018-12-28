package tags

import (
	sdk "github.com/ftlnetwork/ftlnetwork-sdk/types"
)

// Slashing tags
var (
	ActionValidatorUnjailed = []byte("validator-unjailed")

	Action    = sdk.TagAction
	Validator = "validator"
)
