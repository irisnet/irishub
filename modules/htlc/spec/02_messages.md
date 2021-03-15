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
    SenderOnOtherChain   string
    Amount               sdk.Coins
    HashLock             string
    Timestamp            uint64
    TimeLock             uint64
    Transfer             bool
}
```

## MsgClaimHTLC

The HTLC can be claimed using the `MsgClaimHTLC` message

```go
type MsgClaimHTLC struct {
    Sender   string
    Id       string
    Secret   string
}
```
