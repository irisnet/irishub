package service

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	servicecli "mods.irisnet.org/modules/service/client/cli"
	servicetypes "mods.irisnet.org/modules/service/types"
	"mods.irisnet.org/simapp"
)

func DefineServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdDefineService(), args)
}

func BindServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdBindService(), args)
}

func UpdateBindingExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdUpdateServiceBinding(), args)
}

func RefundDepositExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdRefundServiceDeposit(), args)
}

func DisableServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdDisableServiceBinding(), args)
}

func EnableServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdEnableServiceBinding(), args)
}

func CallServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdCallService(), args)
}

func RespondServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdRespondService(), args)
}

func SetWithdrawAddrExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	withdrawalAddress,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		withdrawalAddress,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdSetWithdrawAddr(), args)
}

func WithdrawEarnedFeesExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdWithdrawEarnedFees(), args)
}

func QueryServiceDefinitionExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName string,
	extraArgs ...string,
) *servicetypes.ServiceDefinition {
	args := []string{
		serviceName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.ServiceDefinition{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceDefinition(), args, response)
	return response
}

func QueryServiceBindingExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider string,
	extraArgs ...string,
) *servicetypes.ServiceBinding {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.ServiceBinding{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceBinding(), args, response)
	return response
}

func QueryServiceBindingsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName string,
	extraArgs ...string,
) *servicetypes.QueryBindingsResponse {
	args := []string{
		serviceName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.QueryBindingsResponse{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceBindings(), args, response)
	return response
}

func QueryServiceRequestsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider string,
	extraArgs ...string,
) *servicetypes.QueryRequestsResponse {
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.QueryRequestsResponse{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceRequests(), args, response)
	return response
}

func QueryServiceRequestsByReqCtx(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	requestContextID,
	batchCounter string,
	extraArgs ...string,
) *servicetypes.QueryRequestsResponse {
	args := []string{
		requestContextID,
		batchCounter,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.QueryRequestsResponse{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceRequests(), args, response)
	return response
}

func QueryEarnedFeesExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	extraArgs ...string,
) *servicetypes.QueryEarnedFeesResponse {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.QueryEarnedFeesResponse{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryEarnedFees(), args, response)
	return response
}

func QueryRequestContextExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	contextId string,
	extraArgs ...string,
) *servicetypes.RequestContext {
	args := []string{
		contextId,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.RequestContext{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryRequestContext(), args, response)
	return response
}

func QueryServiceRequestExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	requestId string,
	extraArgs ...string,
) *servicetypes.Request {
	args := []string{
		requestId,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.Request{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceRequest(), args, response)
	return response
}

func QueryServiceResponseExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	requestId string,
	extraArgs ...string,
) *servicetypes.Response {
	args := []string{
		requestId,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.Response{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceResponse(), args, response)
	return response
}
