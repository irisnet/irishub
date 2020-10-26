package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	tokencli "github.com/irisnet/irismod/modules/token/client/cli"
)

func IssueTokenExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdIssueToken(), args)
}
func EditTokenExec(clientCtx client.Context, from string, symbol string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		symbol,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdEditToken(), args)
}
func MintTokenExec(clientCtx client.Context, from string, symbol string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		symbol,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdMintToken(), args)
}
func TransferTokenOwnerExec(clientCtx client.Context, from string, symbol string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		symbol,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdTransferTokenOwner(), args)
}

func QueryTokenExec(clientCtx client.Context, denom string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryToken(), args)
}
func QueryTokensExec(clientCtx client.Context, owner string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		owner,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryTokens(), args)
}

func QueryFeeExec(clientCtx client.Context, symbol string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		symbol,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryFee(), args)
}

func QueryParamsExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryParams(), args)
}
