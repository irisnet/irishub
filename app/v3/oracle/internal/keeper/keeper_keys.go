package keeper

import "encoding/binary"

var (
	// Keys for store prefixes
	separator           = []byte("/")
	PrefixFeedKey       = []byte{0x01}
	PrefixReqCtxIdKey   = []byte{0x02}
	PrefixFeedResultKey = []byte{0x03}
)

func GetFeedKey(feedName string) []byte {
	return append(append(PrefixFeedKey, separator...), []byte(feedName)...)
}

func GetReqCtxIDKey(requestContextID []byte) []byte {
	return append(append(PrefixReqCtxIdKey, separator...), requestContextID...)
}

func GetFeedResultKey(feedName string, batchCounter uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key[:], batchCounter)
	return append(GetFeedResultPrefixKey(feedName), key...)
}

func GetFeedResultPrefixKey(feedName string) []byte {
	return append(append(PrefixFeedResultKey, []byte(feedName)...), separator...)
}
