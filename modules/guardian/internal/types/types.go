package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/exported"
)

// Guardian represents a Guardian
type Guardian struct {
	Description string         `json:"description" yaml:"description"` // description of guardian
	AccountType AccountType    `json:"type"  yaml:"account_type"`      // account type of guardian
	Address     sdk.AccAddress `json:"address" yaml:"address"`         // address of guardian
	AddedBy     sdk.AccAddress `json:"added_by"  yaml:"added_by"`      // address that initiated the AddGuardian tx
}

// Profilers is a collection of Guardian
type Profilers []Guardian

// String implements string
func (ps Profilers) String() (out string) {
	if len(ps) == 0 {
		return "[]"
	}
	for _, val := range ps {
		out += fmt.Sprintf(`Profiler
  Address:       %s
  Type:          %s
  Description:   %s
  AddedBy:       %s
`, val.Address, val.AccountType, val.Description, val.AddedBy)
	}
	return strings.TrimSpace(out)
}

// Trustees is a collection of Guardian
type Trustees []Guardian

// String implements string
func (ts Trustees) String() (out string) {
	if len(ts) == 0 {
		return "[]"
	}
	for _, val := range ts {
		out += fmt.Sprintf(`Trustee
  Address:       %s
  Type:          %s
  Description:   %s
  AddedBy:       %s
`, val.Address, val.AccountType, val.Description, val.AddedBy)
	}
	return strings.TrimSpace(out)
}

// NewGuardian constructs a Guardian
func NewGuardian(description string, accountType AccountType, address, addedBy sdk.AccAddress) Guardian {
	return Guardian{
		Description: description,
		AccountType: accountType,
		Address:     address,
		AddedBy:     addedBy,
	}
}

// Equal returns if the guardian is equal to specified guardian
func (g Guardian) Equal(guardian Guardian) bool {
	return g.Address.Equals(guardian.Address) &&
		g.AddedBy.Equals(guardian.AddedBy) &&
		g.Description == guardian.Description &&
		g.AccountType == guardian.AccountType
}

// GetDescription returns Description of the guardian
func (g Guardian) GetDescription() string {
	return g.Description
}

// GetAccountType returns AccountType of the guardian
func (g Guardian) GetAccountType() exported.AccountTypeI {
	return g.AccountType
}

// GetAddress returns Address of the guardian
func (g Guardian) GetAddress() sdk.AccAddress {
	return g.Address
}

// GetAddedBy returns AddedBy of the guardian
func (g Guardian) GetAddedBy() sdk.AccAddress {
	return g.AddedBy
}

// AccountType represents the type of account
type AccountType byte

const (
	Genesis  AccountType = 0x01 // account type of genesis
	Ordinary AccountType = 0x02 // account type of ordinary
)

// AccountTypeFromString converts string to AccountType byte, Returns ff if invalid.
func AccountTypeFromString(str string) (AccountType, error) {
	switch str {
	case "Genesis":
		return Genesis, nil
	case "Ordinary":
		return Ordinary, nil
	default:
		return AccountType(0xff), errors.Errorf("'%s' is not a valid account type", str)
	}
}

// is defined AccountType?
func validAccountType(bt AccountType) bool {
	return bt == Genesis || bt == Ordinary
}

// Format for Printf / Sprintf, returns bech32 when using %s
func (bt AccountType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(fmt.Sprintf("%s", bt.String())))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%v", byte(bt))))
	}
}

// String converts BindingType byte to String
func (bt AccountType) String() string {
	switch bt {
	case Genesis:
		return "Genesis"
	case Ordinary:
		return "Ordinary"
	default:
		return ""
	}
}

// MarshalJSON marshals AccountType to JSON using string
func (bt AccountType) MarshalJSON() ([]byte, error) {
	return json.Marshal(bt.String())
}

// UnmarshalJSON unmarshals AccountType from JSON assuming Bech32 encoding
func (bt *AccountType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz2, err := AccountTypeFromString(s)
	if err != nil {
		return err
	}
	*bt = bz2
	return nil
}
