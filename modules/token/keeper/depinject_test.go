package keeper_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	authzmodulev1 "cosmossdk.io/api/cosmos/authz/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	capabilitymodulev1 "cosmossdk.io/api/cosmos/capability/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	crisismodulev1 "cosmossdk.io/api/cosmos/crisis/module/v1"
	distrmodulev1 "cosmossdk.io/api/cosmos/distribution/module/v1"
	evidencemodulev1 "cosmossdk.io/api/cosmos/evidence/module/v1"
	feegrantmodulev1 "cosmossdk.io/api/cosmos/feegrant/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	govmodulev1 "cosmossdk.io/api/cosmos/gov/module/v1"
	groupmodulev1 "cosmossdk.io/api/cosmos/group/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	slashingmodulev1 "cosmossdk.io/api/cosmos/slashing/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
	upgrademodulev1 "cosmossdk.io/api/cosmos/upgrade/module/v1"
	vestingmodulev1 "cosmossdk.io/api/cosmos/vesting/module/v1"
	"cosmossdk.io/core/appconfig"
	"google.golang.org/protobuf/types/known/durationpb"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"

	tokenmodule "mods.irisnet.org/api/irismod/token/module/v1"
	_ "mods.irisnet.org/token"
	"mods.irisnet.org/token/contracts"
	tokentypes "mods.irisnet.org/token/types"
)

