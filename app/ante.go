package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	ibcante "github.com/cosmos/ibc-go/v3/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	oraclekeeper "github.com/irisnet/irismod/modules/oracle/keeper"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"

	guardiankeeper "github.com/irisnet/irishub/modules/guardian/keeper"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper.
type HandlerOptions struct {
	parent         ante.HandlerOptions
	ibcKeeper      *ibckeeper.Keeper
	tokenKeeper    tokenkeeper.Keeper
	oracleKeeper   oraclekeeper.Keeper
	guardianKeeper guardiankeeper.Keeper
	bankKeeper     bankkeeper.Keeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.parent.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.parent.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	var sigGasConsumer = options.parent.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = ante.DefaultSigVerificationGasConsumer
	}

	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(options.parent.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.parent.AccountKeeper),
		ante.NewDeductFeeDecorator(options.parent.AccountKeeper, options.bankKeeper, options.parent.FeegrantKeeper),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.parent.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.parent.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.parent.AccountKeeper, sigGasConsumer),
		ante.NewSigVerificationDecorator(options.parent.AccountKeeper, options.parent.SignModeHandler),

		NewValidateTokenDecorator(options.tokenKeeper),
		tokenkeeper.NewValidateTokenFeeDecorator(options.tokenKeeper, options.bankKeeper),
		oraclekeeper.NewValidateOracleAuthDecorator(options.oracleKeeper, options.guardianKeeper),
		NewValidateServiceDecorator(),

		ante.NewIncrementSequenceDecorator(options.parent.AccountKeeper),
		ibcante.NewAnteDecorator(options.ibcKeeper),
	), nil
}
