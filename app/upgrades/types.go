package upgrades

import (
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/irisnet/irishub/v3/app/keepers"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
// must have written, in order for the state migration to go smoothly.
// An upgrade must implement this struct, and then set it in the app.go.
// The app.go will then define the handler.
type Upgrade struct {
	// Upgrade version name, for the upgrade handler, e.g. `v7`
	UpgradeName string

	// UpgradeHandlerConstructor defines the function that creates an upgrade handler
	UpgradeHandlerConstructor func(*module.Manager, module.Configurator, Toolbox) upgradetypes.UpgradeHandler

	// Store upgrades, should be used for any new modules introduced, new modules deleted, or store names renamed.
	StoreUpgrades *store.StoreUpgrades
}

// ConsensusParamsReaderWriter defines the interface for reading and writing consensus params
type ConsensusParamsReaderWriter interface {
	StoreConsensusParams(ctx sdk.Context, cp *tmproto.ConsensusParams)
	GetConsensusParams(ctx sdk.Context) *tmproto.ConsensusParams
}

// Toolbox contains all the modules necessary for an upgrade
type Toolbox struct {
	AppCodec      codec.Codec
	ModuleManager *module.Manager
	ReaderWriter  ConsensusParamsReaderWriter
	keepers.AppKeepers
}

type upgradeRouter struct {
	mu map[string]Upgrade
}

// NewUpgradeRouter creates a new upgrade router.
//
// No parameters.
// Returns a pointer to upgradeRouter.
func NewUpgradeRouter() *upgradeRouter {
	return &upgradeRouter{make(map[string]Upgrade)}
}

func (r *upgradeRouter) Register(u Upgrade) *upgradeRouter {
	if _, has := r.mu[u.UpgradeName]; has {
		panic(u.UpgradeName + " already registered")
	}
	r.mu[u.UpgradeName] = u
	return r
}

func (r *upgradeRouter) Routers() map[string]Upgrade {
	return r.mu
}

func (r *upgradeRouter) UpgradeInfo(planName string) Upgrade {
	return r.mu[planName]
}
