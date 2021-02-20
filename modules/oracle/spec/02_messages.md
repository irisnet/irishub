<!--
order: 2
-->

# Messages

## MsgCreateFeed

The feed can be created using the `MsgCreateFeed` message

```go
type MsgCreateFeed struct {
    FeedName          string
    LatestHistory     uint64
    Description       string
    Creator           string
    ServiceName       string
    Providers         []string
    Input             string
    Timeout           int64
    ServiceFeeCap     sdk.Coins
    RepeatedFrequency uint64
    AggregateFunc     string
    ValueJsonPath     string
    ResponseThreshold uint32
}
```

## MsgStartFeed

The feed can be started using the `MsgStartFeed` message

```go
type MsgStartFeed struct {
    FeedName string
    Creator  string
}
```

## MsgPauseFeed

The feed can be paused using the `MsgPauseFeed` message

```go
type MsgPauseFeed struct {
    FeedName string
    Creator  string
}
```

## MsgEditFeed

The feed can be edited using the `MsgEditFeed` message

```go
type MsgEditFeed struct {
    FeedName          string
    Description       string
    LatestHistory     uint64
    Providers         []string
    Timeout           int64
    ServiceFeeCap     sdk.Coins
    RepeatedFrequency uint64
    ResponseThreshold uint32
    Creator           string
}
```
