package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	randomcli "github.com/irisnet/irismod/modules/random/client/cli"
)

// MsgRedelegateExec creates a redelegate message.
func RequestRandomExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, randomcli.GetCmdRequestRandom(), args)
}

func QueryRandomExec(clientCtx client.Context, requestID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		requestID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, randomcli.GetCmdQueryRandom(), args)
}

func QueryRandomRequestQueueExec(clientCtx client.Context, genHeight string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		genHeight,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, randomcli.GetCmdQueryRandomRequestQueue(), args)
}
