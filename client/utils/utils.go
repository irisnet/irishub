package utils

import (
	"bytes"
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/keys"
)

// CompleteAndBroadcastTxCli implements a utility function that
// facilitates sending a series of messages in a signed
// transaction given a TxContext and a QueryContext. It ensures
// that the account exists, has a proper number and sequence
// set. In addition, it builds and signs a transaction with the
// supplied messages.  Finally, it broadcasts the signed
// transaction to a node.
// NOTE: Also see CompleteAndBroadcastTxREST.
func CompleteAndBroadcastTxCli(TxCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	TxCtx, err := prepareTxBuilder(TxCtx, cliCtx)
	if err != nil {
		return err
	}

	name, err := cliCtx.GetFromName()
	if err != nil {
		return err
	}

	if TxCtx.SimulateGas || cliCtx.DryRun {
		TxCtx, err = EnrichCtxWithGas(TxCtx, cliCtx, name, msgs)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "estimated gas = %v\n", TxCtx.Gas)
	}
	if cliCtx.DryRun {
		return nil
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return err
	}

	// build and sign the transaction
	txBytes, err := TxCtx.BuildAndSign(name, passphrase, msgs)
	if err != nil {
		return err
	}
	// broadcast to a Tendermint node
	_, err = cliCtx.BroadcastTx(txBytes)
	return err
}

// EnrichCtxWithGas calculates the gas estimate that would be consumed by the
// transaction and set the transaction's respective value accordingly.
func EnrichCtxWithGas(TxCtx context.TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg) (context.TxContext, error) {
	_, adjusted, err := simulateMsgs(TxCtx, cliCtx, name, msgs)
	if err != nil {
		return TxCtx, err
	}
	return TxCtx.WithGas(adjusted), nil
}

// CalculateGas simulates the execution of a transaction and returns
// both the estimate obtained by the query and the adjusted amount.
func CalculateGas(queryFunc func(string, common.HexBytes) ([]byte, error), cdc *amino.Codec, txBytes []byte, adjustment float64) (estimate, adjusted int64, err error) {
	// run a simulation (via /app/simulate query) to
	// estimate gas and update TxContext accordingly
	rawRes, err := queryFunc("/app/simulate", txBytes)
	if err != nil {
		return
	}
	estimate, err = parseQueryResponse(cdc, rawRes)
	if err != nil {
		return
	}
	adjusted = adjustGasEstimate(estimate, adjustment)
	return
}

// PrintUnsignedStdTx builds an unsigned StdTx and prints it to os.Stdout.
// Don't perform online validation or lookups if offline is true.
func PrintUnsignedStdTx(TxCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg, offline bool) (err error) {
	var stdTx auth.StdTx
	if offline {
		stdTx, err = buildUnsignedStdTxOffline(TxCtx, cliCtx, msgs)
	} else {
		stdTx, err = buildUnsignedStdTx(TxCtx, cliCtx, msgs)
	}
	if err != nil {
		return
	}
	var json []byte
	if cliCtx.Indent {
		json, err = TxCtx.Codec.MarshalJSONIndent(stdTx, "", "  ")
	} else {
		json, err = TxCtx.Codec.MarshalJSON(stdTx)
	}
	if err == nil {
		fmt.Printf("%s\n", json)
	}
	return
}

// SignStdTx appends a signature to a StdTx and returns a copy of a it. If appendSig
// is false, it replaces the signatures already attached with the new signature.
// Don't perform online validation or lookups if offline is true.
func SignStdTx(TxCtx context.TxContext, cliCtx context.CLIContext, name string, stdTx auth.StdTx, appendSig bool, offline bool) (auth.StdTx, error) {
	var signedStdTx auth.StdTx

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return signedStdTx, err
	}
	info, err := keybase.Get(name)
	if err != nil {
		return signedStdTx, err
	}
	addr := info.GetPubKey().Address()

	// Check whether the address is a signer
	if !isTxSigner(sdk.AccAddress(addr), stdTx.GetSigners()) {
		fmt.Fprintf(os.Stderr, "WARNING: The generated transaction's intended signer does not match the given signer: '%v'\n", name)
	}

	if !offline && TxCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(addr)
		if err != nil {
			return signedStdTx, err
		}
		TxCtx = TxCtx.WithAccountNumber(accNum)
	}

	if !offline && TxCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(addr)
		if err != nil {
			return signedStdTx, err
		}
		TxCtx = TxCtx.WithSequence(accSeq)
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return signedStdTx, err
	}
	return TxCtx.SignStdTx(name, passphrase, stdTx, appendSig)
}

// nolint
// SimulateMsgs simulates the transaction and returns the gas estimate and the adjusted value.
func simulateMsgs(TxCtx context.TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg) (estimated, adjusted int64, err error) {
	txBytes, err := TxCtx.BuildWithPubKey(name, msgs)
	if err != nil {
		return
	}
	estimated, adjusted, err = CalculateGas(cliCtx.Query, cliCtx.Codec, txBytes, TxCtx.GasAdjustment)
	return
}

func adjustGasEstimate(estimate int64, adjustment float64) int64 {
	return int64(adjustment * float64(estimate))
}

func parseQueryResponse(cdc *amino.Codec, rawRes []byte) (int64, error) {
	var simulationResult sdk.Result
	if err := cdc.UnmarshalBinary(rawRes, &simulationResult); err != nil {
		return 0, err
	}
	return simulationResult.GasUsed, nil
}

func prepareTxBuilder(TxCtx context.TxContext, cliCtx context.CLIContext) (context.TxContext, error) {
	if err := cliCtx.EnsureAccountExists(); err != nil {
		return TxCtx, err
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return TxCtx, err
	}

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if TxCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(from)
		if err != nil {
			return TxCtx, err
		}
		TxCtx = TxCtx.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if TxCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(from)
		if err != nil {
			return TxCtx, err
		}
		TxCtx = TxCtx.WithSequence(accSeq)
	}
	return TxCtx, nil
}

// buildUnsignedStdTx builds a StdTx as per the parameters passed in the
// contexts. Gas is automatically estimated if gas wanted is set to 0.
func buildUnsignedStdTx(TxCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx auth.StdTx, err error) {
	TxCtx, err = prepareTxBuilder(TxCtx, cliCtx)
	if err != nil {
		return
	}
	return buildUnsignedStdTxOffline(TxCtx, cliCtx, msgs)
}

func buildUnsignedStdTxOffline(TxCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx auth.StdTx, err error) {
	if TxCtx.SimulateGas {
		var name string
		name, err = cliCtx.GetFromName()
		if err != nil {
			return
		}

		TxCtx, err = EnrichCtxWithGas(TxCtx, cliCtx, name, msgs)
		if err != nil {
			return
		}
		fmt.Fprintf(os.Stderr, "estimated gas = %v\n", TxCtx.Gas)
	}
	stdSignMsg, err := TxCtx.Build(msgs)
	if err != nil {
		return
	}
	return auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, nil, stdSignMsg.Memo), nil
}

func isTxSigner(user sdk.AccAddress, signers []sdk.AccAddress) bool {
	for _, s := range signers {
		if bytes.Equal(user.Bytes(), s.Bytes()) {
			return true
		}
	}
	return false
}
