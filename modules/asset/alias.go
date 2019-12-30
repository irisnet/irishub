package asset

import (
	"github.com/irisnet/irishub/modules/asset/keeper"
	"github.com/irisnet/irishub/modules/asset/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

var (
	// functions aliases
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)

type (
	Keeper = keeper.Keeper
)
