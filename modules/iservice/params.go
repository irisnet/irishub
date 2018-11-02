package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	"github.com/irisnet/irishub/types"
)

// Params defines the high level settings for iservice
type Params struct {
	MaxTagsNum         int
	MaxRequestTimeout  int
	MinProviderDeposit sdk.Coins
	SlashDeposit       sdk.Coins
}

var iserviceParams Params

func init() {
	iserviceParams = DefaultParams()
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	minDeposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 1000, "iris"))
	slashDeposit, _ := types.NewDefaultCoinType("iris").ConvertToMinCoin(fmt.Sprintf("%d%s", 2, "iris"))
	return Params{
		MaxTagsNum:         5,
		MaxRequestTimeout:  100,
		MinProviderDeposit: sdk.Coins{minDeposit},
		SlashDeposit:       sdk.Coins{slashDeposit},
	}
}
