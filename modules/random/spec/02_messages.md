<!--
order: 2
-->

# Messages

## MsgRequestRandom

The random can be requested using the `MsgRequestRandom` message

```go
type MsgRequestRandom struct {
    BlockInterval uint64
    Consumer      string
    Oracle        bool
    ServiceFeeCap sdk.Coins
}
```
