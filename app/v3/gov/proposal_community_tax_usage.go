package gov

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/irisnet/irishub/app/v1/auth"
	sdk "github.com/irisnet/irishub/types"
)

type UsageType byte

const (
	UsageTypeBurn       UsageType = 0x01
	UsageTypeDistribute UsageType = 0x02
	UsageTypeGrant      UsageType = 0x03
)

// String to UsageType byte.  Returns ff if invalid.
func UsageTypeFromString(str string) (UsageType, error) {
	switch str {
	case "Burn":
		return UsageTypeBurn, nil
	case "Distribute":
		return UsageTypeDistribute, nil
	case "Grant":
		return UsageTypeGrant, nil
	default:
		return UsageType(0xff), errors.Errorf("'%s' is not a valid usage type", str)
	}
}

// is defined UsageType?
func ValidUsageType(ut UsageType) bool {
	if ut == UsageTypeBurn ||
		ut == UsageTypeDistribute ||
		ut == UsageTypeGrant {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (ut UsageType) Marshal() ([]byte, error) {
	return []byte{byte(ut)}, nil
}

// Unmarshal needed for protobuf compatibility
func (ut *UsageType) Unmarshal(data []byte) error {
	*ut = UsageType(data[0])
	return nil
}

// Marshals to JSON using string
func (ut UsageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ut.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (ut *UsageType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := UsageTypeFromString(s)
	if err != nil {
		return err
	}
	*ut = bz2
	return nil
}

// Turns VoteOption byte to String
func (ut UsageType) String() string {
	switch ut {
	case UsageTypeBurn:
		return "Burn"
	case UsageTypeDistribute:
		return "Distribute"
	case UsageTypeGrant:
		return "Grant"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (ut UsageType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(ut.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(ut))))
	}
}

// Implements Proposal Interface
var _ Proposal = (*CommunityTaxUsageProposal)(nil)

type TaxUsage struct {
	Usage       UsageType      `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Percent     sdk.Dec        `json:"percent"`
	Amount      sdk.Coins      `json:"amount"`
}

type CommunityTaxUsageProposal struct {
	BasicProposal
	TaxUsage TaxUsage `json:"tax_usage"`
}

func (tp CommunityTaxUsageProposal) GetTaxUsage() TaxUsage { return tp.TaxUsage }
func (tp *CommunityTaxUsageProposal) SetTaxUsage(taxUsage TaxUsage) {
	tp.TaxUsage = taxUsage
}

func (tp *CommunityTaxUsageProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if err := tp.BasicProposal.Validate(ctx, k, verify); err != nil {
		return err
	}

	taxCoins := k.ck.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	if !taxCoins.IsAllGTE(tp.TaxUsage.Amount) {
		return ErrNotEnoughCommunityTax(k.codespace, taxCoins, tp.TaxUsage.Amount)
	}

	// only check trustee address for distribute usage
	if tp.TaxUsage.Usage == UsageTypeDistribute {
		_, found := k.guardianKeeper.GetTrustee(ctx, tp.TaxUsage.DestAddress)
		if !found {
			return ErrNotTrustee(k.codespace, tp.TaxUsage.DestAddress)
		}
	}
	return nil
}

func (tp *CommunityTaxUsageProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	logger := ctx.Logger()
	if err := tp.Validate(ctx, gk, false); err != nil {
		logger.Error("Execute CommunityTaxUsageProposal Failure", "info",
			"the destination address is not a trustee now", "destinationAddress", tp.TaxUsage.DestAddress)
		return err
	}
	burn := false
	if tp.TaxUsage.Usage == UsageTypeBurn {
		burn = true
	}
	gk.AllocateFeeTax(ctx, tp.TaxUsage.DestAddress, tp.TaxUsage.Amount, burn)
	return nil
}

// Allocate fee tax from the community fee pool, burn or send to trustee account
func (keeper Keeper) AllocateFeeTax(ctx sdk.Context, destAddr sdk.AccAddress, amount sdk.Coins, burn bool) {
	logger := ctx.Logger()
	taxCoins := keeper.ck.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	taxLeftCoins, hasNeg := taxCoins.SafeSub(amount)
	if hasNeg {
		logger.Info(fmt.Sprintf("community tax account [%s] is not enough to cover usage amount [%s]", taxCoins, amount))
		return
	}

	logger.Info("Spend community tax fund", "total_community_tax_fund", taxCoins.String(), "left_community_tax_fund", taxLeftCoins.String())
	if burn {
		logger.Info("Burn community tax", "burn_amount", amount.String())
		_, err := keeper.ck.BurnCoins(ctx, auth.CommunityTaxCoinsAccAddr, amount)
		if err != nil {
			panic(err)
		}
		if !amount.IsZero() {
			ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.CommunityTaxCoinsAccAddr.String(), "", amount.String(), sdk.BurnFlow, "")
		}
	} else {
		logger.Info("Grant community tax to account", "grant_amount", amount.String(), "grant_address", destAddr.String())
		if !amount.IsZero() {
			ctx.CoinFlowTags().AppendCoinFlowTag(ctx, auth.CommunityTaxCoinsAccAddr.String(), destAddr.String(), amount.String(), sdk.CommunityTaxUseFlow, "")
		}
		_, err := keeper.ck.SendCoins(ctx, auth.CommunityTaxCoinsAccAddr, destAddr, amount)
		if err != nil {
			panic(err)
		}
	}

}
