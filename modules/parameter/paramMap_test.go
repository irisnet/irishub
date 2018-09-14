package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/irisnet/irishub/modules/gov"
	govParam "github.com/irisnet/irishub/modules/gov/params"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegisterParamMapping(t *testing.T) {

	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	p1 := gov.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 10)},
		MaxDepositPeriod: 1440}

	p2 := gov.DepositProcedure{
		MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 30)},
		MaxDepositPeriod: 1440}

	govParam.DepositProcedureParameter.SetReadWriter(paramKeeper.Setter())
	RegisterGovParamMapping(&govParam.DepositProcedureParameter)
	InitGenesisParameter(&govParam.DepositProcedureParameter, ctx, nil)

	require.Equal(t, paramMapping[govParam.DepositProcedureParameter.GetStoreKey()].ToJson(), "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"10\"}],\"max_deposit_period\":1440}")
	require.Equal(t, p1, govParam.DepositProcedureParameter.Value)

	paramMapping[govParam.DepositProcedureParameter.GetStoreKey()].Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris\",\"amount\":\"30\"}],\"max_deposit_period\":1440}")
	govParam.DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, govParam.DepositProcedureParameter.Value)
}
