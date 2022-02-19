package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MT multi token interface
type MT interface {
	GetID() string
	GetSupply() uint64
	GetOwner() sdk.AccAddress
	GetData() []byte
}
