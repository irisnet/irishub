package call

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var (
	SvcBindKeyPrefix = []byte{0x0002} // prefix for each key to a service definition
)

//服务调用
type SvcReq struct {
	ChainID         string `json:"chainID"`
	ConsumerAddress []byte `json:"consumerAddress"`
	DefinitionHash  []byte `json:"definitionHash"`
	MethodID        int    `json:"methodID"`
	InputValue      string `json:"inputValue"`
	BindingHash     []byte `json:"bindingHash"`
	MaxServiceFee   int64  `json:"maxServiceFee"`
	Timeout         int64  `json:"timeout"`
}

var _ sdk.Msg = SvcReqMsg{}

type SvcReqMsg struct {
	Sender   sdk.Address
	Receiver sdk.Address
	SvcReq
}

// nolint
func (msg SvcReqMsg) Type() string                            { return "iservice" }
func (msg SvcReqMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg SvcReqMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg SvcReqMsg) String() string {
	return ""
}
func (msg SvcReqMsg) GetSignBytes() []byte {
	cdc := wire.NewCodec()
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}
func (tx SvcReqMsg) ValidateBasic() sdk.Error {
	return nil
}

// GetServiceDefinitionKey - get the key for a service definition
func GetSvcBindKey(name string) []byte {
	cdc := wire.NewCodec()
	bz, err := cdc.MarshalBinary(name)
	if err != nil {
		panic(err)
	}
	return append(SvcBindKeyPrefix, bz...)
}
