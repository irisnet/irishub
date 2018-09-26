package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSvcDef struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	IDLContent        string         `json:"methods"`
	Broadcast         BroadcastEnum  `json:"broadcast"`
}