var (

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	genesisModuleOrder = []string{
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName,
		distrtypes.ModuleName, stakingtypes.ModuleName, slashingtypes.ModuleName, govtypes.ModuleName,
		minttypes.ModuleName, crisistypes.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName, authz.ModuleName,
		feegrant.ModuleName, group.ModuleName, paramstypes.ModuleName, upgradetypes.ModuleName,
		vestingtypes.ModuleName, consensustypes.ModuleName, tokentypes.ModuleName,
	}

	// module account permissions
	moduleAccPerms = []*authmodulev1.ModuleAccountPermission{
		{Account: authtypes.FeeCollectorName},
		{Account: distrtypes.ModuleName},
		{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
		{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
		{Account: govtypes.ModuleName, Permissions: []string{authtypes.Burner}},
		{Account: tokentypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}},
	}

	// blocked account addresses
	blockAccAddrs = []string{
		authtypes.FeeCollectorName,
		distrtypes.ModuleName,
		minttypes.ModuleName,
		stakingtypes.BondedPoolName,
		stakingtypes.NotBondedPoolName,
		// We allow the following module accounts to receive funds:
		// govtypes.ModuleName
	}

	// application configuration (used by depinject)
	AppConfig = appconfig.Compose(&appv1alpha1.Config{
		Modules: []*appv1alpha1.ModuleConfig{
			{
				Name: "runtime",
				Config: appconfig.WrapAny(&runtimev1alpha1.Module{
					AppName: "SimApp",
					// During begin block slashing happens after distr.BeginBlocker so that
					// there is nothing left over in the validator fee pool, so as to keep the
					// CanWithdrawInvariant invariant.
					// NOTE: staking module is required if HistoricalEntries param > 0
					// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
					BeginBlockers: []string{
						upgradetypes.ModuleName,
						capabilitytypes.ModuleName,
						minttypes.ModuleName,
						distrtypes.ModuleName,
						slashingtypes.ModuleName,
						evidencetypes.ModuleName,
						stakingtypes.ModuleName,
						authtypes.ModuleName,
						banktypes.ModuleName,
						govtypes.ModuleName,
						crisistypes.ModuleName,
						genutiltypes.ModuleName,
						authz.ModuleName,
						feegrant.ModuleName,
						group.ModuleName,
						paramstypes.ModuleName,
						vestingtypes.ModuleName,
						tokentypes.ModuleName,
						consensustypes.ModuleName,
					},
					EndBlockers: []string{
						crisistypes.ModuleName,
						govtypes.ModuleName,
						stakingtypes.ModuleName,
						capabilitytypes.ModuleName,
						authtypes.ModuleName,
						banktypes.ModuleName,
						distrtypes.ModuleName,
						slashingtypes.ModuleName,
						minttypes.ModuleName,
						genutiltypes.ModuleName,
						evidencetypes.ModuleName,
						authz.ModuleName,
						feegrant.ModuleName,
						group.ModuleName,
						paramstypes.ModuleName,
						consensustypes.ModuleName,
						upgradetypes.ModuleName,
						tokentypes.ModuleName,
						vestingtypes.ModuleName,
					},
					OverrideStoreKeys: []*runtimev1alpha1.StoreKeyConfig{
						{
							ModuleName: authtypes.ModuleName,
							KvStoreKey: "acc",
						},
					},
					InitGenesis: genesisModuleOrder,
					// When ExportGenesis is not specified, the export genesis module order
					// is equal to the init genesis order
					// ExportGenesis: genesisModuleOrder,
					// Uncomment if you want to set a custom migration order here.
					// OrderMigrations: nil,
				}),
			},
			{
				Name: authtypes.ModuleName,
				Config: appconfig.WrapAny(&authmodulev1.Module{
					Bech32Prefix:             "cosmos",
					ModuleAccountPermissions: moduleAccPerms,
					// By default modules authority is the governance module. This is configurable with the following:
					// Authority: "group", // A custom module authority can be set using a module name
					// Authority: "cosmos1cwwv22j5ca08ggdv9c2uky355k908694z577tv", // or a specific address
				}),
			},
			{
				Name:   vestingtypes.ModuleName,
				Config: appconfig.WrapAny(&vestingmodulev1.Module{}),
			},
			{
				Name: banktypes.ModuleName,
				Config: appconfig.WrapAny(&bankmodulev1.Module{
					BlockedModuleAccountsOverride: blockAccAddrs,
				}),
			},
			{
				Name:   stakingtypes.ModuleName,
				Config: appconfig.WrapAny(&stakingmodulev1.Module{}),
			},
			{
				Name:   slashingtypes.ModuleName,
				Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
			},
			{
				Name:   paramstypes.ModuleName,
				Config: appconfig.WrapAny(&paramsmodulev1.Module{}),
			},
			{
				Name:   "tx",
				Config: appconfig.WrapAny(&txconfigv1.Config{}),
			},
			{
				Name:   genutiltypes.ModuleName,
				Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
			},
			{
				Name:   authz.ModuleName,
				Config: appconfig.WrapAny(&authzmodulev1.Module{}),
			},
			{
				Name:   upgradetypes.ModuleName,
				Config: appconfig.WrapAny(&upgrademodulev1.Module{}),
			},
			{
				Name:   distrtypes.ModuleName,
				Config: appconfig.WrapAny(&distrmodulev1.Module{}),
			},
			{
				Name: capabilitytypes.ModuleName,
				Config: appconfig.WrapAny(&capabilitymodulev1.Module{
					SealKeeper: true,
				}),
			},
			{
				Name:   evidencetypes.ModuleName,
				Config: appconfig.WrapAny(&evidencemodulev1.Module{}),
			},
			{
				Name:   minttypes.ModuleName,
				Config: appconfig.WrapAny(&mintmodulev1.Module{}),
			},
			{
				Name: group.ModuleName,
				Config: appconfig.WrapAny(&groupmodulev1.Module{
					MaxExecutionPeriod: durationpb.New(time.Second * 1209600),
					MaxMetadataLen:     255,
				}),
			},
			{
				Name:   feegrant.ModuleName,
				Config: appconfig.WrapAny(&feegrantmodulev1.Module{}),
			},
			{
				Name:   govtypes.ModuleName,
				Config: appconfig.WrapAny(&govmodulev1.Module{}),
			},
			{
				Name:   crisistypes.ModuleName,
				Config: appconfig.WrapAny(&crisismodulev1.Module{}),
			},
			{
				Name:   consensustypes.ModuleName,
				Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
			},
			{
				Name:   tokentypes.ModuleName,
				Config: appconfig.WrapAny(&tokenmodule.Module{
					FeeCollectorName: authtypes.FeeCollectorName,
				}),
			},
		},
	})
)

var (
	_ tokentypes.EVMKeeper   = (*evm)(nil)
	_ tokentypes.ICS20Keeper = (*transferKeeper)(nil)
)


// ProvideEVMKeeper returns an instance of tokentypes.EVMKeeper.
//
// No parameters.
// Returns a tokentypes.EVMKeeper.
func ProvideEVMKeeper() tokentypes.EVMKeeper {
	return &evm{
		erc20s: make(map[common.Address]*erc20),
	}
}

// ProvideICS20Keeper returns an instance of tokentypes.ICS20Keeper.
//
// No parameters.
// Returns a tokentypes.ICS20Keeper.
func ProvideICS20Keeper() tokentypes.ICS20Keeper {
	return &transferKeeper{}
}

type evm struct {
	erc20s map[common.Address]*erc20
}

// ApplyMessage implements types.EVMKeeper.
func (e *evm) ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*tokentypes.Result, error) {
	isCreate := msg.To() == nil
	if isCreate {
		contractAddr := crypto.CreateAddress(msg.From(), msg.Nonce())

		data := msg.Data()[len(contracts.TokenProxyContract.Bin):]
		args, err := contracts.TokenProxyContract.ABI.Constructor.Inputs.Unpack(data)
		if err != nil {
			return nil, err
		}

		data = args[1].([]byte)
		data = data[4:]
		args, err = contracts.ERC20TokenContract.ABI.Methods[contracts.MethodInitialize].Inputs.Unpack(data)
		if err != nil {
			return nil, err
		}

		name, _ := args[0].(string)
		symbol, _ := args[1].(string)
		scale, _ := args[2].(uint8)
		e.erc20s[contractAddr] = &erc20{
			address: contractAddr,
			scale:   scale,
			name:    name,
			symbol:  symbol,
			balance: make(map[common.Address]*big.Int),
		}
		return &tokentypes.Result{
			Hash: contractAddr.Hex(),
		}, nil
	}

	erc20Contract, ok := e.erc20s[*msg.To()]
	if !ok {
		return nil, fmt.Errorf("erc20 contract not found")
	}
	return e.dispatch(erc20Contract, msg.Data())
}

