package upgradeparams

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

const UpgradeParamsKey = "upgradeParams"

var UpgradeParameter UpgradeParams

var _ params.GovParameter = (*UpgradeParams)(nil)

// Procedure around Deposits for governance
type Params struct {
	Threshold sdk.Dec `json:"threshold"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
}

type UpgradeParams struct {
	Value      Params
	paramSpace params.Subspace
}

func NewUpgradeParams() Params {
	return Params{
		Threshold: sdk.NewDecWithPrec(9, 1),
	}
}

func (param *UpgradeParams) GetValueFromRawData(cdc *codec.Codec, res []byte) interface{} {
	cdc.UnmarshalJSON(res, &param.Value)
	return param.Value
}

func (param *UpgradeParams) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(Params); ok {
		param.Value = value
	} else {
		panic("The " + UpgradeParamsKey + " in GenesisState is empty. ")
	}
}

func (param *UpgradeParams) SetReadWriter(paramSpace params.Subspace) {
	param.paramSpace = paramSpace
}

func (param *UpgradeParams) GetStoreKey() []byte {
	return []byte(UpgradeParamsKey)
}

func (param *UpgradeParams) SaveValue(ctx sdk.Context) {
	param.paramSpace.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *UpgradeParams) LoadValue(ctx sdk.Context) bool {
	if param.paramSpace.Has(ctx, param.GetStoreKey()) == false {
		return false
	}
	param.paramSpace.Get(ctx, param.GetStoreKey(), &param.Value)
	return true
}

func (param *UpgradeParams) ToJson(jsonStr string) string {

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

func (param *UpgradeParams) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *UpgradeParams) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		if param.Value.Threshold.LT(sdk.NewDecWithPrec(67, 2)) || param.Value.Threshold.GT(sdk.NewDec(1)) {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidUpgradeParams, fmt.Sprintf("Invalid Upgrade Threshold( "+param.Value.Threshold.String()+" ) should be between 0.67 and 1"))
		}
		return nil

	}
	return sdk.NewError(params.DefaultCodespace, params.CodeInvalidMinDeposit, fmt.Sprintf("Json is not valid"))
}
