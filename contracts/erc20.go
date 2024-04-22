package contracts

import (
	_ "embed"
	"encoding/json"

	"github.com/irisnet/irismod/types"
)

const (
	EventSwapToNative = "SwapToNative"
	MethodMint        = "mint"
	MethodBurn        = "burn"
	MethodBalanceOf   = "balanceOf"
	MethodInitialize   = "initialize"
)

var (
	//go:embed compiled_contracts/Token.json
	ERC20JSON []byte //nolint: golint

	//go:embed compiled_contracts/TokenProxy.json
	TokenProxyJSON []byte //nolint: golint

	// ERC20TokenContract is the compiled erc20 contract
	ERC20TokenContract types.CompiledContract

	// TokenProxyContract is the compiled erc20 contract proxy
	TokenProxyContract types.CompiledContract
)

func init() {
	err := json.Unmarshal(ERC20JSON, &ERC20TokenContract)
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
