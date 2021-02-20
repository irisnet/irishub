<!--
order: 1
-->

# State

## HTLC

`HTLC` defines the struct of an HTLC

```go
type HTLC struct {
    Sender               string
    To                   string
    ReceiverOnOtherChain string
    Amount               sdk.Coins
    Secret               string
    Timestamp            uint64
    ExpirationHeight     uint64
    State                HTLCState
}
```

`HTLCState` defines the state of an HTLC

- `HTLC_STATE_OPEN` defines an open state
- `HTLC_STATE_COMPLETED` defines a completed state
- `HTLC_STATE_EXPIRED` defines an expired state
- `HTLC_STATE_REFUNDED` defines a refunded state

```go
type HTLCState int32

const (
    // HTLC_STATE_OPEN defines an open state.
    Open HTLCState = 0
    // HTLC_STATE_COMPLETED defines a completed state.
    Completed HTLCState = 1
    // HTLC_STATE_EXPIRED defines an expired state.
    Expired HTLCState = 2
    // HTLC_STATE_REFUNDED defines a refunded state.
    Refunded HTLCState = 3
)

var HTLCState_name = map[int32]string{
    0: "HTLC_STATE_OPEN",
    1: "HTLC_STATE_COMPLETED",
    2: "HTLC_STATE_EXPIRED",
    3: "HTLC_STATE_REFUNDED",
}

var HTLCState_value = map[string]int32{
    "HTLC_STATE_OPEN":      0,
    "HTLC_STATE_COMPLETED": 1,
    "HTLC_STATE_EXPIRED":   2,
    "HTLC_STATE_REFUNDED":  3,
}
```
