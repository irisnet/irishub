package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	recordcli "github.com/irisnet/irismod/modules/record/client/cli"
)

// MsgRedelegateExec creates a redelegate message.
func MsgCreateRecordExec(clientCtx client.Context, from string, digest string, digestAlgo string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		digest,
		digestAlgo,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, recordcli.GetCmdCreateRecord(), args)
}

func QueryRecordExec(clientCtx client.Context, recordID string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		recordID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, recordcli.GetCmdQueryRecord(), args)
}
