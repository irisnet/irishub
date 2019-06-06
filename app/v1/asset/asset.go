package asset

import sdk "github.com/irisnet/irishub/types"

type Asset struct {
	Family     FamilyKind       `json:"family"`
	Source     string           `json:"source"`
	Symbol     string           `json:"symbol"`
	Name       string           `json:"name"`
	InitSupply uint64           `json:"init_supply"`
	MaxSupply  uint64           `json:"max_supply"`
	Decimal    uint8            `json:"decimal"`
	Mintable   bool             `json:"mintable"`
	Owner      sdk.AccAddress   `json:"owner"`
	Operators  []sdk.AccAddress `json:"operators"`
}
