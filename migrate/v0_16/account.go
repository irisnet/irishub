package v0_16

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/migrate/v0_16/types"
)

type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         types.Coins    `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`
}
