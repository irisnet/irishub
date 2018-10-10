package iservice

import "github.com/pkg/errors"

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
	ID            int64             `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	OutputPrivacy OutputPrivacyEnum `json:"output_privacy"`
	OutputCached  OutputCachedEnum  `json:"output_cached"`
}

// String to broadcastEnum byte, Returns ff if invalid.
func BroadcastEnumFromString(str string) (BroadcastEnum, error) {
	switch str {
	case "broadcast":
		return Broadcast, nil
	case "unicast":
		return Unicast, nil
	default:
		return BroadcastEnum(0xff), errors.Errorf("'%s' is not a valid broadcast type", str)
	}
}
