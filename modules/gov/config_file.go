package gov

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/upgrade/params"
	sdk "github.com/irisnet/irishub/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"path"
	"github.com/irisnet/irishub/modules/service/params"
)

type ParameterConfigFile struct {
	UpgradeParams upgradeparams.Params `json:"upgrade"`
	ServiceParams serviceparams.Params `json:"service"`
}

func (pd *ParameterConfigFile) ReadFile(cdc *codec.Codec, pathStr string) error {
	pathStr = path.Join(pathStr, "params.json")

	jsonBytes, err := cmn.ReadFile(pathStr)

	fmt.Println("Open ", pathStr)

	if err != nil {
		return err
	}

	err = cdc.UnmarshalJSON(jsonBytes, &pd)
	return err
}
func (pd *ParameterConfigFile) WriteFile(cdc *codec.Codec, res []sdk.KVPair, pathStr string) error {
	for _, kv := range res {
		switch string(kv.Key) {
		case "Gov/"+upgradeparams.UpgradeParamsKey:
			err := cdc.UnmarshalJSON(kv.Value, &pd.UpgradeParams)
			if err != nil {
				return err
			}
		case "Gov/"+serviceparams.ServiceParamsKey:
			err := cdc.UnmarshalJSON(kv.Value, &pd.ServiceParams)
			if err != nil {
				return err
			}
		}
	}

	output, err := cdc.MarshalJSONIndent(pd, "", "  ")

	if err != nil {
		return err
	}

	pathStr = path.Join(pathStr, "params.json")
	err = cmn.WriteFile(pathStr, output, 0644)
	if err != nil {

		return err
	}

	fmt.Println("Save the parameter config file in ", pathStr)
	return nil
}

func (pd *ParameterConfigFile) GetParamFromKey(keyStr string, opStr string) (Param, error) {
	var param Param
	var err error
	var jsonBytes []byte

	if len(keyStr) == 0 {
		return param, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("Key can't be empty!"))
	}

	switch keyStr {
	case "Gov/"+upgradeparams.UpgradeParamsKey:
		jsonBytes, err = json.Marshal(pd.UpgradeParams)
	case "Gov/"+serviceparams.ServiceParamsKey:
		jsonBytes, err = json.Marshal(pd.ServiceParams)
	default:
		return param, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf(keyStr+" is not found"))
	}

	if err != nil {
		return param, err
	}
	param.Value = string(jsonBytes)
	param.Key = keyStr
	param.Op = opStr

	jsonBytes, _ = json.MarshalIndent(param, "", " ")

	return param, err
}
