<!--
order: 2
-->

# Messages

## MsgDefineService

The service definition can be created by any user via a `MsgDefineService` message.

```go
type MsgDefineService struct {
    Name              string
    Description       string
    Tags              []string
    Author            string
    AuthorDescription string
    Schemas           string
}
```

## MsgBindService

Any user who wants to provide a service can create a service binding via `MsgBindService`.

```go
type MsgBindService struct {
    ServiceName string
    Provider    string
    Deposit     sdk.Coins
    Pricing     string
    QoS         uint64
    Options     string
    Owner       string
}
```

## MsgUpdateServiceBinding

The service binding can be updated via `MsgUpdateServiceBinding`

```go
type MsgUpdateServiceBinding struct {
    ServiceName string
    Provider    string
    Deposit     sdk.Coins
    Pricing     string
    QoS         uint64
    Options     string
    Owner       string
}
```

## MsgDisableServiceBinding

The service binding can be diabled via `MsgDisableServiceBinding`

```go
type MsgDisableServiceBinding struct {
    ServiceName string
    Provider    string
    Owner       string
}
```

## MsgEnableServiceBinding

The service binding can be enabled via `MsgEnableServiceBinding`

```go
type MsgEnableServiceBinding struct {
    ServiceName string
    Provider    string
    Deposit     sdk.Coins
    Owner       string
}
```

## MsgRefundServiceDeposit

The owner can refund deposit from an unavailable service binding after a period of time since disabled. The operation is via `MsgRefundServiceDeposit`

```go
type MsgRefundServiceDeposit struct {
    ServiceName string
    Provider    string
    Owner       string
}
```

## MsgSetWithdrawAddress

An owner can set an address to withdraw fees earned by its providers. The corresponding message is `MsgSetWithdrawAddress`

```go
type MsgSetWithdrawAddress struct {
    Owner           string
    WithdrawAddress string
}
```

## MsgCallService

A consumer can initiate a service invocation via `MsgCallService`.

```go
type MsgCallService struct {
    ServiceName       string
    Providers         []string
    Consumer          string
    Input             string
    ServiceFeeCap     sdk.Coins
    Timeout           int64
    Repeated          bool
    RepeatedFrequency uint64
    RepeatedTotal     int64
}
```

## MsgRespondService

The targeting provider can respond to the request via `MsgRespondService`

```go
type MsgRespondService struct {
    RequestID string   `json:"request_id"`
    Provider  string   `json:"provider"`
    Result    string   `json:"result"`
    Output    string   `json:"output"`
}
```

## MsgUpdateRequestContext

After invocation, the consumer can update the created request context via `MsgUpdateRequestContext`.

```go
type MsgUpdateRequestContext struct {
    RequestContextID  string
    Providers         []string
    ServiceFeeCap     sdk.Coins
    Timeout           int64
    RepeatedFrequency uint64
    RepeatedTotal     int64
    Consumer          string
}
```

## MsgPauseRequestContext

The request context can be paused via `MsgPauseRequestContext`

```go
type MsgPauseRequestContext struct {
    RequestContextID string
    Consumer         string
}
```

## MsgStartRequestContext

The request context can be started via `MsgStartRequestContext`

```go
type MsgStartRequestContext struct {
    RequestContextID string
    Consumer         string
}
```

## MsgKillRequestContext

The request context can be terminated via `MsgKillReqeustContext`

```go
type MsgKillRequestContext struct {
    RequestContextID string
    Consumer         string
}
```

## MsgWithdrawEarnedFees

The owner of the provider can withdraw the earned fees via `MsgWithdrawEarnedFees`

```go
type MsgWithdrawEarnedFees struct {
    Owner    string
    Provider string
}
```
