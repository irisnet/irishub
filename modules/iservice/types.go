package iservice

type OutputPrivacyEnum byte
type OutputCachedEnum byte
type BroadcastEnum byte

const (
	NoPrivacy        OutputPrivacyEnum   = 1
	PubKeyEncryption OutputPrivacyEnum   = 2
	OffChainCached   OutputCachedEnum    = 3
	NoCached         OutputCachedEnum    = 4
	Broadcast        BroadcastEnum = 5
	Unicast          BroadcastEnum = 6
)

type MethodProperty struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Input         string              `json:"input"`
	Output        string              `json:"output"`
	Error         string              `json:"error"`
	OutputPrivacy OutputPrivacyEnum   `json:"outputprivacy"`
	OutputCached  OutputCachedEnum    `json:"outputcached"`
}
