<!--
order: 1
-->

# State

## Feed

`Feed` defines the feed standard

```go
type Feed struct {
    FeedName         string
    Description      string
    AggregateFunc    string
    ValueJsonPath    string
    LatestHistory    uint64
    RequestContextID string
    Creator          string
}
```

## FeedValue

`FeedValue` defines the feed result standard

```go
type FeedValue struct {
    Data      string
    Timestamp time.Time
}
```
