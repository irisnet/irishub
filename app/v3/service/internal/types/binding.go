package types

import (
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
	return fmt.Sprintf(`ServiceBinding:
		ServiceName:             %s
		Provider:                %s
		Deposit:                 %s
		Pricing:                 %s
		WithdrawAddress:         %s
		Available:               %v,
		DisabledTime:            %s`,
		binding.ServiceName, binding.Provider, binding.Deposit,
		binding.Pricing, binding.WithdrawAddress, binding.Available,
		binding.DisabledTime)
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
