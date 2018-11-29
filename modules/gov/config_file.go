package gov

import (
	"fmt"
	"encoding/json"
	sdk "github.com/irisnet/irishub/types"
	"path"
	"github.com/irisnet/irishub/codec"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/modules/params"
)

type ParameterConfigFile struct {
	Govparams govparams.ParamSet `json:"gov"`
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
func (pd *ParameterConfigFile) WriteFile(cdc *codec.Codec, res []sdk.KVPair , pathStr string) error {
	for _, kv := range res {
		switch string(kv.Key) {
		case "Gov/govDepositProcedure":
			err := cdc.UnmarshalJSON(kv.Value, &pd.Govparams.DepositProcedure)
			if err != nil {
				return err
			}
		case "Gov/govVotingProcedure":
			err := cdc.UnmarshalJSON(kv.Value, &pd.Govparams.VotingProcedure)
			if err != nil {
				return err
			}
		case "Gov/govTallyingProcedure":
			err := cdc.UnmarshalJSON(kv.Value, &pd.Govparams.TallyingProcedure)
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
	case "Gov/govDepositProcedure":
		jsonBytes, err = json.Marshal(pd.Govparams.DepositProcedure)
	case "Gov/govVotingProcedure":
		jsonBytes, err = json.Marshal(pd.Govparams.VotingProcedure)
	case "Gov/govTallyingProcedure":
		jsonBytes, err = json.Marshal(pd.Govparams.TallyingProcedure)
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
