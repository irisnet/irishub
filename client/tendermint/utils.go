package tendermint

import (
	sdk "github.com/irisnet/irishub/types"
)

type ReadableTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func MakeTagsHumanReadable(tags sdk.Tags) []ReadableTag {
	readableTags := make([]ReadableTag, len(tags))
	for i, kv := range tags {
		readableTags[i] = ReadableTag{
			Key:   string(kv.Key),
			Value: string(kv.Value),
		}
	}
	return readableTags
}
