package keeper

var (
	// Keys for store prefixes
	separator     = []byte("/")
	PrefixFeedKey = []byte{0x01}
)

func GetFeedKey(feedName string) []byte {
	return append(append(PrefixFeedKey, separator...), []byte(feedName)...)
}
