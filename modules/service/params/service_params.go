package serviceparams

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

const ServiceParamsKey = "serviceParams"

var ServiceParameter ServiceParams

var _ params.GovParameter = (*ServiceParams)(nil)

type Params struct {
	MaxRequestTimeout  int64   `json:"max_request_timeout"`
	MinDepositMultiple int64   `json:"min_deposit_multiple"`
	ServiceFeeTax      sdk.Dec `json:"service_fee_tax"`
	SlashFraction      sdk.Dec `json:"slash_fraction"`
}

type ServiceParams struct {
	Value      Params
	paramSpace params.Subspace
}

func NewSericeParams() Params {
	return Params{
		MaxRequestTimeout:  100,
		MinDepositMultiple: 1000,
		ServiceFeeTax:      sdk.NewDecWithPrec(2, 2), //2%
		SlashFraction:      sdk.NewDecWithPrec(1, 2), //1%
	}
}

func (param *ServiceParams) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(Params); ok {
		param.Value = value
	} else {
		param.Value = NewSericeParams()
	}
}

func (param *ServiceParams) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *ServiceParams) GetStoreKey() []byte {
	return []byte(ServiceParamsKey)
}

func (param *ServiceParams) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *ServiceParams) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

func (param *ServiceParams) ToJson(jsonStr string) string {
	var jsonBytes []byte

	if len(jsonStr) == 0 {
		jsonBytes, _ = json.Marshal(param.Value)
		return string(jsonBytes)
	}

	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		jsonBytes, _ = json.Marshal(param.Value)
		return string(jsonBytes)
	}
	return string(jsonBytes)
}

func (param *ServiceParams) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *ServiceParams) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *ServiceParams) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		if param.Value.MaxRequestTimeout <= 0 {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMaxRequestTimeout, fmt.Sprintf("Invalid MaxRequestTimeout [%d] should be greater than 0", param.Value.MaxRequestTimeout))
		}
		if param.Value.MinDepositMultiple <= 0 {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDepositMultiple, fmt.Sprintf("Invalid MinDepositMultiple [%d] should be greater than 0", param.Value.MinDepositMultiple))
		}
		if param.Value.ServiceFeeTax.LTE(sdk.ZeroDec()) || param.Value.ServiceFeeTax.GTE(sdk.NewDec(1)) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidServiceFeeTax, fmt.Sprintf("Invalid ServiceFeeTax ( "+param.Value.ServiceFeeTax.String()+" ) should be between 0 and 1"))
		}
		if param.Value.SlashFraction.LTE(sdk.ZeroDec()) || param.Value.SlashFraction.GTE(sdk.NewDec(1)) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashFraction, fmt.Sprintf("Invalid SlashFraction ( "+param.Value.SlashFraction.String()+" ) should be between 0 and 1"))
		}
		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidServiceParams, fmt.Sprintf("ServiceParams Json is not valid. "))
}
