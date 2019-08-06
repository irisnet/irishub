package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	//AddLiquidity
	flagDeposit   = "deposit"
	flagAmount    = "amount"
	flagMinReward = "min-reward"
	flagPeriod    = "period"

	//RemoveLiquidity
	flagMinToken  = "min-token"
	flagMinNative = "min-native"

	//SwapTokens
	flagInput      = "input"
	flagOutput     = "output"
	flagDeadline   = "deadline"
	FlagRecipient  = "recipient"
	FlagIsBuyOrder = "is-buy-order"
)

var (
	FsSwapTokens      = flag.NewFlagSet("", flag.ContinueOnError)
	FsAddLiquidity    = flag.NewFlagSet("", flag.ContinueOnError)
	FsRemoveLiquidity = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	//AddLiquidity
	FsAddLiquidity.String(flagDeposit, "", "Token to be deposit")
	FsAddLiquidity.String(flagAmount, "", "Amount of iris token")
	FsAddLiquidity.String(flagMinReward, "", "lower bound UNI sender is willing to accept for deposited coins")
	FsAddLiquidity.String(flagPeriod, "", "deadline of transaction")

	//RemoveLiquidity
	FsRemoveLiquidity.String(flagMinToken, "", "coin to be withdrawn with a lower bound for its amount")
	FsRemoveLiquidity.String(flagMinNative, "", "minimum amount of the native asset the sender is willing to accept")
	FsRemoveLiquidity.String(flagAmount, "", "amount of UNI to be burned to withdraw liquidity from a reserve pool")
	FsRemoveLiquidity.String(flagPeriod, "", "deadline of transaction")

	//SwapTokens
	FsSwapTokens.String(flagInput, "", "Amount of coins to swap, if --is-buy-order=true,it means this is the exact amount of coins sold,otherwise it means this is the max amount of coins sold")
	FsSwapTokens.String(flagOutput, "", "Amount of coins to swap, if --is-buy-order=true,it means this is the min amount of coins bought,otherwise it means this is the exact amount of coins bought")
	FsSwapTokens.String(flagDeadline, "", "deadline for the transaction to still be considered valid,such as 6m")
	FsSwapTokens.String(FlagRecipient, "", "address of coins bought")
	FsSwapTokens.Bool(FlagIsBuyOrder, false, "boolean indicating whether the order should be treated as a buy or sell")
}
