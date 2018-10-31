package context

import (
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"fmt"
)

// TxContext implements a transaction context created in SDK modules.
type TxContext struct {
	Codec         *codec.Codec
	cliCtx        CLIContext
	AccountNumber int64
	Sequence      int64
	Gas           int64 // TODO: should this turn into uint64? requires further discussion - see #2173
	GasAdjustment float64
	SimulateGas   bool
	ChainID       string
	Memo          string
	Fee           string
}

// NewTxBuilderFromCLI returns a new initialized TxContext with parameters from
// the command line using Viper.
func NewTxContextFromCLI() TxContext {
	// if chain ID is not specified manually, read default chain ID
	chainID := viper.GetString(client.FlagChainID)
	if chainID == "" {
		defaultChainID, err := sdk.DefaultChainID()
		if err != nil {
			chainID = defaultChainID
		}
	}

	return TxContext{
		ChainID:       chainID,
		AccountNumber: viper.GetInt64(client.FlagAccountNumber),
		Gas:           client.GasFlagVar.Gas,
		GasAdjustment: viper.GetFloat64(client.FlagGasAdjustment),
		Sequence:      viper.GetInt64(client.FlagSequence),
		SimulateGas:   client.GasFlagVar.Simulate,
		Fee:           viper.GetString(client.FlagFee),
		Memo:          viper.GetString(client.FlagMemo),
	}
}

// WithCodec returns a copy of the context with an updated codec.
func (bldr TxContext) WithCliCtx(ctx CLIContext) TxContext {
	bldr.cliCtx = ctx
	return bldr
}

// WithCodec returns a copy of the context with an updated codec.
func (bldr TxContext) WithCodec(cdc *codec.Codec) TxContext {
	bldr.Codec = cdc
	return bldr
}

// WithChainID returns a copy of the context with an updated chainID.
func (bldr TxContext) WithChainID(chainID string) TxContext {
	bldr.ChainID = chainID
	return bldr
}

// WithGas returns a copy of the context with an updated gas.
func (bldr TxContext) WithGas(gas int64) TxContext {
	bldr.Gas = gas
	return bldr
}

// WithFee returns a copy of the context with an updated fee.
func (bldr TxContext) WithFee(fee string) TxContext {
	bldr.Fee = fee
	return bldr
}

// WithSequence returns a copy of the context with an updated sequence number.
func (bldr TxContext) WithSequence(sequence int64) TxContext {
	bldr.Sequence = sequence
	return bldr
}

// WithMemo returns a copy of the context with an updated memo.
func (bldr TxContext) WithMemo(memo string) TxContext {
	bldr.Memo = memo
	return bldr
}

// WithAccountNumber returns a copy of the context with an account number.
func (bldr TxContext) WithAccountNumber(accnum int64) TxContext {
	bldr.AccountNumber = accnum
	return bldr
}

// Build builds a single message to be signed from a TxContext given a set of
// messages. It returns an error if a fee is supplied but cannot be parsed.
func (bldr TxContext) Build(msgs []sdk.Msg) (authtxb.StdSignMsg, error) {
	chainID := bldr.ChainID
	if chainID == "" {
		return authtxb.StdSignMsg{}, errors.Errorf("chain ID required but not specified")
	}

	fee := sdk.Coins{}
	if bldr.Fee != "" {
		parsedFee, err := bldr.cliCtx.ParseCoins(bldr.Fee)
		if err != nil {
			return authtxb.StdSignMsg{}, fmt.Errorf("encountered error in parsing transaction fee: %s", err.Error())
		}

		fee = parsedFee
	}

	return authtxb.StdSignMsg{
		ChainID:       bldr.ChainID,
		AccountNumber: bldr.AccountNumber,
		Sequence:      bldr.Sequence,
		Memo:          bldr.Memo,
		Msgs:          msgs,
		Fee:           auth.NewStdFee(bldr.Gas, fee...),
	}, nil
}

// Sign signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (bldr TxContext) Sign(name, passphrase string, msg authtxb.StdSignMsg) ([]byte, error) {
	sig, err := MakeSignature(name, passphrase, msg)
	if err != nil {
		return nil, err
	}
	return bldr.Codec.MarshalBinary(auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo))
}

// BuildAndSign builds a single message to be signed, and signs a transaction
// with the built message given a name, passphrase, and a set of
// messages.
func (bldr TxContext) BuildAndSign(name, passphrase string, msgs []sdk.Msg) ([]byte, error) {
	msg, err := bldr.Build(msgs)
	if err != nil {
		return nil, err
	}

	return bldr.Sign(name, passphrase, msg)
}

// BuildWithPubKey builds a single message to be signed from a TxContext given a set of
// messages and attach the public key associated to the given name.
// It returns an error if a fee is supplied but cannot be parsed or the key cannot be
// retrieved.
func (bldr TxContext) BuildWithPubKey(name string, msgs []sdk.Msg) ([]byte, error) {
	msg, err := bldr.Build(msgs)
	if err != nil {
		return nil, err
	}

	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	info, err := keybase.Get(name)
	if err != nil {
		return nil, err
	}

	sigs := []auth.StdSignature{{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
		PubKey:        info.GetPubKey(),
	}}

	return bldr.Codec.MarshalBinary(auth.NewStdTx(msg.Msgs, msg.Fee, sigs, msg.Memo))
}

// SignStdTx appends a signature to a StdTx and returns a copy of a it. If append
// is false, it replaces the signatures already attached with the new signature.
func (bldr TxContext) SignStdTx(name, passphrase string, stdTx auth.StdTx, appendSig bool) (signedStdTx auth.StdTx, err error) {
	stdSignature, err := MakeSignature(name, passphrase, authtxb.StdSignMsg{
		ChainID:       bldr.ChainID,
		AccountNumber: bldr.AccountNumber,
		Sequence:      bldr.Sequence,
		Fee:           stdTx.Fee,
		Msgs:          stdTx.GetMsgs(),
		Memo:          stdTx.GetMemo(),
	})
	if err != nil {
		return
	}

	sigs := stdTx.GetSignatures()
	if len(sigs) == 0 || !appendSig {
		sigs = []auth.StdSignature{stdSignature}
	} else {
		sigs = append(sigs, stdSignature)
	}
	signedStdTx = auth.NewStdTx(stdTx.GetMsgs(), stdTx.Fee, sigs, stdTx.GetMemo())
	return
}

// MakeSignature builds a StdSignature given key name, passphrase, and a StdSignMsg.
func MakeSignature(name, passphrase string, msg authtxb.StdSignMsg) (sig auth.StdSignature, err error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return
	}
	sigBytes, pubkey, err := keybase.Sign(name, passphrase, msg.Bytes())
	if err != nil {
		return
	}
	return auth.StdSignature{
		AccountNumber: msg.AccountNumber,
		Sequence:      msg.Sequence,
		PubKey:        pubkey,
		Signature:     sigBytes,
	}, nil
}
