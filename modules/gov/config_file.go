package gov

import (
	"fmt"
	"github.com/irisnet/irishub/modules/iparam"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"path"
	"github.com/cosmos/cosmos-sdk/codec"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/irisnet/irishub/modules/gov/params"
)

type ParameterConfigFile struct {
	Govparams govparams.ParamSet `json:"gov"`
}

func (pd *ParameterConfigFile) ReadFile(cdc *codec.Codec, pathStr string) error {
	pathStr = path.Join(pathStr, "config/params.json")

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
		case "Gov/gov/DepositProcedure":
			cdc.MustUnmarshalBinary(kv.Value, &pd.Govparams.DepositProcedure)
		case "Gov/gov/VotingProcedure":
			cdc.MustUnmarshalBinary(kv.Value, &pd.Govparams.VotingProcedure)
		case "Gov/gov/TallyingProcedure":
			cdc.MustUnmarshalBinary(kv.Value, &pd.Govparams.TallyingProcedure)
		default:
			return sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidTallyingProcedure, fmt.Sprintf(string(kv.Key)+" is not found"))
		}
	}
	output, err := cdc.MarshalJSONIndent(pd, "", "  ")

	if err != nil {
		return err
	}

	pathStr = path.Join(pathStr, "config/params.json")
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
		return param, sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidKey, fmt.Sprintf("Key can't be empty!"))
	}

	switch keyStr {
	case "Gov/gov/DepositProcedure":
		jsonBytes, err = json.Marshal(pd.Govparams.DepositProcedure)
	case "Gov/gov/VotingProcedure":
		jsonBytes, err = json.Marshal(pd.Govparams.VotingProcedure)
	case "Gov/gov/TallyingProcedure":
		jsonBytes, err = json.Marshal(pd.Govparams.TallyingProcedure)
	default:
		return param, sdk.NewError(iparam.DefaultCodespace, iparam.CodeInvalidKey, fmt.Sprintf(keyStr+" is not found"))
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
