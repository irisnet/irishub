package keeper

import (
	"encoding/binary"
	"github.com/irisnet/irishub/app/v3/oracle/internal/types"
)

var (
	// Keys for store prefixes
	separator                 = []byte("/")
	PrefixFeedKey             = []byte{0x01}
	PrefixReqCtxIdKey         = []byte{0x02}
	PrefixFeedValueKey        = []byte{0x03}
	PrefixFeedRunningStateKey = []byte{0x04}
	PrefixFeedPauseStateKey   = []byte{0x05}
)

func GetFeedKey(feedName string) []byte {
	return append(append(PrefixFeedKey, separator...), []byte(feedName)...)
}

func GetReqCtxIDKey(requestContextID []byte) []byte {
	return append(append(PrefixReqCtxIdKey, separator...), requestContextID...)
}

func GetFeedValueKey(feedName string, batchCounter uint64) []byte {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key[:], batchCounter)
	return append(GetFeedValuePrefixKey(feedName), key...)
}

func GetFeedValuePrefixKey(feedName string) []byte {
	return append(append(PrefixFeedValueKey, []byte(feedName)...), separator...)
}

func GetFeedStateKey(feedName string, state types.RequestContextState) []byte {
	return append(append(GetFeedStatePrefixKey(state), separator...), []byte(feedName)...)
}

func GetFeedStatePrefixKey(state types.RequestContextState) []byte {
	if state == types.Running {
		return PrefixFeedRunningStateKey
	}
	return PrefixFeedPauseStateKey
}
