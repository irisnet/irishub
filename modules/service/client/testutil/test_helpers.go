package testutil

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	servicecli "github.com/irisnet/irismod/modules/service/client/cli"
)

func DefineServiceExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdDefineService(), args)
}

func BindServiceExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdBindService(), args)
}

func UpdateBindingExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdUpdateServiceBinding(), args)
}

func RefundDepositExec(clientCtx client.Context, serviceName, provider, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdRefundServiceDeposit(), args)
}

func DisableServiceExec(clientCtx client.Context, serviceName, provider, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdDisableServiceBinding(), args)
}

func EnableServiceExec(clientCtx client.Context, serviceName, provider, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdEnableServiceBinding(), args)
}

func CallServiceExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdCallService(), args)
}

func RespondServiceExec(clientCtx client.Context, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdRespondService(), args)
}

func SetWithdrawAddrExec(clientCtx client.Context, withdrawalAddress, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		withdrawalAddress,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdSetWithdrawAddr(), args)
}

func WithdrawEarnedFeesExec(clientCtx client.Context, provider, from string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdWithdrawEarnedFees(), args)
}

func QueryServiceDefinitionExec(clientCtx client.Context, serviceName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		serviceName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdQueryServiceDefinition(), args)
}

func QueryServiceBindingExec(clientCtx client.Context, serviceName, provider string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdQueryServiceBinding(), args)
}

func QueryServiceBindingsExec(clientCtx client.Context, serviceName string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		serviceName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdQueryServiceBindings(), args)
}

func QueryServiceRequestsExec(clientCtx client.Context, serviceName, provider string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdQueryServiceRequests(), args)
}

func QueryEarnedFeesExec(clientCtx client.Context, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, servicecli.GetCmdQueryEarnedFees(), args)
}
