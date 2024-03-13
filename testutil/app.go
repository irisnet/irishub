package testutil

import (
	"encoding/json"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"

	"github.com/irisnet/irishub/v3/app"
	"github.com/irisnet/irishub/v3/app/params"
)

var (
	_ runtime.AppI            = (*AppBuilder)(nil)
	_ servertypes.Application = (*AppBuilder)(nil)
)

// AppBuilder extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type AppBuilder struct {
	*app.IrisApp
}

func (ab AppBuilder) build(
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *app.IrisApp {
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
	return app
}

// MakeCodecs returns the application codec and tx codec
func MakeCodecs() params.EncodingConfig {
	return app.RegisterEncodingConfig()
}

// DefaultGenesis returns default genesis state as raw bytes
func DefaultGenesis(cdc codec.JSONCodec) map[string]json.RawMessage {
	return app.ModuleBasics.DefaultGenesis(cdc)
}
