package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	farmcli "github.com/irisnet/irismod/modules/farm/client/cli"
)

// CreateFarmPoolExec creates a redelegate message.
func CreateFarmPoolExec(clientCtx client.Context,
	creator string,
	poolName string,
	extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		poolName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdCreateFarmPool(), args)
}

func QueryFarmPoolsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdQueryFarmPools(), args)
}

func QueryFarmPoolExec(clientCtx client.Context, poolName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		poolName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdQueryFarmPool(), args)
}

// AppendRewardExec creates a redelegate message.
func AppendRewardExec(clientCtx client.Context,
	creator,
	poolName string,
	extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		poolName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdAdjustPool(), args)
}

// StakeExec creates a redelegate message.
func StakeExec(clientCtx client.Context,
	creator,
	poolName,
	lpToken string,
	extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		poolName,
		lpToken,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdStake(), args)
}

// UnstakeExec creates a redelegate message.
func UnstakeExec(clientCtx client.Context,
	creator,
	poolName,
	lpToken string,
	extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		poolName,
		lpToken,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdUnstake(), args)
}

// HarvestExec creates a redelegate message.
func HarvestExec(clientCtx client.Context,
	creator,
	poolName string,
	extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		poolName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdHarvest(), args)
}

// DestroyExec creates a redelegate message.
func DestroyExec(clientCtx client.Context,
	creator,
	poolName string,
	extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		poolName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, creator),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdDestroyFarmPool(), args)
}

// QueryFarmerExec creates a redelegate message.
func QueryFarmerExec(clientCtx client.Context,
	creator string,
	extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		creator,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, farmcli.GetCmdQueryFarmer(), args)
}
