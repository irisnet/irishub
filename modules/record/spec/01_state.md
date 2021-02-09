<!--
order: 1
-->

# State

## Record

Definition of data structure of Record

- Record: `0x01 -> amino(Record)`

```go
type Record struct {
    TxHash   string    
    Contents []Content 
    Creator  string    // the creator of the record
}
```

