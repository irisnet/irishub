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

// DefineServiceExec defines a service execution.
//
// Parameters:
// - t: *testing.T
// - network: simapp.Network
// - clientCtx: client.Context
// - from: string
// - extraArgs: ...string
// Returns *simapp.ResponseTx
func DefineServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdDefineService(), args)
}

// BindServiceExec executes the command to bind a service.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - from: string - the address of the user binding the service
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func BindServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdBindService(), args)
}

// UpdateBindingExec executes the command to update a service binding.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - from: string - the address of the user updating the service binding
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func UpdateBindingExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdUpdateServiceBinding(), args)
}

// RefundDepositExec executes the command to refund the deposit for a service.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - serviceName: string - the name of the service
// - provider: string - the provider of the service
// - from: string - the address of the user refunding the deposit
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func RefundDepositExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdRefundServiceDeposit(), args)
}

// DisableServiceExec executes the command to disable a service binding.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - serviceName: string - the name of the service
// - provider: string - the provider of the service
// - from: string - the address of the user disabling the service
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func DisableServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdDisableServiceBinding(), args)
}

// EnableServiceExec executes the command to enable a service binding.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - serviceName: string - the name of the service
// - provider: string - the provider of the service
// - from: string - the address of the user enabling the service
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func EnableServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		serviceName,
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdEnableServiceBinding(), args)
}

// CallServiceExec executes the command to call a service.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - from: string - the address of the user calling the service
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func CallServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdCallService(), args)
}

// RespondServiceExec executes the command to respond to a service request.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - from: string - the address of the user responding to the service request
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func RespondServiceExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdRespondService(), args)
}

// SetWithdrawAddrExec executes the command to set a withdrawal address.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - withdrawalAddress: string - the address to set for withdrawal
// - from: string - the address of the user setting the withdrawal address
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func SetWithdrawAddrExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	withdrawalAddress,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		withdrawalAddress,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdSetWithdrawAddr(), args)
}

// WithdrawEarnedFeesExec executes the command to withdraw earned fees.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - provider: string - the provider of the service
// - from: string - the address of the user withdrawing the fees
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *simapp.ResponseTx - the response transaction object
func WithdrawEarnedFeesExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	provider,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		provider,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, servicecli.GetCmdWithdrawEarnedFees(), args)
}

// QueryServiceDefinitionExec executes a query to retrieve a service definition from the network.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - serviceName: string - the name of the service
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *servicetypes.ServiceDefinition - the response object containing the service definition
func QueryServiceDefinitionExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName string,
	extraArgs ...string,
) *servicetypes.ServiceDefinition {
	t.Helper()
	args := []string{
		serviceName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.ServiceDefinition{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceDefinition(), args, response)
	return response
}

// QueryServiceBindingExec executes a query to retrieve a service binding from the network.
//
// Parameters:
// - t: *testing.T - the testing object.
// - network: simapp.Network - the network object.
// - clientCtx: client.Context - the client context object.
// - serviceName: string - the name of the service.
// - provider: string - the provider of the service.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *servicetypes.ServiceBinding - the response object containing the service binding.
func QueryServiceBindingExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider string,
	extraArgs ...string,
) *servicetypes.ServiceBinding {
	t.Helper()
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

// QueryServiceBindingsExec executes a query to retrieve service bindings from the network.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - serviceName: string - the name of the service
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *servicetypes.QueryBindingsResponse - the response object containing the service bindings
func QueryServiceBindingsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName string,
	extraArgs ...string,
) *servicetypes.QueryBindingsResponse {
	t.Helper()
	args := []string{
		serviceName,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.QueryBindingsResponse{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceBindings(), args, response)
	return response
}

// QueryServiceRequestsExec queries the service requests by service name and provider.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - serviceName: string - the name of the service
// - provider: string - the provider of the service
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *servicetypes.QueryRequestsResponse - the response object containing the service requests
func QueryServiceRequestsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	serviceName,
	provider string,
	extraArgs ...string,
) *servicetypes.QueryRequestsResponse {
	t.Helper()
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

// QueryServiceRequestsByReqCtx queries the service requests by request context ID.
//
// Parameters:
// - t: *testing.T
// - network: simapp.Network
// - clientCtx: client.Context
// - requestContextID: string
// - batchCounter: string
// - extraArgs: ...string
//
// Returns:
// - *servicetypes.QueryRequestsResponse
func QueryServiceRequestsByReqCtx(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	requestContextID,
	batchCounter string,
	extraArgs ...string,
) *servicetypes.QueryRequestsResponse {
	t.Helper()
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

// QueryEarnedFeesExec executes the command to query the earned fees.
//
// Parameters:
// - t: *testing.T - the testing object
// - network: simapp.Network - the network object
// - clientCtx: client.Context - the client context object
// - extraArgs: ...string - additional arguments for the command
//
// Returns:
// - *servicetypes.QueryEarnedFeesResponse - the response object containing the earned fees
func QueryEarnedFeesExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	extraArgs ...string,
) *servicetypes.QueryEarnedFeesResponse {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.QueryEarnedFeesResponse{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryEarnedFees(), args, response)
	return response
}

// QueryRequestContextExec executes a query to retrieve a request context from the network.
//
// Parameters:
// - t: *testing.T
// - network: simapp.Network
// - clientCtx: client.Context
// - contextId: string
// - extraArgs: ...string
// Returns *servicetypes.RequestContext
func QueryRequestContextExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	contextID string,
	extraArgs ...string,
) *servicetypes.RequestContext {
	t.Helper()
	args := []string{
		contextID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.RequestContext{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryRequestContext(), args, response)
	return response
}

// QueryServiceRequestExec executes a query to retrieve a service request from the network.
//
// Parameters:
// - t: *testing.T
// - network: simapp.Network
// - clientCtx: client.Context
// - requestID: string
// - extraArgs: ...string
// Returns *servicetypes.Request
func QueryServiceRequestExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	requestID string,
	extraArgs ...string,
) *servicetypes.Request {
	t.Helper()
	args := []string{
		requestID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.Request{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceRequest(), args, response)
	return response
}

// QueryServiceResponseExec executes a query to retrieve a service response from the network.
//
// Parameters:
// - t: *testing.T
// - network: simapp.Network
// - clientCtx: client.Context
// - requestID: string
// - extraArgs: ...string
// Returns *servicetypes.Response
func QueryServiceResponseExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	requestID string,
	extraArgs ...string,
) *servicetypes.Response {
	t.Helper()
	args := []string{
		requestID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &servicetypes.Response{}
	network.ExecQueryCmd(t, clientCtx, servicecli.GetCmdQueryServiceResponse(), args, response)
	return response
}
