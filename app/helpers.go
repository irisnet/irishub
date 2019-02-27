package app

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/abci/server"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// nolint - Mostly for testing
func (app *BaseApp) Check(tx sdk.Tx, txBytes []byte) (result sdk.Result) {
	return app.runTx(RunTxModeCheck, txBytes, tx)
}

// nolint - full tx execution
func (app *BaseApp) Simulate(tx sdk.Tx, txBytes []byte) (result sdk.Result) {
	return app.runTx(RunTxModeSimulate, txBytes, tx)
}

// nolint
func (app *BaseApp) Deliver(tx sdk.Tx, txBytes []byte) (result sdk.Result) {
	return app.runTx(RunTxModeDeliver, txBytes, tx)
}

// RunForever - BasecoinApp execution and cleanup
func RunForever(app abci.Application) {

	// Start the ABCI server
	srv, err := server.NewServer("0.0.0.0:26658", "socket", app)
	if err != nil {
		cmn.Exit(err.Error())
		return
	}
	err = srv.Start()
	if err != nil {
		cmn.Exit(err.Error())
		return
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		err := srv.Stop()
		if err != nil {
			cmn.Exit(err.Error())
		}
	})
}
