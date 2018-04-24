package bind

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var (
	SvcBindKeyPrefix = []byte{0x0002} // prefix for each key to a service definition
)

//服务绑定
type SvcBind struct {
	DefinitionHash    []byte `json:"definitionHash"`
	ChainID           string `json:"chainID"`
	ProviderAddress   []byte `json:"providerAddress"`
	BindingType       string `json:"bindingType"`
	ProviderDeposit   int64  `json:"providerDeposit"`
	ServicePricing    string `json:"servicePricing"`
	ServiceLevel      string `json:"serviceLevel"`
	BindingExpiration int64  `json:"bindingExpiration"`
	IsValid           bool   `json:"isValid"`
}

var _ sdk.Msg = SvcBindMsg{}

type SvcBindMsg struct {
	Sender sdk.Address
	SvcBind
}

// nolint
func (msg SvcBindMsg) Type() string                            { return "iservice" }
func (msg SvcBindMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg SvcBindMsg) GetSigners() []sdk.Address               { return []sdk.Address{msg.Sender} }
func (msg SvcBindMsg) String() string {
	return ""
}
func (msg SvcBindMsg) GetSignBytes() []byte {
	cdc := wire.NewCodec()
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}
func (tx SvcBindMsg) ValidateBasic() sdk.Error {
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
