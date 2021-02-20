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
    ExactStandardAmt sdk.Int
    MinLiquidity     sdk.Int
    Deadline         int64
    Sender           string
}
```

## MsgRemoveLiquidity

The liquidity can be removed using the `MsgAddLiquidity` message

```go
type MsgRemoveLiquidity struct {
    WithdrawLiquidity types.Coin
    MinToken          sdk.Int
    MinStandardAmt    sdk.Int
    Deadline          int64
    Sender            string
}
```
