<!--
order: 1
-->

# State

## Record

`Record` defines the data structure of Record

```go
type Record struct {
    TxHash   string    
    Contents []Content 
    Creator  string
}
```

`Content` defines the data structure of Content

```go
type Content struct {
    Digest     string
    DigestAlgo string
    URI        string
    Meta       string
}
```
