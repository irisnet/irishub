<!--
order: 2
-->

# Messages

## Service Definition

The service definition can be created by any user via a `MsgDefineService`
message.

```go
type MsgDefineService struct {
    Name              string
    Description       string
    Tags              []string
    Author            sdk.AccAddress
    AuthorDescription string
    Schemas           string
}
```

**State modifications:**

- Create a new `ServiceDefinition`

This message is expected to fail if:

- `Name` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]
- the length of `Description` exceeds 128 bytes
- the length of `AuthorDescription` exceeds 128 bytes
- `Tags` does not satisfy the following:
  - length of a single tag does not exceeds 70 bytes
  - total number does not exceeds 10
  - does not contain duplicate tags
- `Schemas` does not satisfy the following:
  - is a valid JSON object
  - contains the input and output object which are both valid JSON Schema
- the service definition with the `Name` already exists

## Service Binding

Any user who wants to provide a service can create a service binding via `MsgBindService`. Later, the service binding can be updated via `MsgUpdateServiceBinding`, diabled via `MsgDisableServiceBinding`, and enabled via `MsgEnableServiceBinding`.

The owner can refund deposit from an unavailable service binding after a period of time since disabled. The operation is via `MsgRefundServiceDeposit`

An owner can set an address to withdraw fees earned by its providers. The corresponding message is `MsgSetWithdrawAddress`

```go
type MsgBindService struct {
    ServiceName string
    Provider    sdk.AccAddress
    Deposit     sdk.Coins
    Pricing     string
    QoS         uint64
    Options     string
    Owner       sdk.AccAddress
}
```

**State modifications:**

- Create a new `ServiceBinding`

This message is expected to fail if:

- `ServiceName` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]
- `Deposit` is invalid coins or not positive
- `Pricing` does not conform the `Pricing Schema`
- `QoS` is equal to 0 or greater than the system parameter `MaxRequestTimeout`
- `Options` is non-functional options
- the service definition with the `ServiceName` does not exist
- the service binding with the `ServiceName` and `Provider` already exists

```go
type MsgUpdateServiceBinding struct {
    ServiceName string
    Provider    sdk.AccAddress
    Deposit     sdk.Coins
    Pricing     string
    QoS         uint64
    Options     string
    Owner       sdk.AccAddress
}
```

**State modifications:**

- Update the deposit if provided
- Update the pricing if provided
- Update the QoS if provided

This message is expected to fail if:

- `ServiceName` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]
- `Deposit` is invalid coins when not empty
- `Pricing` does not conform the `Pricing Schema` if not empty
- the service binding with the `ServiceName` and `Provider` does not exist
- owner of the servic binding is not `Owner`

```go
type MsgDisableServiceBinding struct {
    ServiceName string
    Provider    sdk.AccAddress
    Owner       sdk.AccAddress
}
```

**State modifications:**

- Disable the service binding

This message is expected to fail if:

- `ServiceName` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]
- the service binding with the `ServiceName` and `Provider` does not exist
- owner of the servic binding is not `Owner`
- the service binding is unvailable

```go
type MsgEnableServiceBinding struct {
    ServiceName string
    Provider    sdk.AccAddress
    Deposit     sdk.Coins
    Owner       sdk.AccAddress
}
```

**State modifications:**

- Enable the service binding
- Increase the deposit by `Deposit` if provided

This message is expected to fail if:

- `ServiceName` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]
- `Deposit` is invalid coins when not empty
- the service binding with the `ServiceName` and `Provider` does not exist
- owner of the servic binding is not `Owner`
- the service binding is available

```go
type MsgRefundServiceDeposit struct {
    ServiceName string
    Provider    sdk.AccAddress
    Owner       sdk.AccAddress
}
```

**State modifications:**

- Change the deposit to zero

This message is expected to fail if:

- `ServiceName` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]
- the service binding with the `ServiceName` and `Provider` does not exist
- owner of the servic binding is not `Owner`
- the service binding is available
- the deposit is zero
- the block time is earlier than the refundable time

This message is expected to fail if:

- `ServiceName` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]

```go
type MsgSetWithdrawAddress struct {
    Owner           sdk.AccAddress
    WithdrawAddress sdk.AccAddress
}
```

