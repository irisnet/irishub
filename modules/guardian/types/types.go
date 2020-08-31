package types

import (
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

// ValidAccountType returns true if the AccountType option is valid and false otherwise.
func ValidAccountType(option AccountType) bool {
	if option == Genesis ||
		option == Ordinary {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility.
func (at AccountType) Marshal() ([]byte, error) {
	return []byte{byte(at)}, nil
}

// Unmarshal needed for protobuf compatibility.
func (at *AccountType) Unmarshal(data []byte) error {
	*at = AccountType(data[0])
	return nil
}

// Format implements the fmt.Formatter interface.
func (at AccountType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(at.String()))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(at))))
	}
}
