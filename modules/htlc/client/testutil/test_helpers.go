package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	htlccli "github.com/irisnet/irismod/modules/htlc/client/cli"
)

// MsgRedelegateExec creates a redelegate message.
func CreateHTLCExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, htlccli.GetCmdCreateHTLC(), args)
}

func ClaimHTLCExec(clientCtx client.Context, from string, id string, secret string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		id,
		secret,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, htlccli.GetCmdClaimHTLC(), args)
}

func QueryHTLCExec(clientCtx client.Context, id string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		id,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, htlccli.GetCmdQueryHTLC(), args)
}
