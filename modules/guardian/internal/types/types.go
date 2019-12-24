package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian/exported"
)

type Guardian struct {
	Description string         `json:"description" yaml:"description"`
	AccountType AccountType    `json:"type"  yaml:"account_type"`
	Address     sdk.AccAddress `json:"address" yaml:"address"`    // this guardian's address
	AddedBy     sdk.AccAddress `json:"added_by"  yaml:"added_by"` // address that initiated the AddGuardian tx
}

type Profilers []Guardian

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

type Trustees []Guardian

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

func NewGuardian(description string, accountType AccountType, address, addedBy sdk.AccAddress) Guardian {
	return Guardian{
		Description: description,
		AccountType: accountType,
		Address:     address,
		AddedBy:     addedBy,
	}
}

func (g Guardian) Equal(guardian Guardian) bool {
	return g.Address.Equals(guardian.Address) &&
		g.AddedBy.Equals(guardian.AddedBy) &&
		g.Description == guardian.Description &&
		g.AccountType == guardian.AccountType
}

func (g Guardian) GetDescription() string {
	return g.Description
}

func (g Guardian) GetAccountType() exported.AccountTypeI {
	return g.AccountType
}

func (g Guardian) GetAddress() sdk.AccAddress {
	return g.Address
}

func (g Guardian) GetAddedBy() sdk.AccAddress {
	return g.AddedBy
}

type AccountType byte

const (
	Genesis  AccountType = 0x01
	Ordinary AccountType = 0x02
)

// String to AccountType byte, Returns ff if invalid.
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

// For Printf / Sprintf, returns bech32 when using %s
func (bt AccountType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", bt.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(bt))))
	}
}

// Turns BindingType byte to String
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

// Marshals to JSON using string
func (bt AccountType) MarshalJSON() ([]byte, error) {
	return json.Marshal(bt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
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
