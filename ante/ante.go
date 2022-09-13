package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	guardiankeeper "github.com/irisnet/irishub/modules/guardian/keeper"

	oraclekeeper "github.com/irisnet/irismod/modules/oracle/keeper"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	ante.HandlerOptions
	BankKeeper           bankkeeper.Keeper
	TokenKeeper          tokenkeeper.Keeper
	OracleKeeper         oraclekeeper.Keeper
	GuardianKeeper       guardiankeeper.Keeper
	BypassMinFeeMsgTypes []string
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(opts HandlerOptions) (sdk.AnteHandler, error) {
	var sigGasConsumer = opts.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewExtensionOptionsDecorator(opts.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(opts.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(opts.AccountKeeper),
		ante.NewDeductFeeDecorator(opts.AccountKeeper, opts.BankKeeper, opts.FeegrantKeeper, opts.TxFeeChecker),
		ante.NewSetPubKeyDecorator(opts.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(opts.AccountKeeper),
		ante.NewSigGasConsumeDecorator(opts.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(opts.AccountKeeper, opts.SignModeHandler),
		NewValidateTokenDecorator(opts.TokenKeeper),
		tokenkeeper.NewValidateTokenFeeDecorator(opts.TokenKeeper, opts.BankKeeper),
		oraclekeeper.NewValidateOracleAuthDecorator(opts.OracleKeeper, opts.GuardianKeeper),
		NewValidateServiceDecorator(),
		ante.NewIncrementSequenceDecorator(opts.AccountKeeper),
	), nil
}
