package types

import (
	"encoding/json"
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

// ParsePricing parses the given pricing string
func ParsePricing(pricing string) (Pricing, sdk.Error) {
	var p Pricing
	if err := json.Unmarshal([]byte(pricing), &p); err != nil {
		return p, ErrInvalidPricing(DefaultCodespace, fmt.Sprintf("failed to unmarshal the pricing: %s", err))
	}

	return p, nil
}

// GetPrice gets the current price by the specified volume and time
// Note: ensure that the pricing is valid
func GetPrice(pricing string, time time.Time, volume uint64) sdk.Coins {
	p, _ := ParsePricing(pricing)

	// get discounts
	discountByTime := GetDiscountByTime(p.PromotionsByTime, time)
	discountByVolume := GetDiscountByVolume(p.PromotionsByVolume, volume)

	// compute the price
	basePrice := p.Price.AmountOf(sdk.IrisAtto)
	price := sdk.NewDecFromInt(basePrice).Mul(discountByTime).Mul(discountByVolume)

	return sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, price.TruncateInt()))
}

// GetDiscountByTime gets the current discount level by the specified time
func GetDiscountByTime(promotionsByTime []PromotionByTime, time time.Time) sdk.Dec {
	for _, p := range promotionsByTime {
		if (time.Equal(p.StartTime) || time.After(p.StartTime)) &&
			(time.Equal(p.EndTime) || time.Before(p.EndTime)) {
			return p.Discount
		}
	}

	return sdk.OneDec()
}

// GetDiscountByVolume gets the current discount level by the specified volume
func GetDiscountByVolume(promotionsByVolume []PromotionByVolume, volume uint64) sdk.Dec {
	if len(promotionsByVolume) > 0 {
		if volume < promotionsByVolume[0].Volume {
			return sdk.OneDec()
		}

		if volume >= promotionsByVolume[len(promotionsByVolume)-1].Volume {
			return promotionsByVolume[len(promotionsByVolume)-1].Discount
		}
	}

	for i, p := range promotionsByVolume {
		if volume < p.Volume {
			return promotionsByVolume[i-1].Discount
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

	if !validServiceCoins(binding.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", binding.Deposit))
	}

	if len(binding.Pricing) == 0 {
		return ErrInvalidPricing(DefaultCodespace, "pricing missing")
	}

	return validatePricing(binding.Pricing)
}
