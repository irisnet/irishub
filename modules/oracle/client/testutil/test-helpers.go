package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	oraclecli "github.com/irisnet/irismod/modules/oracle/client/cli"
)

// MsgRedelegateExec creates a redelegate message.
func CreateFeedExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, oraclecli.GetCmdCreateFeed(), args)
}

func EditFeedExec(clientCtx client.Context, from string, feedName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		feedName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, oraclecli.GetCmdEditFeed(), args)
}

func StartFeedExec(clientCtx client.Context, from string, feedName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		feedName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, oraclecli.GetCmdStartFeed(), args)
}

func PauseFeedExec(clientCtx client.Context, from string, feedName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		feedName,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, oraclecli.GetCmdPauseFeed(), args)
}

func QueryFeedExec(clientCtx client.Context, feedName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		feedName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, oraclecli.GetCmdQueryFeed(), args)
}

func QueryFeedsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, oraclecli.GetCmdQueryFeeds(), args)
}

func QueryFeedValueExec(clientCtx client.Context, feedName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		feedName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, oraclecli.GetCmdQueryFeedValue(), args)
}
