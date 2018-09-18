package context

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

// TxContext implements a transaction context created in SDK modules.
type TxContext struct {
	Codec         *wire.Codec
	cliCtx        CLIContext
	AccountNumber int64
	Sequence      int64
	Gas           int64
	ChainID       string
	Memo          string
	Fee           string
}

// NewTxContextFromCLI returns a new initialized TxContext with parameters from
// the command line using Viper.
func NewTxContextFromCLI() TxContext {
	// if chain ID is not specified manually, read default chain ID
	chainID := viper.GetString(client.FlagChainID)
	if chainID == "" {
		fmt.Printf("must specify --chain-id")
		os.Exit(1)
	}

	return TxContext{
		ChainID:       chainID,
		Gas:           viper.GetInt64(client.FlagGas),
		AccountNumber: viper.GetInt64(client.FlagAccountNumber),
		Sequence:      viper.GetInt64(client.FlagSequence),
		Fee:           viper.GetString(client.FlagFee),
		Memo:          viper.GetString(client.FlagMemo),
	}
}

// WithCodec returns a copy of the context with an updated codec.
func (txCtx TxContext) WithCodec(cdc *wire.Codec) TxContext {
	txCtx.Codec = cdc
	return txCtx
}

// WithChainID returns a copy of the context with an updated chainID.
func (txCtx TxContext) WithChainID(chainID string) TxContext {
	txCtx.ChainID = chainID
	return txCtx
}

// WithGas returns a copy of the context with an updated gas.
func (txCtx TxContext) WithGas(gas int64) TxContext {
	txCtx.Gas = gas
	return txCtx
}

// WithFee returns a copy of the context with an updated fee.
func (txCtx TxContext) WithFee(fee string) TxContext {
	txCtx.Fee = fee
	return txCtx
}

// WithSequence returns a copy of the context with an updated sequence number.
func (txCtx TxContext) WithSequence(sequence int64) TxContext {
	txCtx.Sequence = sequence
	return txCtx
}

// WithMemo returns a copy of the context with an updated memo.
func (txCtx TxContext) WithMemo(memo string) TxContext {
	txCtx.Memo = memo
	return txCtx
}

// WithAccountNumber returns a copy of the context with an account number.
func (txCtx TxContext) WithAccountNumber(accnum int64) TxContext {
	txCtx.AccountNumber = accnum
	return txCtx
}

// WithCliCtx returns a copy of the context with a CLIContext
func (txCtx TxContext) WithCliCtx(cliCtx CLIContext) TxContext {
	txCtx.cliCtx = cliCtx
	return txCtx
}

// Build builds a single message to be signed from a TxContext given a set of
// messages. It returns an error if a fee is supplied but cannot be parsed.
func (txCtx TxContext) Build(msgs []sdk.Msg) (auth.StdSignMsg, error) {
	chainID := txCtx.ChainID
	if chainID == "" {
		return auth.StdSignMsg{}, errors.Errorf("chain ID required but not specified")
	}

	fee := sdk.Coins{}
	if txCtx.Fee != "" {
		parsedFee, err := txCtx.cliCtx.ParseCoins(txCtx.Fee)
		if err != nil {
			return auth.StdSignMsg{}, err
		}

		fee = parsedFee
	}

	return auth.StdSignMsg{
		ChainID:       txCtx.ChainID,
		AccountNumber: txCtx.AccountNumber,
		Sequence:      txCtx.Sequence,
		Memo:          txCtx.Memo,
		Msgs:          msgs,
		Fee:           auth.NewStdFee(txCtx.Gas, fee...),
	}, nil
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (txCtx TxContext) Sign(name, passphrase string, msg auth.StdSignMsg) ([]byte, error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	sig, pubkey, err := keybase.Sign(name, passphrase, msg.Bytes())
	if err != nil {
		return nil, err
	}

	sigs := []auth.StdSignature{{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
		PubKey:        pubkey,
		Signature:     sig,
	}}

	return txCtx.Codec.MarshalBinary(auth.NewStdTx(msg.Msgs, msg.Fee, sigs, msg.Memo))
}

// BuildAndSign builds a single message to be signed, and signs a transaction
// with the built message given a name, passphrase, and a set of
// messages.
func (txCtx TxContext) BuildAndSign(name, passphrase string, msgs []sdk.Msg) ([]byte, error) {
	msg, err := txCtx.Build(msgs)
	if err != nil {
		return nil, err
	}

	return txCtx.Sign(name, passphrase, msg)
}
