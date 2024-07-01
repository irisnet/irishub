package contracts

import (
	_ "embed" //embed

	"encoding/json"

	"mods.irisnet.org/modules/token/types"
)

const (
	// EventSwapToNative is the name of the event emitted when swapping to native
	EventSwapToNative = "SwapToNative"
	// MethodMint is the name of the method mint of the ERC20 contract
	MethodMint        = "mint"
	// MethodBurn is the name of the method burn of the ERC20 contract
	MethodBurn        = "burn"
	// MethodBalanceOf is the name of the method balanceOf of the ERC20 contract
	MethodBalanceOf   = "balanceOf"
	// MethodInitialize is the name of the method initialize of the ERC20 contract
	MethodInitialize  = "initialize"
	// MethodUpgradeTo is the name of the method upgradeTo of the Beacon contract
	MethodUpgradeTo   = "upgradeTo"
)

var (
	// ERC20JSON is the json abi of the ERC20 contract
	//go:embed compiled_contracts/Token.json
	ERC20JSON []byte //nolint: golint // Generate by solc

	// TokenProxyJSON is the json abi of the ERC20 contract
	//go:embed compiled_contracts/TokenProxy.json
	TokenProxyJSON []byte //nolint: golint // Generate by solc

	// BeaconJSON is the json abi of the Beacon contract
	//go:embed compiled_contracts/UpgradeableBeacon.json
	BeaconJSON []byte //nolint: golint // Generate by solc

	// ERC20TokenContract is the compiled erc20 contract
	ERC20TokenContract types.CompiledContract

	// TokenProxyContract is the compiled erc20 contract proxy
	TokenProxyContract types.CompiledContract

	// BeaconContract is the compiled beacon contract proxy
	BeaconContract types.CompiledContract
)

func init() {
	err := json.Unmarshal(ERC20JSON, &ERC20TokenContract)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(BeaconJSON, &BeaconContract)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(TokenProxyJSON, &TokenProxyContract)
	if err != nil {
		panic(err)
	}

	if len(TokenProxyContract.Bin) == 0 {
		panic("load token proxy contract failed")
	}
}
