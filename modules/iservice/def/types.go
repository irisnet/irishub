package def

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	Unicast          = "Unicast"
	Multicast        = "Multicast"
	NoPrivacy        = "NoPrivacy"
	PubKeyEncryption = "PubKeyEncryption"
)

var (
	SvcDefKeyPrefix = []byte{0x0001} // prefix for each key to a service definition
)

var _ sdk.Msg = SvcDefMsg{}

//_________________________________________________________________________
//服务定义
type SvcDefMsg struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        string   `json:"tags"` // ('k1'='v1','k2'='v2')
	Creator     string   `json:"creator"`
	ChainID     string   `json:"chainID"`
	Messaging   string   `json:"messaging"`
	Methods     []Method `json:"methods"`
}

type Method struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Input         string `json:"input"`
	Output        string `json:"output"`
	Error         string `json:"error"`
	OutputPrivacy string `json:"outputPrivacy"`
}

// nolint
func (msg SvcDefMsg) Type() string                            { return "iservice" }
func (msg SvcDefMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg SvcDefMsg) GetSigners() []sdk.Address {
	address, _ := sdk.GetAddress(msg.Creator)
	return []sdk.Address{address}
}
func (msg SvcDefMsg) String() string {
	return ""
}
func (msg SvcDefMsg) GetSignBytes() []byte {
	cdc := wire.NewCodec()
	bz, err := cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	return bz
}
func (tx SvcDefMsg) ValidateBasic() sdk.Error {
	if tx.Name == "" {
		return NotEmpty("Invalid service name")
	}
	if tx.Description == "" {
		return NotEmpty("Invalid service description")
	}
	if tx.ChainID == "" {
		return NotEmpty("Invalid chainID")
	}
	if tx.Messaging != Unicast && tx.Messaging != Multicast {
		return NotEmpty("Invalid Messaging type ,only {unicast,multicast}")
	}
	if len(tx.Methods) < 0 {
		return NotEmpty("Invalid svc Methods")
	}

	var mthId = make(map[int]bool)
	var mthNm = make(map[string]bool)

	for _, method := range tx.Methods {
		if method.ID < 0 {
			return NotEmpty("Invalid Methods Id")
		}
		if mthId[method.ID] {
			return NotEmpty("Methods Id repeat")
		}
		mthId[method.ID] = true

		if mthNm[method.Name] {
			return NotEmpty("Methods Name repeat")
		}
		if method.Name == "" {
			return NotEmpty("Invalid Methods Name")
		}
		mthNm[method.Name] = true

		if method.Description == "" {
			return NotEmpty("Invalid Methods description")
		}
		//if method.Input == "" {
		//	return NotEmpty("Invalid service description")
		//}
		if method.Output == "" {
			return NotEmpty("Invalid Methods Output")
		}
		if method.OutputPrivacy != NoPrivacy && method.OutputPrivacy != PubKeyEncryption {
			return NotEmpty("Invalid OutputPrivacy only {outputPrivacy ,pubKeyEncryption}")
		}
	}
	return nil
}

// GetSvcDefKey - get the key for a service definition
func GetSvcDefKey(name string) []byte {
	cdc := wire.NewCodec()
	bz, err := cdc.MarshalBinary(name)
	if err != nil {
		panic(err)
	}
	return append(SvcDefKeyPrefix, bz...)
}
