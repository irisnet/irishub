package types

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/irisnet/irishub/types"
)

// ServiceBinding defines a struct for service binding
type ServiceBinding struct {
	ServiceName     string         `json:"service_name"`
	Provider        sdk.AccAddress `json:"provider"`
	Deposit         sdk.Coins      `json:"deposit"`
	Pricing         string         `json:"pricing"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
	Available       bool           `json:"available"`
	DisabledTime    time.Time      `json:"disabled_time"`
}

// NewServiceBinding creates a new ServiceBinding instance
func NewServiceBinding(
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
	withdrawAddr sdk.AccAddress,
	available bool,
	disabledTime time.Time,
) ServiceBinding {
	return ServiceBinding{
		ServiceName:     serviceName,
		Provider:        provider,
		Deposit:         deposit,
		Pricing:         pricing,
		WithdrawAddress: withdrawAddr,
		Available:       available,
		DisabledTime:    disabledTime,
	}
}

// String implements fmt.Stringer
func (binding ServiceBinding) String() string {
	disabledTimeStr := ""

	if !binding.DisabledTime.IsZero() {
		disabledTimeStr = binding.DisabledTime.Format("1970-01-01 00:00:00")
	}

	return fmt.Sprintf(`ServiceBinding:
		ServiceName:             %s
		Provider:                %s
		Deposit:                 %s
		Pricing:                 %s
		WithdrawAddress:         %s
		Available:               %v,
		DisabledTime:            %s`,
		binding.ServiceName, binding.Provider, binding.Deposit.MainUnitString(),
		binding.Pricing, binding.WithdrawAddress,
		binding.Available, disabledTimeStr,
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
	Discount  float32   `json:"discount"`   // discount during the promotion
}

// PromotionByVolume defines the promotion activity by volume
type PromotionByVolume struct {
	Volume   uint64  `json:"volume"`   // minimal volume for the promotion
	Discount float32 `json:"discount"` // discount for the promotion
}

// ParsePricing parses the given pricing string
func ParsePricing(pricing string) (Pricing, sdk.Error) {
	var p Pricing
	if err := json.Unmarshal([]byte(pricing), &p); err != nil {
		return p, ErrInvalidPricing(DefaultCodespace, fmt.Sprintf("failed to unmarshal the pricing: %s", err))
	}

	return p, nil
}
