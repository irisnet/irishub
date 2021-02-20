<!--
order: 1
-->

# State

## Random

`Random` defines the feed standard

```go
type Random struct {
    RequestTxHash string
    Height        int64
    Value         string
}

```

## Request

`Request` defines the random request standard

```go
type Request struct {
    Height           int64
    Consumer         string
    TxHash           string
    Oracle           bool
    ServiceFeeCap    sdk.Coins
    ServiceContextID string
}
```
