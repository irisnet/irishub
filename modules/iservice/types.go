package iservice

import (
	"github.com/pkg/errors"
	"fmt"
	"encoding/json"
)

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

type BroadcastEnum byte

const (
	Broadcast BroadcastEnum = 0x01
	Unicast   BroadcastEnum = 0x02
)

type MethodProperty struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	OutputPrivacy OutputPrivacyEnum `json:"output_privacy"`
	OutputCached  OutputCachedEnum  `json:"output_cached"`
}

// String to broadcastEnum byte, Returns ff if invalid.
func BroadcastEnumFromString(str string) (BroadcastEnum, error) {
	switch str {
	case "Broadcast":
		return Broadcast, nil
	case "Unicast":
		return Unicast, nil
	default:
		return BroadcastEnum(0xff), errors.Errorf("'%s' is not a valid broadcastEnum type", str)
	}
}

// is defined BroadcastEnum?
func validBroadcastEnum(be BroadcastEnum) bool {
	if be == Broadcast ||
		be == Unicast {
		return true
	}
	return false
}

// For Printf / Sprintf, returns bech32 when using %s
func (be BroadcastEnum) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", be.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(be))))
	}
}

// Turns BroadcastEnum byte to String
func (be BroadcastEnum) String() string {
	switch be {
	case Broadcast:
		return "Broadcast"
	case Unicast:
		return "unicast"
	default:
		return ""
	}
}

// Marshals to JSON using string
func (be BroadcastEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(be.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (be *BroadcastEnum) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := BroadcastEnumFromString(s)
	if err != nil {
		return err
	}
	*be = bz2
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

