package testutil

import (
	"encoding/json"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/irisnet/irishub/v4/app"
	"github.com/irisnet/irishub/v4/app/params"
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
	encCdc := app.RegisterEncodingConfig()
	if appOpts == nil {
		appOpts = EmptyAppOptions{}
	}
	app := app.NewIrisApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		encCdc,
		appOpts,
		baseAppOptions...,
	)
	return &AppWrapper{app}
}

// MakeCodecs returns the application codec and tx codec
func MakeCodecs() params.EncodingConfig {
	return app.RegisterEncodingConfig()
}

// DefaultGenesis returns default genesis state as raw bytes
func DefaultGenesis(cdc codec.JSONCodec) map[string]json.RawMessage {
	return app.ModuleBasics.DefaultGenesis(cdc)
}
