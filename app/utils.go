//nolint
package app

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
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

// SHA256 wraps sha256.Sum256 with result converted to slice
func SHA256(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
