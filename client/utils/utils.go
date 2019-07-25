package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/stake/types"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
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
func SendOrPrintTx(txCtx TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) error {
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
		var result sdk.Result
		txCtx, result, err = EnrichCtxWithGas(txCtx, cliCtx, name, msgs)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "estimated gas = %v\n", txCtx.Gas)
		fmt.Fprintf(os.Stderr, "simulation code = %v\n", result.Code)
		fmt.Fprintf(os.Stderr, "simulation log = %v\n", result.Log)
		fmt.Fprintf(os.Stderr, "simulation gas wanted = %v\n", result.GasWanted)
		fmt.Fprintf(os.Stderr, "simulation gas used = %v\n", result.GasUsed)
		fmt.Fprintf(os.Stderr, "simulation fee amount = %v\n", result.FeeAmount)
		fmt.Fprintf(os.Stderr, "simulation fee denom = %v\n", result.FeeDenom)
		for _, tag := range result.Tags {
			fmt.Fprintf(os.Stderr, "simulation tag %s = %s\n", string(tag.Key), string(tag.Value))
		}
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
func EnrichCtxWithGas(txCtx TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg) (TxContext, sdk.Result, error) {
	_, adjusted, result, err := simulateMsgs(txCtx, cliCtx, name, msgs)
	if err != nil {
		return txCtx, sdk.Result{}, err
	}
	return txCtx.WithGas(adjusted), result, nil
}

// CalculateGas simulates the execution of a transaction and returns
// both the estimate obtained by the query and the adjusted amount.
func CalculateGas(queryFunc func(string, common.HexBytes) ([]byte, error), cdc *amino.Codec, txBytes []byte, adjustment float64) (estimate, adjusted uint64, simulationResult sdk.Result, err error) {
	// run a simulation (via /app/simulate query) to
	// estimate gas and update TxContext accordingly
	rawRes, err := queryFunc("/app/simulate", txBytes)
	if err != nil {
		return
	}
	if err := cdc.UnmarshalBinaryLengthPrefixed(rawRes, &simulationResult); err != nil {
		return 0, 0, sdk.Result{}, err
	}
	estimate = simulationResult.GasUsed
	adjusted = adjustGasEstimate(estimate, adjustment)
	return
}

// PrintUnsignedStdTx builds an unsigned StdTx and prints it to os.Stdout.
// Don't perform online validation or lookups if offline is true.
func PrintUnsignedStdTx(txCtx TxContext, cliCtx context.CLIContext, msgs []sdk.Msg, offline bool) (err error) {
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
func SignStdTx(txCtx TxContext, cliCtx context.CLIContext, name string, stdTx auth.StdTx, appendSig bool, offline bool) (auth.StdTx, error) {
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

	if !offline {
		txCtx, err = populateAccountFromState(txCtx, cliCtx, sdk.AccAddress(addr))
		if err != nil {
			return signedStdTx, err
		}
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return signedStdTx, err
	}
	return txCtx.SignStdTx(name, passphrase, stdTx, appendSig)
}

// SignStdTxWithSignerAddress attaches a signature to a StdTx and returns a copy of a it.
// Don't perform online validation or lookups if offline is true, else
// populate account and sequence numbers from a foreign account.
func SignStdTxWithSignerAddress(txCtx TxContext, cliCtx context.CLIContext,
	addr sdk.AccAddress, name string, stdTx auth.StdTx,
	offline bool) (signedStdTx auth.StdTx, err error) {

	// Check whether the address is a signer
	if !isTxSigner(sdk.AccAddress(addr), stdTx.GetSigners()) {
		fmt.Fprintf(os.Stderr, "WARNING: The generated transaction's intended signer does not match the given signer: '%v'\n", name)
	}

	if !offline {
		txCtx, err = populateAccountFromState(txCtx, cliCtx, addr)
		if err != nil {
			return signedStdTx, err
		}
	}

	passphrase, err := keys.GetPassphrase(name)
	if err != nil {
		return signedStdTx, err
	}

	return txCtx.SignStdTx(name, passphrase, stdTx, false)
}

// nolint
// SimulateMsgs simulates the transaction and returns the gas estimate and the adjusted value.
func simulateMsgs(txCtx TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg) (estimated, adjusted uint64, result sdk.Result, err error) {
	txBytes, err := txCtx.BuildWithPubKey(name, msgs)
	if err != nil {
		return
	}
	estimated, adjusted, result, err = CalculateGas(cliCtx.Query, cliCtx.Codec, txBytes, txCtx.GasAdjustment)
	return
}

func adjustGasEstimate(estimate uint64, adjustment float64) uint64 {
	return uint64(adjustment * float64(estimate))
}

func prepareTxContext(txCtx TxContext, cliCtx context.CLIContext) (TxContext, error) {
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
func buildUnsignedStdTx(txCtx TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx auth.StdTx, err error) {
	txCtx, err = prepareTxContext(txCtx, cliCtx)
	if err != nil {
		return
	}
	return buildUnsignedStdTxOffline(txCtx, cliCtx, msgs)
}

func populateAccountFromState(
	txCtx TxContext, cliCtx context.CLIContext, addr sdk.AccAddress,
) (TxContext, error) {

	accNum, err := cliCtx.GetAccountNumber(addr)
	if err != nil {
		return txCtx, err
	}

	accSeq, err := cliCtx.GetAccountSequence(addr)
	if err != nil {
		return txCtx, err
	}

	return txCtx.WithAccountNumber(accNum).WithSequence(accSeq), nil
}

func buildUnsignedStdTxOffline(txCtx TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx auth.StdTx, err error) {
	if txCtx.SimulateGas {
		var name string
		name, err = cliCtx.GetFromName()
		if err != nil {
			return
		}

		txCtx, _, err = EnrichCtxWithGas(txCtx, cliCtx, name, msgs)
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

func ExRateFromStakeTokenToMainUnit(cliCtx context.CLIContext) sdk.Rat {
	stakeToken, err := cliCtx.GetCoinType(types.StakeTokenName)
	if err != nil {
		panic(err)
	}
	decimalDiff := stakeToken.MinUnit.Decimal - stakeToken.GetMainUnit().Decimal
	exRate := sdk.NewRat(1).Quo(sdk.NewRatFromInt(sdk.NewIntWithDecimal(1, int(decimalDiff))))
	return exRate
}

func ConvertDecToRat(input sdk.Dec) sdk.Rat {
	output, err := sdk.NewRatFromDecimal(input.String(), 10)
	if err != nil {
		panic(err.Error())
	}
	return output
}

// GetAccountDecoder gets the account decoder for auth.DefaultAccount.
func GetAccountDecoder(cdc *codec.Codec) auth.AccountDecoder {
	return func(accBytes []byte) (acct auth.Account, err error) {
		err = cdc.UnmarshalBinaryBare(accBytes, &acct)
		if err != nil {
			panic(err)
		}

		return acct, err
	}
}
