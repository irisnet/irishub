package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	//AddLiquidity
	flagDeposit   = "deposit"
	flagAmount    = "native"
	flagMinReward = "min-reward"
	flagPeriod    = "period"

	//SwapTokens
	flagInputToken  = "input-token"
	flagOutputToken = "output-token"
	flagDeadline    = "deadline"
	FlagRecipient   = "recipient"
	FlagIsBuyOrder  = "is-buy-order"
)

var (
	FsSwapTokens   = flag.NewFlagSet("", flag.ContinueOnError)
	FsAddLiquidity = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	//AddLiquidity
	FsAddLiquidity.String(flagDeposit, "", "Token to be deposit")
	FsAddLiquidity.String(flagAmount, "", "Amount of iris token")
	FsAddLiquidity.String(flagMinReward, "", "lower bound UNI sender is willing to accept for deposited coins")
	FsAddLiquidity.String(flagPeriod, "", "deadline of transaction")

	//SwapTokens
	FsSwapTokens.String(flagInputToken, "", "Amount of coins to swap, if --is-buy-order=true,it means this is the exact amount of coins sold,otherwise it means this is the max amount of coins sold")
	FsSwapTokens.String(flagOutputToken, "", "Amount of coins to swap, if --is-buy-order=true,it means this is the min amount of coins bought,otherwise it means this is the exact amount of coins bought")
	FsSwapTokens.String(flagDeadline, "", "deadline for the transaction to still be considered valid,such as 6m")
	FsSwapTokens.String(FlagRecipient, "", "address of coins bought")
	FsSwapTokens.Bool(FlagIsBuyOrder, false, "boolean indicating whether the order should be treated as a buy or sell")
}
