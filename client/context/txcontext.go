package context

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/auth"
	authtxb "github.com/irisnet/irishub/client/auth/txbuilder"
	"github.com/irisnet/irishub/client"
	"github.com/irisnet/irishub/client/keys"
)

//----------------------------------------
// Building / Sending utilities

// BaseReq defines a structure that can be embedded in other request structures
// that all share common "base" fields.
type BaseTx struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	ChainID       string `json:"chain_id"`
	AccountNumber int64  `json:"account_number"`
	Sequence      int64  `json:"sequence"`
	Gas           string `json:"gas"`
	GasAdjustment string `json:"gas_adjustment"`
	Fee           string `json:"fee"`
}

// Sanitize performs basic sanitization on a BaseReq object.
func (br BaseTx) Sanitize() BaseTx {
	return BaseTx{
		Name:          strings.TrimSpace(br.Name),
		Password:      strings.TrimSpace(br.Password),
		ChainID:       strings.TrimSpace(br.ChainID),
		Gas:           strings.TrimSpace(br.Gas),
		Fee:           strings.TrimSpace(br.Fee),
		GasAdjustment: strings.TrimSpace(br.GasAdjustment),
		AccountNumber: br.AccountNumber,
		Sequence:      br.Sequence,
	}
}

// ValidateBasic performs basic validation of a BaseReq. If custom validation
// logic is needed, the implementing request handler should perform those
// checks manually.
func (br BaseTx) ValidateBasic(w http.ResponseWriter, cliCtx CLIContext) bool {
	switch {
	case !cliCtx.GenerateOnly && len(br.Name) == 0:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("name required but not specified"))
		return false

	case !cliCtx.DryRun && !cliCtx.GenerateOnly && len(br.Password) == 0:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("password required but not specified"))
		return false

	case len(br.ChainID) == 0:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("chainID required but not specified"))
		return false
	}

	return true
}

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
func (txCtx TxContext) WithCliCtx(ctx CLIContext) TxContext {
	txCtx.cliCtx = ctx
	return txCtx
}

// WithCodec returns a copy of the context with an updated codec.
func (txCtx TxContext) WithCodec(cdc *codec.Codec) TxContext {
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

// Build builds a single message to be signed from a TxContext given a set of
// messages. It returns an error if a fee is supplied but cannot be parsed.
func (txCtx TxContext) Build(msgs []sdk.Msg) (authtxb.StdSignMsg, error) {
	chainID := txCtx.ChainID
	if chainID == "" {
		return authtxb.StdSignMsg{}, errors.Errorf("chain ID required but not specified")
	}

	fee := sdk.Coins{}
	if txCtx.Fee != "" {
		parsedFee, err := txCtx.cliCtx.ParseCoins(txCtx.Fee)
		if err != nil {
			return authtxb.StdSignMsg{}, fmt.Errorf("encountered error in parsing transaction fee: %s", err.Error())
		}

		fee = parsedFee
	}

	return authtxb.StdSignMsg{
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
func (txCtx TxContext) Sign(name, passphrase string, msg authtxb.StdSignMsg) ([]byte, error) {
	sig, err := MakeSignature(name, passphrase, msg)
	if err != nil {
		return nil, err
	}
	return txCtx.Codec.MarshalBinaryLengthPrefixed(auth.NewStdTx(msg.Msgs, msg.Fee, []auth.StdSignature{sig}, msg.Memo))
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

// BuildWithPubKey builds a single message to be signed from a TxContext given a set of
// messages and attach the public key associated to the given name.
// It returns an error if a fee is supplied but cannot be parsed or the key cannot be
// retrieved.
func (txCtx TxContext) BuildWithPubKey(name string, msgs []sdk.Msg) ([]byte, error) {
	msg, err := txCtx.Build(msgs)
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

	return txCtx.Codec.MarshalBinaryLengthPrefixed(auth.NewStdTx(msg.Msgs, msg.Fee, sigs, msg.Memo))
}

// SignStdTx appends a signature to a StdTx and returns a copy of a it. If append
// is false, it replaces the signatures already attached with the new signature.
func (txCtx TxContext) SignStdTx(name, passphrase string, stdTx auth.StdTx, appendSig bool) (signedStdTx auth.StdTx, err error) {
	stdSignature, err := MakeSignature(name, passphrase, authtxb.StdSignMsg{
		ChainID:       txCtx.ChainID,
		AccountNumber: txCtx.AccountNumber,
		Sequence:      txCtx.Sequence,
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
