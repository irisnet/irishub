package mt

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	mtcli "mods.irisnet.org/modules/mt/client/cli"
	mttypes "mods.irisnet.org/modules/mt/types"
	"mods.irisnet.org/simapp"
)

// IssueDenomExec executes the IssueDenom command with the specified parameters.
//
// Parameters:
// - t: The testing.T object for logging and reporting.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object representing the client context.
// - from: The address of the account issuing the denom.
// - extraArgs: Additional command line arguments.
//
// Returns:
// - *simapp.ResponseTx: The response transaction object.
func IssueDenomExec(
	t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdIssueDenom(), args)
}


// BurnMTExec executes a burn token transaction.
//
// Parameters:
// - t: The testing.T object for logging and reporting.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object representing the client context.
// - from: The address of the account initiating the burn transaction.
// - denomID: The unique identifier of the denomination to burn.
// - mtID: The unique identifier of the multi-token to burn.
// - amount: The amount of tokens to burn.
// - extraArgs: Additional command line arguments.
//
// Returns:
// - *simapp.ResponseTx: The response transaction object.
func BurnMTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	mtID string,
	amount string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		denomID,
		mtID,
		amount,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdBurnMT(), args)
}

// MintMTExec executes a mint token transaction.
//
// Parameters:
// - t: The testing.T object for logging and reporting.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object representing the client context.
// - from: The address of the account initiating the mint transaction.
// - denomID: The unique identifier of the denomination to mint.
// - extraArgs: Additional command line arguments.
//
// Returns:
// - *simapp.ResponseTx: The response transaction object.
func MintMTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		denomID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdMintMT(), args)
}

// EditMTExec executes an edit MT transaction.
//
// Parameters:
// - t: The testing.T object for logging and reporting.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object representing the client context.
// - from: The address of the account initiating the edit transaction.
// - denomID: The unique identifier of the denomination to edit.
// - mtID: The unique identifier of the MT to edit.
// - extraArgs: Additional command line arguments.
//
// Returns:
// - *simapp.ResponseTx: The response transaction object.
func EditMTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	mtID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		denomID,
		mtID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdEditMT(), args)
}

// TransferMTExec executes a transfer MT transaction.
//
// Parameters:
// - t: The testing.T object for logging and reporting.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object representing the client context.
// - from: The address of the account initiating the transfer transaction.
// - recipient: The address of the account receiving the transferred tokens.
// - denomID: The unique identifier of the denomination.
// - mtID: The unique identifier of the MT being transferred.
// - amount: The amount of tokens to transfer.
// - extraArgs: Additional command line arguments.
//
// Returns:
// - *simapp.ResponseTx: The response transaction object.
func TransferMTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	recipient string,
	denomID string,
	mtID string,
	amount string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		from,
		recipient,
		denomID,
		mtID,
		amount,
	}
	args = append(args, extraArgs...)

	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdTransferMT(), args)
}

// QueryDenomExec executes a query command to retrieve a specific denom from the network.
//
// Parameters:
// - t: The testing.T object for testing.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object for the client.
// - denomID: The ID of the denom to query.
// - extraArgs: Additional arguments to be passed to the command.
//
// Returns:
// - *mttypes.Denom: The response object containing the queried denom.
func QueryDenomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	extraArgs ...string) *mttypes.Denom {
	args := []string{
		denomID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &mttypes.Denom{}
	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryDenom(), args, response)
	return response
}

// QueryDenomsExec executes a query command to retrieve all denoms from the network.
//
// Parameters:
// - t: The testing.T object for testing.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object for the client.
// - extraArgs: Additional arguments to be passed to the command.
//
// Returns:
// - *mttypes.QueryDenomsResponse: The response object containing the queried denoms.
func QueryDenomsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	extraArgs ...string) *mttypes.QueryDenomsResponse {
	args := []string{
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &mttypes.QueryDenomsResponse{}
	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryDenoms(), args, response)
	return response
}

// QueryMTsExec executes a query command to retrieve all MTs from the network.
//
// Parameters:
// - t: The testing.T object for testing.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object for the client.
// - denomID: The ID of the denom to query.
// - resp: The response object to store the queried MTs.
// - extraArgs: Additional arguments to be passed to the command.
//
// Returns: None.
func QueryMTsExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	resp proto.Message,
	extraArgs ...string,
) {
	args := []string{
		denomID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryMTs(), args, resp)
}

// QueryMTExec executes a query command to retrieve a specific MT from the network.
//
// Parameters:
// - t: The testing.T object for testing.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object for the client.
// - denomID: The ID of the denom containing the MT.
// - mtID: The ID of the MT to query.
// - extraArgs: Additional arguments to be passed to the command.
//
// Returns:
// - *mttypes.MT: The response object containing the queried MT.
func QueryMTExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	denomID string,
	mtID string,
	extraArgs ...string) *mttypes.MT {
	args := []string{
		denomID,
		mtID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)
	response := &mttypes.MT{}
	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryMT(), args, response)
	return response
}

// QueryBlancesExec executes a query command to retrieve the balances of a specific account for a given denomination.
//
// Parameters:
// - t: The testing.T object for testing.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object for the client.
// - from: The address of the account to query balances for.
// - denomID: The ID of the denomination to query balances for.
// - extraArgs: Additional arguments to be passed to the command.
//
// Returns:
// - *mttypes.QueryBalancesResponse: The response object containing the queried balances.
func QueryBlancesExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	denomID string,
	extraArgs ...string) *mttypes.QueryBalancesResponse {
	args := []string{
		from,
		denomID,
		fmt.Sprintf("--%s=json", cli.OutputFlag),
	}
	args = append(args, extraArgs...)

	response := &mttypes.QueryBalancesResponse{}
	network.ExecQueryCmd(t, clientCtx, mtcli.GetCmdQueryBalances(), args, response)
	return response
}

// TransferDenomExec executes a transfer denomination transaction.
//
// Parameters:
// - t: The testing.T object for logging and reporting.
// - network: The simapp.Network object representing the network.
// - clientCtx: The client.Context object representing the client context.
// - from: The address of the account initiating the transfer transaction.
// - recipient: The address of the account receiving the transferred tokens.
// - denomID: The unique identifier of the denomination.
// - extraArgs: Additional command line arguments.
//
// Returns:
// - *simapp.ResponseTx: The response transaction object.
func TransferDenomExec(t *testing.T,
	network simapp.Network,
	clientCtx client.Context,
	from string,
	recipient string,
	denomID string,
	extraArgs ...string,
) *simapp.ResponseTx {
	args := []string{
		from,
		recipient,
		denomID,
	}

	args = append(args, extraArgs...)
	return network.ExecTxCmdWithResult(t, clientCtx, mtcli.GetCmdTransferDenom(), args)
}
