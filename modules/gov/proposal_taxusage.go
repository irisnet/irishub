package gov

import (
	"encoding/json"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
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
		return nil
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
var _ Proposal = (*ParameterProposal)(nil)

type TaxUsage struct {
	Usage       UsageType      `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Percent     sdk.Dec        `json:"percent"`
}

type TaxUsageProposal struct {
	BasicProposal
	TaxUsage TaxUsage `json:"tax_usage"`
}

func (tp TaxUsageProposal) GetTaxUsage() TaxUsage { return tp.TaxUsage }
func (tp *TaxUsageProposal) SetTaxUsage(taxUsage TaxUsage) {
	tp.TaxUsage = taxUsage
}
