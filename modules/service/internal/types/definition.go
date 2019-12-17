package types

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SvcDef struct {
	Name              string         `json:"name" yaml:"name"`
	ChainId           string         `json:"chain_id" yaml:"chain_id"`
	Description       string         `json:"description" yaml:"description"`
	Tags              []string       `json:"tags" yaml:"tags"`
	Author            sdk.AccAddress `json:"author" yaml:"author"`
	AuthorDescription string         `json:"author_description" yaml:"author_description"`
	IDLContent        string         `json:"idl_content" yaml:"idl_content"`
}

type MethodProperty struct {
	ID            int16             `json:"id" yaml:"id"`
	Name          string            `json:"name" yaml:"name"`
	Description   string            `json:"description" yaml:"description"`
	OutputPrivacy OutputPrivacyEnum `json:"output_privacy" yaml:"output_privacy"`
	OutputCached  OutputCachedEnum  `json:"output_cached" yaml:"output_cached"`
}

func NewSvcDef(name, chainId, description string, tags []string, author sdk.AccAddress, authorDescription, idlContent string) SvcDef {
	return SvcDef{
		Name:              name,
		ChainId:           chainId,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		IDLContent:        idlContent,
	}
}

type OutputPrivacyEnum byte

const (
	NoPrivacy        OutputPrivacyEnum = 0x01
	PubKeyEncryption OutputPrivacyEnum = 0x02
)

type OutputCachedEnum byte

const (
	OffChainCached OutputCachedEnum = 0x01
	NoCached       OutputCachedEnum = 0x02
)

type MessagingType byte

const (
	Unicast   MessagingType = 0x01
	Multicast MessagingType = 0x02
)

// String to messagingType byte, Returns ff if invalid.
func MessagingTypeFromString(str string) (MessagingType, error) {
	switch str {
	case "Multicast":
		return Multicast, nil
	case "Unicast":
		return Unicast, nil
	default:
		return MessagingType(0xff), errors.Errorf("'%s' is not a valid messaging type", str)
	}
}

// is defined messagingType?
func validMessagingType(mt MessagingType) bool {
	if mt == Multicast ||
		mt == Unicast {
		return true
	}
	return false
}

// For Printf / Sprintf, returns bech32 when using %s
func (mt MessagingType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", mt.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(mt))))
	}
}

// Turns MessagingType byte to String
func (mt MessagingType) String() string {
	switch mt {
	case Multicast:
		return "Multicast"
	case Unicast:
		return "Unicast"
	default:
		return ""
	}
}

// Marshals to JSON using string
func (mt MessagingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (mt *MessagingType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := MessagingTypeFromString(s)
	if err != nil {
		return err
	}
	*mt = bz2
	return nil
}

// String to outputCachedEnum byte, Returns ff if invalid.
func OutputCachedEnumFromString(str string) (OutputCachedEnum, error) {
	switch str {
	case "OffChainCached":
		return OffChainCached, nil
	case "NoCached":
		return NoCached, nil
	default:
		return OutputCachedEnum(0xff), errors.Errorf("'%s' is not a valid outputCachedEnum type", str)
	}
}

// is defined OutputCachedEnum?
func validOutputCachedEnum(oe OutputCachedEnum) bool {
	if oe == OffChainCached ||
		oe == NoCached {
		return true
	}
	return false
}

// For Printf / Sprintf, returns bech32 when using %s
func (oe OutputCachedEnum) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", oe.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(oe))))
	}
}

// Turns OutputCachedEnum byte to String
func (oe OutputCachedEnum) String() string {
	switch oe {
	case OffChainCached:
		return "OffChainCached"
	case NoCached:
		return "NoCached"
	default:
		return ""
	}
}

// Marshals to JSON using string
func (oe OutputCachedEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(oe.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (oe *OutputCachedEnum) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := OutputCachedEnumFromString(s)
	if err != nil {
		return err
	}
	*oe = bz2
	return nil
}

// String to outputPrivacyEnum byte, Returns ff if invalid.
func OutputPrivacyEnumFromString(str string) (OutputPrivacyEnum, error) {
	switch str {
	case "NoPrivacy":
		return NoPrivacy, nil
	case "PubKeyEncryption":
		return PubKeyEncryption, nil
	default:
		return OutputPrivacyEnum(0xff), errors.Errorf("'%s' is not a valid outputPrivacyEnum type", str)
	}
}

// is defined OutputPrivacyEnum?
func validOutputPrivacyEnum(oe OutputPrivacyEnum) bool {
	if oe == NoPrivacy ||
		oe == PubKeyEncryption {
		return true
	}
	return false
}

// For Printf / Sprintf, returns bech32 when using %s
func (oe OutputPrivacyEnum) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", oe.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(oe))))
	}
}

// Turns OutputCachedEnum byte to String
func (oe OutputPrivacyEnum) String() string {
	switch oe {
	case NoPrivacy:
		return "NoPrivacy"
	case PubKeyEncryption:
		return "PubKeyEncryption"
	default:
		return ""
	}
}

// Marshals to JSON using string
func (oe OutputPrivacyEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(oe.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (oe *OutputPrivacyEnum) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := OutputPrivacyEnumFromString(s)
	if err != nil {
		return err
	}
	*oe = bz2
	return nil
}
