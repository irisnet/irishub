package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GuardianI interface {
	GetDescription() string
	GetAccountType() AccountType
	GetAddress() sdk.AccAddress
	GetAddedBy() sdk.AccAddress
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
