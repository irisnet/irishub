package types

import (
	"fmt"
	"time"

	sdk "github.com/irisnet/irishub/types"
)

// ServiceBinding defines a struct for service binding
type ServiceBinding struct {
	ServiceName  string         `json:"service_name"`
	Provider     sdk.AccAddress `json:"provider"`
	Deposit      sdk.Coins      `json:"deposit"`
	Pricing      string         `json:"pricing"`
	Available    bool           `json:"available"`
	DisabledTime time.Time      `json:"disabled_time"`
}

// NewServiceBinding creates a new ServiceBinding instance
func NewServiceBinding(
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	available bool,
	disabledTime time.Time,
) ServiceBinding {
	return ServiceBinding{
		ServiceName:  serviceName,
		Provider:     provider,
		Deposit:      deposit,
		Pricing:      pricing,
		Available:    available,
		DisabledTime: disabledTime,
	}
}

// String implements fmt.Stringer
func (binding ServiceBinding) String() string {
	return fmt.Sprintf(`ServiceBinding:
		ServiceName:             %s
		Provider:                %s
		Deposit:                 %s
		Pricing:                 %s
		Available:               %v,
		DisabledTime:            %v`,
		binding.ServiceName, binding.Provider, binding.Deposit.MainUnitString(),
		binding.Pricing, binding.Available, binding.DisabledTime,
	)
}

// ServiceBindings is a set of service bindings
type ServiceBindings []ServiceBinding

// String implements fmt.Stringer
func (bindings ServiceBindings) String() string {
	if len(bindings) == 0 {
		return "[]"
	}

	var str string
	for _, binding := range bindings {
		str += binding.String() + "\n"
	}

	return str
}

// RawPricing represents the raw pricing of a service binding
type RawPricing struct {
	Price              string              `json:"price"`                // base price string
	PromotionsByTime   []PromotionByTime   `json:"promotions_by_time"`   // promotions by time
	PromotionsByVolume []PromotionByVolume `json:"promotions_by_volume"` // promotions by volume
}

// Pricing represents the pricing of a service binding
type Pricing struct {
	Price              sdk.Coins           `json:"price"`                // base price
	PromotionsByTime   []PromotionByTime   `json:"promotions_by_time"`   // promotions by time
	PromotionsByVolume []PromotionByVolume `json:"promotions_by_volume"` // promotions by volume
}

// PromotionByTime defines the promotion activity by time
type PromotionByTime struct {
	StartTime time.Time `json:"start_time"` // starting time of the promotion
	EndTime   time.Time `json:"end_time"`   // ending time of the promotion
	Discount  sdk.Dec   `json:"discount"`   // discount during the promotion
}

// PromotionByVolume defines the promotion activity by volume
type PromotionByVolume struct {
	Volume   uint64  `json:"volume"`   // minimal volume for the promotion
	Discount sdk.Dec `json:"discount"` // discount for the promotion
}

// GetDiscountByTime gets the discount level by the specified time
func GetDiscountByTime(pricing Pricing, time time.Time) sdk.Dec {
	for _, p := range pricing.PromotionsByTime {
		if (time.Equal(p.StartTime) || time.After(p.StartTime)) &&
			(time.Equal(p.EndTime) || time.Before(p.EndTime)) {
			return p.Discount
		}
	}

	return sdk.OneDec()
}

// GetDiscountByVolume gets the discount level by the specified volume
func GetDiscountByVolume(pricing Pricing, volume uint64) sdk.Dec {
	promotionsByVol := pricing.PromotionsByVolume

	if len(promotionsByVol) > 0 {
		if volume < promotionsByVol[0].Volume {
			return sdk.OneDec()
		}

		if volume >= promotionsByVol[len(promotionsByVol)-1].Volume {
			return promotionsByVol[len(promotionsByVol)-1].Discount
		}

		for i, p := range promotionsByVol {
			if volume < p.Volume {
				return promotionsByVol[i-1].Discount
			}
		}
	}

	return sdk.OneDec()
}

func (binding ServiceBinding) Validate() sdk.Error {
	if len(binding.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(binding.ServiceName); err != nil {
		return err
	}

	if !ValidServiceCoins(binding.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", binding.Deposit))
	}

	if len(binding.Pricing) == 0 {
		return ErrInvalidPricing(DefaultCodespace, "pricing missing")
	}

	return ValidateBindingPricing(binding.Pricing)
}