**State modifications:**

- Set a new withdrawal address for the owner

This message is expected to fail if:

- `WithdrawAddress` is empty

## Service Invocation

A consumer can initiate a service invocation via `MsgCallService`, and the targeting provider can respond to the request via `MsgRespondService`. After invocation, the consumer can update the created request context via `MsgUpdateRequestContext`. The request context can be paused via `MsgPauseRequestContext`, started via `MsgStartRequestContext`, terminated via `MsgKillReqeustContext` as well. The owner of the provider can withdraw the earned fees via `MsgWithdrawEarnedFees`

```go
type MsgCallService struct {
    ServiceName       string
    Providers         []sdk.AccAddress
    Consumer          sdk.AccAddress
    Input             string
    ServiceFeeCap     sdk.Coins
    Timeout           int64
    Repeated          bool
    RepeatedFrequency uint64
    RepeatedTotal     int64
}
```

**State modifications:**

- Create a `RequestContext` by which the requests are generated every EndBlocker

This message is expected to fail if:

- `ServiceName` does not satisfy the following:
  - begins with alphabetic charactors
  - consists of only alphanumerics, - and _
  - length ranges in (0,70]
- `Providers` contain duplicate addresses
- `Input` does not conform to the service input schema
- `ServiceFeeCap` is invalid coins
- `Timeout` is equal to 0 or greater than the system parameter `MaxRequestTimeout`
- `RepeatedFrequency` is less than `Timeout` if `Repeated` is true
- `RepeatedTotal` is less than -1 if `Repeated` is true

```go
type MsgRespondService struct {
    RequestID tmbytes.HexBytes `json:"request_id"`
    Provider  sdk.AccAddress   `json:"provider"`
    Result    string           `json:"result"`
    Output    string           `json:"output"`
}
```

**State modifications:**

- Create a `Response` if succeeded
- Slash the provider and refund the service fee to the consumer if the request times out

This message is expected to fail if:

- `RequestID` is invalid
- the request is not active
- the provider of the corresponding request is not `Provider`
- `Result` does not conform to the result schema
- `Output` is not provided if the `Result` code is 200
- `Output` does not conform to the service output schema when required

```go
type MsgUpdateRequestContext struct {
    RequestContextID  tmbytes.HexBytes
    Providers         []sdk.AccAddress
    ServiceFeeCap     sdk.Coins
    Timeout           int64
    RepeatedFrequency uint64
    RepeatedTotal     int64
    Consumer          sdk.AccAddress
}
```

**State modifications:**

- Update the providers if provided
- Update the service fee cap if provided
- Update the timeout if provided
- Update the frequency if provided
- Update the total count if provided

This message is expected to fail if:

- `RequestContextID` is invalid
- the request context does not exist
- `Providers` contain duplicate addresses
- `ServiceFeeCap` is invalid coins if not empty
- `Timeout` is less than the frequency or greater than the system parameter `MaxRequestTimeout` if non zero
- `RepeatedFrequency` is less than the timeout if non zero
- `RepeatedTotal` is less than -1 if non zero

```go
type MsgPauseRequestContext struct {
    RequestContextID tmbytes.HexBytes
    Consumer         sdk.AccAddress
}
```

**State modifications:**

- Pause the request context

This message is expected to fail if:

- `RequestContextID` is invalid
- the request context does not exist
- the request context is not running

```go
type MsgStartRequestContext struct {
    RequestContextID tmbytes.HexBytes
    Consumer         sdk.AccAddress
}
```

**State modifications:**

- Start the request context

This message is expected to fail if:

- `RequestContextID` is invalid
- the request context does not exist
- the request context is not paused

```go
type MsgKillRequestContext struct {
    RequestContextID tmbytes.HexBytes
    Consumer         sdk.AccAddress
}
```

**State modifications:**

- Terminate the request context

This message is expected to fail if:

- `RequestContextID` is invalid
- the request context does not exist

```go
type MsgWithdrawEarnedFees struct {
    Owner    sdk.AccAddress
    Provider sdk.AccAddress
}
```

**State modifications:**

- Change the earned fees

This message is expected to fail if:

- `Owner` does not have the earned fees or being zero
- `Provider` does not have the earned fees or being zero
