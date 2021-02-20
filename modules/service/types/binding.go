package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewServiceBinding creates a new ServiceBinding instance
func NewServiceBinding(
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	qos uint64,
	options string,
	available bool,
	disabledTime time.Time,
	owner sdk.AccAddress,
) ServiceBinding {
	return ServiceBinding{
		ServiceName:  serviceName,
		Provider:     provider.String(),
		Deposit:      deposit,
		Pricing:      pricing,
		QoS:          qos,
		Options:      options,
		Available:    available,
		DisabledTime: disabledTime,
		Owner:        owner.String(),
	}
}

// RawPricing represents the raw pricing of a service binding
type RawPricing struct {
	Price              string              `json:"price"`                // base price string
	PromotionsByTime   []PromotionByTime   `json:"promotions_by_time"`   // promotions by time
	PromotionsByVolume []PromotionByVolume `json:"promotions_by_volume"` // promotions by volume
}

// GetDiscountByTime gets the discount level by the specified time
func GetDiscountByTime(pricing Pricing, time time.Time) sdk.Dec {
	for _, p := range pricing.PromotionsByTime {
		if !time.Before(p.StartTime) && time.Before(p.EndTime) {
			return p.Discount
		}
	}

	return sdk.OneDec()
}

// GetDiscountByVolume gets the discount level by the specified volume
// Note: make sure that the promotions by volume are sorted in ascending order
func GetDiscountByVolume(pricing Pricing, volume uint64) sdk.Dec {
	promotionsByVol := pricing.PromotionsByVolume

	for i, p := range promotionsByVol {
		if volume < p.Volume {
			if i == 0 {
				return sdk.OneDec()
			}

			return promotionsByVol[i-1].Discount
		}

		if i == len(promotionsByVol)-1 {
			return p.Discount
		}
	}

	return sdk.OneDec()
}

// ParsePricing parses the given string to Pricing
func ParsePricing(pricing string) (p Pricing, err error) {
	var rawPricing RawPricing
	if err := json.Unmarshal([]byte(pricing), &rawPricing); err != nil {
		return p, sdkerrors.Wrapf(ErrInvalidPricing, "failed to unmarshal the pricing: %s", err.Error())
	}

	priceCoin, err := sdk.ParseCoinNormalized(rawPricing.Price)
	if err != nil {
		return p, sdkerrors.Wrapf(ErrInvalidPricing, "invalid price: %s", err.Error())
	}

	p.Price = sdk.Coins{priceCoin}
	p.PromotionsByTime = rawPricing.PromotionsByTime
	p.PromotionsByVolume = rawPricing.PromotionsByVolume

	return p, nil
}

// CheckPricing checks if the given pricing complies with the specific rules
func CheckPricing(pricing Pricing) error {
	// CONTRACT:
	// p.EndTime > p.StartTime
	// p[i].StartTime >= p[i-1].EndTime
	for i, p := range pricing.PromotionsByTime {
		if !p.EndTime.After(p.StartTime) || (i > 0 && p.StartTime.Before(pricing.PromotionsByTime[i-1].EndTime)) {
			return sdkerrors.Wrapf(ErrInvalidPricing, "invalid timing promotion %d", i)
		}
	}

	// CONTRACT:
	// p[i].Volume > p[i-1].Volume
	for i, p := range pricing.PromotionsByVolume {
		if i > 0 && p.Volume < pricing.PromotionsByVolume[i-1].Volume {
			return sdkerrors.Wrapf(ErrInvalidPricing, "invalid volume promotion %d", i)
		}
	}

	return nil
}

// Validate validates the service binding
func (binding ServiceBinding) Validate() error {
	if err := ValidateProvider(binding.Provider); err != nil {
		return err
	}

	if err := ValidateOwner(binding.Owner); err != nil {
		return err
	}

	if err := ValidateServiceName(binding.ServiceName); err != nil {
		return err
	}

	if err := ValidateServiceDeposit(binding.Deposit); err != nil {
		return err
	}

	if err := ValidateQoS(binding.QoS); err != nil {
		return err
	}

	if err := ValidateOptions(binding.Options); err != nil {
		return err
	}

	return ValidateBindingPricing(binding.Pricing)
}
