# State

## MT

MT defines the tokenData of multi token, mainly including ID, owner, and supply. MT can be transferred through `TransferMT`, or you can edit `MetaData` information through `EditMT` transaction. The name of the collection and the id of mt identify the unique assets in the system. The `MT` Interface inherits the MT struct and includes getter functions for the asset data. It also includes a Stringer function in order to print the struct. The interface may change if tokenData is moved to its own module as it might no longer be necessary for the flexibility of an interface.

```go
// MT multi token interface
type MT interface {
    GetID() string
    GetSupply() uint64
    GetData() []byte
}
```

## Collection

As all MTs belong to a specific `Collection`, however, considering the performance issue, we did not store the structure, but used `{denomID}/{tokenID}` as the key to identify each mt ’s own collection, use `{denom}` as the key to store the object array of mt in the current collection, which is convenient for statistics and query. Collection is defined as follows:

```go
// Collection defines a type of collection
type Collection struct {
    Denom *Denom `json:"denom"`
    Mts   []MT   `json:"mts"`
}
```

## Balance

Balance is a data structure specifically designed for MT owned by statistical model owners. The ownership of an MT is set initially when an MT is minted and needs to be updated every time there's a transfer or when an MT is burned,defined as follows:

```go
// Balance defines multi token balance for owners
type Balance struct {
    MtId   string `json:"mt_id"`
    Amount uint64 `json:"amount"`
}
```
