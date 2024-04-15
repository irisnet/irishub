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
)

var (
	//go:embed compiled_contracts/Token.json
	ERC20JSON []byte //nolint: golint

	// ERC20TokenContract is the compiled erc20 contract
	ERC20TokenContract types.CompiledContract
)

func init() {
	err := json.Unmarshal(ERC20JSON, &ERC20TokenContract)
	if err != nil {
		panic(err)
	}

	if len(ERC20TokenContract.Bin) == 0 {
		panic("load contract failed")
	}
}
