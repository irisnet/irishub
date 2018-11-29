package serviceparams

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/iparam"
	"fmt"
	"encoding/json"
	"github.com/irisnet/irishub/codec"
)

var MaxRequestTimeoutParameter MaxRequestTimeoutParam

var _ iparam.SignalParameter = (*MaxRequestTimeoutParam)(nil)

type MaxRequestTimeoutParam struct {
	Value      int64
	paramSpace params.Subspace
}

func (param *MaxRequestTimeoutParam) InitGenesis(genesisState interface{}) {
	param.Value = genesisState.(int64)
}

func (param *MaxRequestTimeoutParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *MaxRequestTimeoutParam) GetStoreKey() []byte {
	return []byte("serviceMaxRequestTimeout")
}

func (param *MaxRequestTimeoutParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *MaxRequestTimeoutParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

func (param *MaxRequestTimeoutParam) ToJson(jsonStr string) string {
	var jsonBytes []byte

	if len(jsonStr) == 0 {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}
	return string(jsonBytes)
}

func (param *MaxRequestTimeoutParam) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *MaxRequestTimeoutParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *MaxRequestTimeoutParam) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		if param.Value <= 0{
			return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidMaxRequestTimeout, fmt.Sprintf("Invalid MaxRequestTimeout [%d] should be greater than 0",param.Value))
		}
		return nil

	}
	return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidMaxRequestTimeout, fmt.Sprintf("Json is not valid"))
}

var MinDepositMultipleParameter MinDepositMultipleParam
var _ iparam.SignalParameter = (*MinDepositMultipleParam)(nil)

type MinDepositMultipleParam struct {
	Value      int64
	paramSpace params.Subspace
}

func (param *MinDepositMultipleParam) InitGenesis(genesisState interface{}) {
	param.Value = genesisState.(int64)
}

func (param *MinDepositMultipleParam) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *MinDepositMultipleParam) GetStoreKey() []byte {
	return []byte("serviceMinDepositMultiple")
}

func (param *MinDepositMultipleParam) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *MinDepositMultipleParam) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

func (param *MinDepositMultipleParam) ToJson(jsonStr string) string {
	var jsonBytes []byte

	if len(jsonStr) == 0 {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _  = json.Marshal(param.Value)
		return string(jsonBytes)
	}
	return string(jsonBytes)
}

func (param *MinDepositMultipleParam) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *MinDepositMultipleParam) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *MinDepositMultipleParam) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		if param.Value <= 0{
			return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidMinDepositMultiple, fmt.Sprintf("Invalid MinDepositMultiple [%d] should be greater than 0",param.Value))
		}
		return nil

	}
	return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidMinDepositMultiple, fmt.Sprintf("Json is not valid"))
}
