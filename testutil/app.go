package testutil

import (
	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/irisnet/irishub/v4/app"
)

var (
	_ runtime.AppI            = (*AppWrapper)(nil)
	_ servertypes.Application = (*AppWrapper)(nil)
)

// AppWrapper extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type AppWrapper struct {
	*app.IrisApp
}

func setup(
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *AppWrapper {
	db := dbm.NewMemDB()
	if appOpts == nil {
		appOpts = EmptyAppOptions{}
	}
	app := app.NewIrisApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		appOpts,
		baseAppOptions...,
	)
	return &AppWrapper{app}
}