// ChainID implements types.EVMKeeper.
func (e *evm) ChainID() *big.Int {
	return big.NewInt(16688)
}

// EstimateGas implements types.EVMKeeper.
func (e *evm) EstimateGas(ctx context.Context, req *tokentypes.EthCallRequest) (uint64, error) {
	return 3000000, nil
}

// SupportedKey implements types.EVMKeeper.
func (e *evm) SupportedKey(pubKey cryptotypes.PubKey) bool {
	return true
}

func (e *evm) dispatch(contract *erc20, data []byte) (*tokentypes.Result, error) {
	method, err := contracts.ERC20TokenContract.ABI.MethodById(data[0:4])
	if err != nil {
		return nil, err
	}

	ret, err := contract.call(method, data[4:])
	if err != nil {
		return nil, err
	}
	return &tokentypes.Result{
		Hash: contract.address.Hex(),
		Ret:  ret,
	}, nil
}

type erc20 struct {
	address      common.Address
	scale        uint8
	name, symbol string

	balance map[common.Address]*big.Int
}

func (erc20 erc20) call(method *abi.Method, data []byte) ([]byte, error) {
	args, err := method.Inputs.Unpack(data)
	if err != nil {
		return nil, err
	}

	switch method.Name {
	case "name":
		return method.Outputs.Pack(erc20.name)
	case "symbol":
		return method.Outputs.Pack(erc20.symbol)
	case "decimals":
		return method.Outputs.Pack(erc20.scale)
	case "balanceOf":
		balance, ok := erc20.balance[args[0].(common.Address)]
		if !ok {
			return method.Outputs.Pack(big.NewInt(0))
		}
		return method.Outputs.Pack(balance)
	case "mint":
		to := args[0].(common.Address)
		balance, ok := erc20.balance[args[0].(common.Address)]
		if !ok {
			balance = big.NewInt(0)
		}
		erc20.balance[to] = new(big.Int).Add(balance, args[1].(*big.Int))
		return nil, nil
	case "burn":
		from := args[0].(common.Address)
		erc20.balance[from] = new(big.Int).Sub(erc20.balance[from], args[1].(*big.Int))
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown method %s", method.Name)
	}
}

type transferKeeper struct{}

// HasTrace implements types.ICS20Keeper.
func (t *transferKeeper) HasTrace(ctx sdk.Context, denom string) bool {
	return true
}
