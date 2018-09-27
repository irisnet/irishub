package iservice

type OutputPrivacyEnum byte
const (
	NoPrivacy        OutputPrivacyEnum = 0x01
	PubKeyEncryption OutputPrivacyEnum = 0x02
)

type OutputCachedEnum byte
const (
	OffChainCached   OutputCachedEnum  = 0x01
	NoCached         OutputCachedEnum  = 0x02
)

type BroadcastEnum byte
const (
	Broadcast        BroadcastEnum     = 0x01
	Unicast          BroadcastEnum     = 0x02
)

type MethodProperty struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	OutputPrivacy OutputPrivacyEnum `json:"output_privacy"`
	OutputCached  OutputCachedEnum  `json:"output_cached"`
}
