<!--
order: 1
-->

# State

## HTLC

`HTLC` defines the struct of an HTLC

```go
type HTLC struct {
    Id                   string
    Sender               string
    To                   string
    ReceiverOnOtherChain string
    SenderOnOtherChain   string
    Amount               skd.Coins
    HashLock             string
    Secret               string
    Timestamp            uint64
    ExpirationHeight     uint64
    State                HTLCState
    ClosedBlock          uint64
    Transfer             bool
    Direction            SwapDirection
}
```

`HTLCState` defines the state of an HTLC

- `HTLC_STATE_OPEN` defines an open state
- `HTLC_STATE_COMPLETED` defines a completed state
- `HTLC_STATE_REFUNDED` defines a refunded state

```go
type HTLCState int32

const (
    // HTLC_STATE_OPEN defines an open state.
    Open HTLCState = 0
    // HTLC_STATE_COMPLETED defines a completed state.
    Completed HTLCState = 1
    // HTLC_STATE_REFUNDED defines a refunded state.
    Refunded HTLCState = 2
)

var HTLCState_name = map[int32]string{
    0: "HTLC_STATE_OPEN",
    1: "HTLC_STATE_COMPLETED",
    2: "HTLC_STATE_REFUNDED",
}

var HTLCState_value = map[string]int32{
    "HTLC_STATE_OPEN":      0,
    "HTLC_STATE_COMPLETED": 1,
    "HTLC_STATE_REFUNDED":  2,
}
```

`SwapDirection` defines the direction of an HTLT

- `NONE` defines an htlt invalid direction
- `INCOMING` defines an htlt incoming direction
- `OUTGOING` defines an htlt outgoing direction

```go
type SwapDirection int32

const (
    // NONE defines an htlt invalid direction.
    None SwapDirection = 0
    // INCOMING defines an htlt incoming direction.
    Incoming SwapDirection = 1
    // OUTGOING defines an htlt outgoing direction.
    Outgoing SwapDirection = 2
)

var SwapDirection_name = map[int32]string{
    0: "NONE",
    1: "INCOMING",
    2: "OUTGOING",
}

var SwapDirection_value = map[string]int32{
    "NONE":  0,
    "INCOMING": 1,
    "OUTGOING": 2,
}
```

## AssetSupply

AssetSupply contains information about an asset's supply

```go
type AssetSupply struct {
    IncomingSupply           sdk.Coin
    OutgoingSupply           sdk.Coin
    CurrentSupply            sdk.Coin
    TimeLimitedCurrentSupply sdk.Coin
    TimeElapsed              time.Duration
}
```
