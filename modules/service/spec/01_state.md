<!--
order: 1
-->

# State

## Parameters

`Parameters` define the rules on which the service module depends to
guarantee the interaction between the consumer and provider.

```go
type Params struct {
    MaxRequestTimeout         int64         // maximum request timeout
    MinDepositMultiple        int64         // minimum deposit multiple
    MinDeposit                sdk.Coins     // minimum deposit
    ServiceFeeTax             sdk.Dec       // service fee tax ratio
    SlashFraction             sdk.Dec       // fraction for stashing
    ComplaintRetrospect       time.Duration // duration for complaint retrospect
    ArbitrationTimeLimit      time.Duration // duration for arbitration
    TxSizeLimit               uint64        // transaction size limitation for service
    BaseDenom                 string        // base denom for deposit
    RestrictedServiceFeeDenom bool          // indicates if the service fee only accepts the base denom
}
```

Parameters are stored in a global GlobalParams KVStore.

## ServiceDefinition

`ServiceDefinition` represents a service definiton object

```go
type ServiceDefinition struct {
    Name              string         // service name
    Description       string         // service description
    Tags              []string       // service tags
    Author            string         // service creator
    AuthorDescription string         // author description
    Schemas           string         // service interface schemas
}
```

## ServiceBinding

`ServiceBinding` is intended for storing the service binding

```go
type ServiceBinding struct {
    ServiceName  string            // service name
    Provider     string            // provider address
    Deposit      sdk.Coins         // deposit for the binding
    Pricing      string            // service pricing
    QoS          uint64            // service quality in terms of minimum response time
    Options      string            // non-functional requirements options
    Available    bool              // indicate if the binding is active
    DisabledTime time.Time         // time when the binding is disabled
    Owner        string            // owner of the binding
}
```

## ServiceInvocation

The `RequestContext` object represents a basic context in which the
requests are initiated.

```go
type RequestContext struct {
    ServiceName            string                   // service name
    Providers              []string                 // provider address list
    Consumer               string                   // consumer address
    ServiceFeeCap          sdk.Coins                // maximum service fee to pay for a single request
    Input                  string                   // service input data conforming to the service input schema
    ModuleName             string                   // name of the module from which the invocation is initiated, which is not necessary from CLI and API
    Timeout                int64                    // request timeout
    RepeatedFrequency      uint64                   // invocation frequency when the request context is repeated
    RepeatedTotal          int64                    // invocation total number when the request context is repeated
    BatchCounter           uint64                   // the current batch number
    BatchRequestCount      uint32                   // request count for the current batch
    BatchResponseCount     uint32                   // response count for the current batch
    BatchResponseThreshold uint32                   // response threshold for the current batch
    ResponseThreshold      uint32                   // initial response threshold for the request context
    Repeated               bool                     // indicate if the request context is repetitive
    BatchState             RequestContextBatchState // state for the current batch
    State                  RequestContextState      // state for the request context
}
```

`CompactRequest` is used to store the compact request which contains the
ID of the request context to which the request belongs.

```go
type CompactRequest struct {
    RequestContextId           string  // ID of the request context from which the request is initiated
    RequestContextBatchCounter uint64            // the batch number of the request
    Provider                   string    // provider address
    ServiceFee                 sdk.Coins         // service fee
    RequestHeight              int64             // block number at which the request is initiated
  	ExpirationHeight           int64             // block number at which the request is expired
}
```

`Response` is an object which is a response to a request targeting the
provider

```go
type Response struct {
    Provider                   string            // provicer address
    Consumer                   string            // consumer address
    Result                     string            // response result according with the result schema
    Output                     string           // response output according with the service output schema
    RequestContextID           string            // ID of the request context to which the response belongs
    RequestContextBatchCounter uint64           // the batch number of the response
}
```

## Stores

_Stores are KVStores in the multi-store.

For pseudocode purposes, here are the two function we will use to read
or write in stores:

- `load(StoreKey, Key)`: Retrieve item stored at key `Key` in store
  found at key `StoreKey` in the multistore
- `store(StoreKey, Key, value)`: Write value `Value` at key `Key` in
  store found at key `StoreKey` in the multistore

