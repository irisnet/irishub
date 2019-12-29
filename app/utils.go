//nolint
package app

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/irisnet/irishub/modules/asset"
	token "github.com/irisnet/irishub/modules/asset/01-token"
)

// ExportStateToJSON util function to export the app state to JSON
func ExportStateToJSON(app *IrisApp, path string) error {
	fmt.Println("exporting app state...")
	appState, _, err := app.ExportAppStateAndValidators(false, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(appState), 0644)
}

// NewIrisAppUNSAFE is used for debugging purposes only.
//
// NOTE: to not use this function with non-test code
func NewIrisAppUNSAFE(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp),
) (gapp *IrisApp, keyMain, keyStaking *sdk.KVStoreKey, stakingKeeper staking.Keeper) {
	gapp = NewIrisApp(logger, db, traceStore, loadLatest, invCheckPeriod, baseAppOptions...)
	return gapp, gapp.keys[bam.MainStoreKey], gapp.keys[staking.StoreKey], gapp.stakingKeeper
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(ak keeper.AccountKeeper, supplyKeeper types.SupplyKeeper, assetKeeper asset.Keeper, sigGasConsumer ante.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, supplyKeeper),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak),
		ante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
		token.NewValidateTokenFeeDecorator(assetKeeper.TokenKeeper, ak),
	)
}
