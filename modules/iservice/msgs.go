package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/iservice/idl"
)

type MsgSvcDef struct {
	Name              string         `json:"name"`
	ChainId           string         `json:"chain_id"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	IDLContent        string         `json:"idl_content"`
	Broadcast         BroadcastEnum  `json:"broadcast"`
}

func NewMsgSvcDef(name, chainId, description string, tags []string, author sdk.AccAddress, authorDescription, idlContent string, broadcast BroadcastEnum) MsgSvcDef {
	return MsgSvcDef{
		Name:              name,
		ChainId:           chainId,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		IDLContent:        idlContent,
		Broadcast:         broadcast,
	}
}

func (msg MsgSvcDef) Type() string {
	return "iservice"
}

func (msg MsgSvcDef) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgSvcDef) ValidateBasic() sdk.Error {
	if !idl.ValidIDL(msg.IDLContent) {
		return NewError(DefaultCodespace, CodeInvalidIDL, "")
	}
	return nil
}

func (msg MsgSvcDef) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}
