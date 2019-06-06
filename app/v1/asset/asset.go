package asset

import sdk "github.com/irisnet/irishub/types"

type Asset struct {
	Family     string           `json:"family"`
	Name       string           `json:"name"`
	Symbol     string           `json:"symbol"`
	Source     string           `json:"source"`
	InitSupply uint64           `json:"init_supply"`
	MaxSupply  uint64           `json:"max_supply"`
	Decimal    uint8            `json:"decimal"`
	Mintable   bool             `json:"mintable"`
	Owner      sdk.AccAddress   `json:"owner"`
	Operators  []sdk.AccAddress `json:"operators"`
}
