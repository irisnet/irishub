package utils

import (
	"bytes"
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/keys"
	irishubType "github.com/irisnet/irishub/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
)

// SendOrPrintTx implements a utility function that
// facilitates sending a series of messages in a signed
// transaction given a TxContext and a QueryContext. It ensures
// that the account exists, has a proper number and sequence
// set. In addition, it builds and signs a transaction with the
// supplied messages.  Finally, it broadcasts the signed
// transaction to a node.
// NOTE: Also see CompleteAndBroadcastTxREST.
func SendOrPrintTx(txCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	if cliCtx.GenerateOnly {
		return PrintUnsignedStdTx(txCtx, cliCtx, msgs, false)
	}
	// Build and sign the transaction, then broadcast to a Tendermint
	// node.
	cliCtx.PrintResponse = true

	txCtx, err := prepareTxContext(txCtx, cliCtx)
	if err != nil {
		return err
	}

	name, err := cliCtx.GetFromName()
	if err != nil {
		return err
	}

	if txCtx.SimulateGas || cliCtx.DryRun {
		txCtx, err = EnrichCtxWithGas(txCtx, cliCtx, name, msgs)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "estimated gas = %v\n", txCtx.Gas)
	}
	if cliCtx.DryRun {
		return nil
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return err
	}

	// build and sign the transaction
	txBytes, err := txCtx.BuildAndSign(name, passphrase, msgs)
	if err != nil {
		return err
	}
	// broadcast to a Tendermint node
	_, err = cliCtx.BroadcastTx(txBytes)
	return err
}

// EnrichCtxWithGas calculates the gas estimate that would be consumed by the
// transaction and set the transaction's respective value accordingly.
func EnrichCtxWithGas(txCtx context.TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg) (context.TxContext, error) {
	_, adjusted, err := simulateMsgs(txCtx, cliCtx, name, msgs)
	if err != nil {
		return txCtx, err
	}
	return txCtx.WithGas(adjusted), nil
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
func PrintUnsignedStdTx(txCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg, offline bool) (err error) {
	var stdTx auth.StdTx
	if offline {
		stdTx, err = buildUnsignedStdTxOffline(txCtx, cliCtx, msgs)
	} else {
		stdTx, err = buildUnsignedStdTx(txCtx, cliCtx, msgs)
	}
	if err != nil {
		return
	}
	var json []byte
	if cliCtx.Indent {
		json, err = txCtx.Codec.MarshalJSONIndent(stdTx, "", "  ")
	} else {
		json, err = txCtx.Codec.MarshalJSON(stdTx)
	}
	if err == nil {
		fmt.Printf("%s\n", json)
	}
	return
}

// SignStdTx appends a signature to a StdTx and returns a copy of a it. If appendSig
// is false, it replaces the signatures already attached with the new signature.
// Don't perform online validation or lookups if offline is true.
func SignStdTx(txCtx context.TxContext, cliCtx context.CLIContext, name string, stdTx auth.StdTx, appendSig bool, offline bool) (auth.StdTx, error) {
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

	if !offline && txCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(addr)
		if err != nil {
			return signedStdTx, err
		}
		txCtx = txCtx.WithAccountNumber(accNum)
	}

	if !offline && txCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(addr)
		if err != nil {
			return signedStdTx, err
		}
		txCtx = txCtx.WithSequence(accSeq)
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return signedStdTx, err
	}
	return txCtx.SignStdTx(name, passphrase, stdTx, appendSig)
}

// nolint
// SimulateMsgs simulates the transaction and returns the gas estimate and the adjusted value.
func simulateMsgs(txCtx context.TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg) (estimated, adjusted int64, err error) {
	txBytes, err := txCtx.BuildWithPubKey(name, msgs)
	if err != nil {
		return
	}
	estimated, adjusted, err = CalculateGas(cliCtx.Query, cliCtx.Codec, txBytes, txCtx.GasAdjustment)
	return
}

func adjustGasEstimate(estimate int64, adjustment float64) int64 {
	return int64(adjustment * float64(estimate))
}

func parseQueryResponse(cdc *amino.Codec, rawRes []byte) (int64, error) {
	var simulationResult sdk.Result
	if err := cdc.UnmarshalBinaryLengthPrefixed(rawRes, &simulationResult); err != nil {
		return 0, err
	}
	return simulationResult.GasUsed, nil
}

func prepareTxContext(txCtx context.TxContext, cliCtx context.CLIContext) (context.TxContext, error) {
	if err := cliCtx.EnsureAccountExists(); err != nil {
		return txCtx, err
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return txCtx, err
	}

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(from)
		if err != nil {
			return txCtx, err
		}
		txCtx = txCtx.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if txCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(from)
		if err != nil {
			return txCtx, err
		}
		txCtx = txCtx.WithSequence(accSeq)
	}
	return txCtx, nil
}

// buildUnsignedStdTx builds a StdTx as per the parameters passed in the
// contexts. Gas is automatically estimated if gas wanted is set to 0.
func buildUnsignedStdTx(txCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx auth.StdTx, err error) {
	txCtx, err = prepareTxContext(txCtx, cliCtx)
	if err != nil {
		return
	}
	return buildUnsignedStdTxOffline(txCtx, cliCtx, msgs)
}

func buildUnsignedStdTxOffline(txCtx context.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx auth.StdTx, err error) {
	if txCtx.SimulateGas {
		var name string
		name, err = cliCtx.GetFromName()
		if err != nil {
			return
		}

		txCtx, err = EnrichCtxWithGas(txCtx, cliCtx, name, msgs)
		if err != nil {
			return
		}
		fmt.Fprintf(os.Stderr, "estimated gas = %v\n", txCtx.Gas)
	}
	stdSignMsg, err := txCtx.Build(msgs)
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

func ExRateFromStakeTokenToMainUnit(cliCtx context.CLIContext) irishubType.Rat {
	stakeTokenDenom, err := cliCtx.GetCoinType(app.Denom)
	if err != nil {
		panic(err)
	}
	decimalDiff := stakeTokenDenom.MinUnit.Decimal - stakeTokenDenom.GetMainUnit().Decimal
	exRate := irishubType.NewRat(1).Quo(irishubType.NewRatFromInt(sdk.NewIntWithDecimal(1, decimalDiff)))
	return exRate
}

func ConvertDecToRat(input sdk.Dec) irishubType.Rat {
	output, err := irishubType.NewRatFromDecimal(input.String(), 10)
	if err != nil {
		panic(err.Error())
	}
	return output
}
