package rand

import sdk "github.com/cosmos/cosmos-sdk/types"

type GenesisState struct {
	PendingRandRequests map[string][]Request `json:"pending_rand_requests"` // pending rand requests: height->[]Request
}

type Request struct {
	Height   int64          `json:"height"`   // the height of the block in which the request tx is included
	Consumer sdk.AccAddress `json:"consumer"` // the request address
	TxHash   []byte         `json:"txhash"`   // the request tx hash
}
