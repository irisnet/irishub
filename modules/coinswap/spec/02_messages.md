<!--
order: 2
-->

# Messages

## MsgSwapOrder

The coins can be swapped using the `MsgSwapOrder` message

```go
type MsgSwapOrder struct {
    Input      Input
    Output     Output
    Deadline   int64
    IsBuyOrder bool
}
```

```go
type Input struct {
    Address string
    Coin    types.Coin
}
```

```go
type Output struct {
    Address string
    Coin    types.Coin
}

```

## MsgAddLiquidity

The liquidity can be added using the `MsgAddLiquidity` message

```go
type MsgAddLiquidity struct {
    MaxToken         types.Coin
    ExactStandardAmt sdkmath.Int
    MinLiquidity     sdkmath.Int
    Deadline         int64
    Sender           string
}
```

## MsgAddUnilateralLiquidity

The liquidity can be added unilaterally using the `MsgAddUnilateralLiquidity` message

```go
type MsgAddUnilateralLiquidity struct {
	CounterpartyDenom  string
	ExactToken         types.Coin
	MinLiquidity       sdkmath.Int
	Deadline           int64
	Sender             string
}
```

## MsgRemoveLiquidity

The liquidity can be removed using the `MsgAddLiquidity` message

```go
type MsgRemoveLiquidity struct {
    WithdrawLiquidity types.Coin
    MinToken          sdkmath.Int
    MinStandardAmt    sdkmath.Int
    Deadline          int64
    Sender            string
}
```

## MsgRemoveUnilateralLiquidity

The liquidity can be removed unilaterally using the `MsgRemoveUnilateralLiquidity` message

```go
type MsgRemoveUnilateralLiquidity struct {
    CounterpartyDenom  string
    MinToken           types.Coin
    ExactLiquidity     sdkmath.Int
    Deadline           int64
    Sender             string
}
```