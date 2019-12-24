package types

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SvcDef
type SvcDef struct {
	Name              string         `json:"name" yaml:"name"`
	ChainID           string         `json:"chain_id" yaml:"chain_id"`
	Description       string         `json:"description" yaml:"description"`
	Tags              []string       `json:"tags" yaml:"tags"`
	Author            sdk.AccAddress `json:"author" yaml:"author"`
	AuthorDescription string         `json:"author_description" yaml:"author_description"`
	IDLContent        string         `json:"idl_content" yaml:"idl_content"`
}

// MethodProperty
type MethodProperty struct {
	ID            int16             `json:"id" yaml:"id"`
	Name          string            `json:"name" yaml:"name"`
	Description   string            `json:"description" yaml:"description"`
	OutputPrivacy OutputPrivacyEnum `json:"output_privacy" yaml:"output_privacy"`
	OutputCached  OutputCachedEnum  `json:"output_cached" yaml:"output_cached"`
}

// NewSvcDef
func NewSvcDef(name, chainID, description string, tags []string, author sdk.AccAddress, authorDescription, idlContent string) SvcDef {
	return SvcDef{
		Name:              name,
		ChainID:           chainID,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		IDLContent:        idlContent,
	}
}

// OutputPrivacyEnum
type OutputPrivacyEnum byte

const (
	NoPrivacy        OutputPrivacyEnum = 0x01
	PubKeyEncryption OutputPrivacyEnum = 0x02
)

// OutputCachedEnum
type OutputCachedEnum byte

const (
	OffChainCached OutputCachedEnum = 0x01
	NoCached       OutputCachedEnum = 0x02
)

// MessagingType
type MessagingType byte

const (
	Unicast   MessagingType = 0x01
	Multicast MessagingType = 0x02
)

// MessagingTypeFromString converts string to messagingType byte, Returns ff if invalid.
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
	return mt == Multicast || mt == Unicast
}

// Format for Printf / Sprintf, returns bech32 when using %s
func (mt MessagingType) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(fmt.Sprintf("%s", mt.String())))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%v", byte(mt))))
	}
}

// String converts MessagingType byte to String
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

// MarshalJSON marshals MessagingType to JSON using string
func (mt MessagingType) MarshalJSON() ([]byte, error) {
	return json.Marshal(mt.String())
}

// UnmarshalJSON unmarshals MessagingType from JSON assuming Bech32 encoding
func (mt *MessagingType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz2, err := MessagingTypeFromString(s)
	if err != nil {
		return err
	}
	*mt = bz2
	return nil
}

// OutputCachedEnumFromString convert string to outputCachedEnum byte, returns ff if invalid.
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
	return oe == OffChainCached || oe == NoCached
}

// Format for Printf / Sprintf, returns bech32 when using %s
func (oe OutputCachedEnum) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(fmt.Sprintf("%s", oe.String())))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%v", byte(oe))))
	}
}

// String convert OutputCachedEnum byte to string
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

// MarshalJSON marshals OutputCachedEnum to JSON using string
func (oe OutputCachedEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(oe.String())
}

// UnmarshalJSON unmarshals OutputCachedEnum from JSON assuming Bech32 encoding
func (oe *OutputCachedEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz2, err := OutputCachedEnumFromString(s)
	if err != nil {
		return err
	}
	*oe = bz2
	return nil
}

// OutputPrivacyEnumFromString convert string to outputPrivacyEnum byte, returns ff if invalid.
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
	return oe == NoPrivacy || oe == PubKeyEncryption
}

// For Printf / Sprintf, returns bech32 when using %s
func (oe OutputPrivacyEnum) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		_, _ = s.Write([]byte(fmt.Sprintf("%s", oe.String())))
	default:
		_, _ = s.Write([]byte(fmt.Sprintf("%v", byte(oe))))
	}
}

// String convert OutputCachedEnum byte to string
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

// MarshalJSON marshals OutputPrivacyEnum to JSON using string
func (oe OutputPrivacyEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(oe.String())
}

// UnmarshalJSON unmarshals OutputPrivacyEnum from JSON assuming Bech32 encoding
func (oe *OutputPrivacyEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil
	}
	bz2, err := OutputPrivacyEnumFromString(s)
	if err != nil {
		return err
	}
	*oe = bz2
	return nil
}
