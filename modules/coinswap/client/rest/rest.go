package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes registers asset-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type AddLiquidityReq struct {
	BaseTx           rest.BaseReq `json:"base_tx" yaml:"base_tx"`
	ID               string       `json:"id" yaml:"id"`                                 // the unique liquidity id
	MaxToken         string       `json:"max_token" yaml:"max_token"`                   // token to be deposited as liquidity with an upper bound for its amount
	ExactStandardAmt string       `json:"exact_standard_amt" yaml:"exact_standard_amt"` // exact amount of standard token being add to the liquidity pool
	MinLiquidity     string       `json:"min_liquidity" yaml:"min_liquidity"`           // lower bound UNI sender is willing to accept for deposited coins
	Deadline         string       `json:"deadline" yaml:"deadline"`                     // deadline duration, e.g. 10m
	Sender           string       `json:"sender" yaml:"sender"`
}

type RemoveLiquidityReq struct {
	BaseTx            rest.BaseReq `json:"base_tx" yaml:"base_tx"`
	ID                string       `json:"id" yaml:"id"`                                 // the unique liquidity id
	MinToken          string       `json:"min_token" yaml:"min_token"`                   // coin to be withdrawn with a lower bound for its amount
	WithdrawLiquidity string       `json:"withdraw_liquidity" yaml:"withdraw_liquidity"` // amount of UNI to be burned to withdraw liquidity from a reserve pool
	MinStandardAmt    string       `json:"min_standard_amt" yaml:"min_standard_amt"`     // minimum amount of the native asset the sender is willing to accept
	Deadline          string       `json:"deadline" yaml:"deadline"`                     // deadline duration, e.g. 10m
	Sender            string       `json:"sender" yaml:"sender"`
}

type Input struct {
	Address string   `json:"address" yaml:"address"`
	Coin    sdk.Coin `json:"coin" yaml:"coin"`
}

type Output struct {
	Address string   `json:"address" yaml:"address"`
	Coin    sdk.Coin `json:"coin" yaml:"coin"`
}

type SwapOrderReq struct {
	BaseTx   rest.BaseReq `json:"base_tx" yaml:"base_tx"`
	Input    Input        `json:"input" yaml:"input"`       // the amount the sender is trading
	Output   Output       `json:"output" yaml:"output"`     // the amount the sender is receiving
	Deadline string       `json:"deadline" yaml:"deadline"` // deadline for the transaction to still be considered valid
}
