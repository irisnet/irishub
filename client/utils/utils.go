package utils

import (
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/app"
)

// SendTx implements a auxiliary handler that facilitates sending a series of
// messages in a signed transaction given a TxContext and a QueryContext. It
// ensures that the account exists, has a proper number and sequence set. In
// addition, it builds and signs a transaction with the supplied messages.
// Finally, it broadcasts the signed transaction to a node.
func SendTx(ctx app.Context, msgs []sdk.Msg) error {
	txCtx := ctx.GetTxCxt()
	if err := ctx.EnsureAccountExists(); err != nil {
		return err
	}

	from, err := ctx.GetFromAddress()
	if err != nil {
		return err
	}

	// automatically doing a manual lookup.
	if txCtx.AccountNumber == 0 {
		accNum, err := ctx.GetAccountNumber(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithAccountNumber(accNum)
	}

	// automatically doing a manual lookup.
	if txCtx.Sequence == 0 {
		accSeq, err := ctx.GetAccountSequence(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithSequence(accSeq)
	}

	passphrase, err := keys.GetPassphrase(ctx.FromAddressName)
	if err != nil {
		return err
	}

	// build and sign the transaction
	ctx = ctx.WithTxContext(txCtx)
	msg, err := ctx.Build(msgs)
	if err != nil {
		return err
	}
	txBytes, err := txCtx.Sign(ctx.FromAddressName, passphrase, msg)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	return ctx.EnsureBroadcastTx(txBytes)
}
