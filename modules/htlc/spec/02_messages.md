<!--
order: 2
-->

# Messages

## MsgCreateHTLC

The HTLC can be created using the `MsgCreateHTLC` message

```go
type MsgCreateHTLC struct {
    Sender               string
    To                   string
    ReceiverOnOtherChain string
    Amount               sdk.Coins
    HashLock             string
    Timestamp            uint64
    TimeLock             uint64
}
```

## MsgClaimHTLC

The HTLC can be claimed using the `MsgClaimHTLC` message

```go
type MsgClaimHTLC struct {
    Sender   string
    HashLock string
    Secret   string
}
```

## MsgRefundHTLC

The HTLC can be refunded using the `MsgRefundHTLC` message

```go
type MsgRefundHTLC struct {
    Sender   string
    HashLock string
}
```
