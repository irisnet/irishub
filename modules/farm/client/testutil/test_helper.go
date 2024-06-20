package testutil

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cometbft/cometbft/libs/cli"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/client/flags"

// 	"github.com/irisnet/irismod/simapp"
// 	farmcli "github.com/irisnet/irismod/farm/client/cli"
// 	farmtypes "github.com/irisnet/irismod/farm/types"
// )

// // CreateFarmPoolExec creates a redelegate message.
// func CreateFarmPoolExec(t *testing.T, network simapp.Network, clientCtx client.Context,
// 	creator string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, farmcli.GetCmdCreateFarmPool(), args)
// }

// func QueryFarmPoolsExec(
// 	t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	extraArgs ...string,
// ) *farmtypes.QueryFarmPoolsResponse {
// 	args := []string{
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)

// 	response := &farmtypes.QueryFarmPoolsResponse{}
// 	network.ExecQueryCmd(t, clientCtx, farmcli.GetCmdQueryFarmPools(), args, response)
// 	return response
// }

// func QueryFarmPoolExec(
// 	t *testing.T,
// 	network simapp.Network,
// 	clientCtx client.Context,
// 	poolId string,
// 	extraArgs ...string,
// ) *farmtypes.QueryFarmPoolResponse {
// 	args := []string{
// 		poolId,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)
// 	response := &farmtypes.QueryFarmPoolResponse{}
// 	network.ExecQueryCmd(t, clientCtx, farmcli.GetCmdQueryFarmPool(), args, response)
// 	return response
// }

// // AppendRewardExec creates a redelegate message.
// func AppendRewardExec(t *testing.T, network simapp.Network, clientCtx client.Context,
// 	creator,
// 	poolId string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		poolId,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, farmcli.GetCmdAdjustPool(), args)
// }

// // StakeExec creates a redelegate message.
// func StakeExec(t *testing.T, network simapp.Network, clientCtx client.Context,
// 	creator,
// 	poolId,
// 	lpToken string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		poolId,
// 		lpToken,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, farmcli.GetCmdStake(), args)
// }

// // UnstakeExec creates a redelegate message.
// func UnstakeExec(t *testing.T, network simapp.Network, clientCtx client.Context,
// 	creator,
// 	poolId,
// 	lpToken string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		poolId,
// 		lpToken,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, farmcli.GetCmdUnstake(), args)
// }

// // HarvestExec creates a redelegate message.
// func HarvestExec(t *testing.T, network simapp.Network, clientCtx client.Context,
// 	creator,
// 	poolId string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		poolId,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, farmcli.GetCmdHarvest(), args)
// }

// // DestroyExec creates a redelegate message.
// func DestroyExec(t *testing.T, network simapp.Network, clientCtx client.Context,
// 	creator,
// 	poolId string,
// 	extraArgs ...string) *simapp.ResponseTx {
// 	args := []string{
// 		poolId,
// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
// 	}
// 	args = append(args, extraArgs...)
// 	return network.ExecTxCmdWithResult(t, clientCtx, farmcli.GetCmdDestroyFarmPool(), args)
// }

// // QueryFarmerExec creates a redelegate message.
// func QueryFarmerExec(t *testing.T, network simapp.Network, clientCtx client.Context,
// 	creator string,
// 	extraArgs ...string) *farmtypes.QueryFarmerResponse {
// 	args := []string{
// 		creator,
// 		fmt.Sprintf("--%s=json", cli.OutputFlag),
// 	}
// 	args = append(args, extraArgs...)
// 	response := &farmtypes.QueryFarmerResponse{}
// 	network.ExecQueryCmd(t, clientCtx, farmcli.GetCmdQueryFarmer(), args, response)
// 	return response
// }
