package token

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	tokencli "mods.irisnet.org/modules/token/client/cli"
	v1 "mods.irisnet.org/modules/token/types/v1"
	"mods.irisnet.org/simapp"
)

// IssueTokenExec executes the command to issue a token on the specified network with the given client context and sender address.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the token will be issued.
// - clientCtx: client.Context - the client context for the transaction.
// - from: string - the address of the sender.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *simapp.ResponseTx - the response transaction from executing the command.
func IssueTokenExec(t *testing.T,
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

	return network.ExecTxCmdWithResult(t, clientCtx, tokencli.GetCmdIssueToken(), args)
}

// EditTokenExec executes the command to edit a token on the specified network with the given client context and sender address.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the token will be edited.
// - clientCtx: client.Context - the client context for the transaction.
// - from: string - the address of the sender.
// - symbol: string - the symbol of the token to be edited.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *simapp.ResponseTx - the response transaction from executing the command.
func EditTokenExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	symbol string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		symbol,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, tokencli.GetCmdEditToken(), args)
}

// MintTokenExec executes the command to mint a token on the specified network with the given client context and sender address.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the token will be minted.
// - clientCtx: client.Context - the client context for the transaction.
// - from: string - the address of the sender.
// - coinStr: string - the amount and coin type to be minted.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *simapp.ResponseTx - the response transaction from executing the command.
func MintTokenExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	coinStr string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		coinStr,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, tokencli.GetCmdMintToken(), args)
}

// BurnTokenExec executes the command to burn a token on the specified network with the given client context and sender address.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the token will be burned.
// - clientCtx: client.Context - the client context for the transaction.
// - from: string - the address of the sender.
// - coinStr: string - the amount and coin type to be burned.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *simapp.ResponseTx - the response transaction from executing the command.
func BurnTokenExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	coinStr string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		coinStr,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, tokencli.GetCmdBurnToken(), args)
}

// TransferTokenOwnerExec executes the command to transfer the ownership of a token on the specified network with the given client context and sender address.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the token ownership will be transferred.
// - clientCtx: client.Context - the client context for the transaction.
// - from: string - the address of the current owner.
// - symbol: string - the symbol of the token for which ownership will be transferred.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *simapp.ResponseTx - the response transaction from executing the command.
func TransferTokenOwnerExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	symbol string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		symbol,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, tokencli.GetCmdTransferTokenOwner(), args)
}

// SwapToERC20Exec executes the command to swap a given coin to ERC20 on the specified network with the given client context and sender address.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the swap will be executed.
// - clientCtx: client.Context - the client context for the transaction.
// - from: string - the address of the sender.
// - coinStr: string - the amount and coin type to be swapped.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *simapp.ResponseTx - the response transaction from executing the command.
func SwapToERC20Exec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	coinStr string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		coinStr,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, tokencli.GetCmdSwapToErc20(), args)
}

// SwapFromERC20Exec executes the command to swap a given coin from ERC20 to native token on the specified network with the given client context and sender address.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the swap will be executed.
// - clientCtx: client.Context - the client context for the transaction.
// - from: string - the address of the sender.
// - coinStr: string - the amount and coin type to be swapped.
// - extraArgs: ...string - additional arguments for the command.
//
// Returns:
// - *simapp.ResponseTx - the response transaction from executing the command.
func SwapFromERC20Exec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	coinStr string,
	extraArgs ...string,
) *simapp.ResponseTx {
	t.Helper()
	args := []string{
		coinStr,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, tokencli.GetCmdSwapFromErc20(), args)
}

// QueryTokenExec executes a query command to retrieve information about a token.
//
// Parameters:
// - t: testing instance
// - network: simapp.Network instance
// - clientCtx: client.Context instance
// - denom: string representing the denomination of the token
// - extraArgs: variadic string arguments
//
// Returns:
// - v1.TokenI interface
func QueryTokenExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denom string,
	extraArgs ...string,
) v1.TokenI {
	t.Helper()
	args := []string{
		denom,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	respType := proto.Message(&codectypes.Any{})
	network.ExecQueryCmd(t, clientCtx, tokencli.GetCmdQueryToken(), args, respType)

	var token v1.TokenI
	err := clientCtx.InterfaceRegistry.UnpackAny(respType.(*codectypes.Any), &token)
	require.NoError(t, err, "QueryTokenExec failed")
	return token
}

// QueryTokensExec executes a query command to retrieve information about tokens owned by a specific owner.
//
// Parameters:
// - t: testing instance
// - network: simapp.Network instance
// - clientCtx: client.Context instance
// - owner: string representing the owner of the tokens
// - extraArgs: variadic string arguments
//
// Returns:
// - []v1.TokenI: a slice of v1.TokenI representing the tokens owned by the specified owner
func QueryTokensExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	owner string,
	extraArgs ...string,
) []v1.TokenI {
	t.Helper()
	args := []string{
		owner,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)
	tokens := []v1.TokenI{}
	buf, err := clitestutil.ExecTestCLICmd(clientCtx, tokencli.GetCmdQueryTokens(), args)
	require.NoError(t, err, "QueryTokensExec failed")
	require.NoError(t, clientCtx.LegacyAmino.UnmarshalJSON(buf.Bytes(), &tokens))
	return tokens
}

// QueryFeeExec executes a query command to retrieve information about a token's fees.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the query will be executed.
// - clientCtx: client.Context - the client context for the query.
// - symbol: string - the symbol of the token.
// - extraArgs: ...string - additional arguments for the query command.
//
// Returns:
// - *v1.QueryFeesResponse - the response containing the token's fees.
func QueryFeeExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	symbol string,
	extraArgs ...string,
) *v1.QueryFeesResponse {
	t.Helper()
	args := []string{
		symbol,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &v1.QueryFeesResponse{}
	network.ExecQueryCmd(t, clientCtx, tokencli.GetCmdQueryFee(), args, response)
	return response
}

// QueryParamsExec executes a query command to retrieve parameters.
//
// Parameters:
// - t: *testing.T - the testing object for running tests.
// - network: simapp.Network - the network on which the query will be executed.
// - clientCtx: client.Context - the client context for the query.
// - extraArgs: ...string - additional arguments for the query command.
//
// Returns:
// - *v1.Params - the response containing the parameters.
func QueryParamsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	extraArgs ...string,
) *v1.Params {
	t.Helper()
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &v1.Params{}
	network.ExecQueryCmd(t, clientCtx, tokencli.GetCmdQueryParams(), args, response)
	return response
}
