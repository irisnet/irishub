# State

## MT

Nft defines the tokenData of non-fungible tokens, mainly including ID, owner, and tokenURI.Nft can be transferred through `MsgTransferMT`, or you can edit `tokenURI` information through `MsgEditMT` transaction. The name of the collection and the id of mt identify the unique assets in the system. The `MT` Interface inherits the MT struct and includes getter functions for the asset data. It also includes a Stringer function in order to print the struct. The interface may change if tokenData is moved to it’s own module as it might no longer be necessary for the flexibility of an interface.

```go
// MT multi token interface
type MT interface {
    GetID() string              // unique identifier of the MT
    GetName() string            // return the name of MT
    GetOwner() sdk.AccAddress   // gets owner account of the MT
    GetURI() string             // tokenData field: URI to retrieve the of chain tokenData of the MT
    GetURIHash() string
    GetData() string            // return the Data of MT
}
```

## Collections

As all MTs belong to a specific `Collection`, however, considering the performance issue, we did not store the structure, but used `{denomID}/{tokenID}` as the key to identify each mt ’s own collection, use `{denom}` as the key to store the number of mt in the current collection, which is convenient for statistics and query.collection is defined as follows

```go
// Collection of multi tokens
type Collection struct {
    Denom Denom     `json:"denom"`  // Denom of the collection; not exported to clients
    MTs  []MT `json:"mts"`   // MTs that belongs to a collection
}
```

## Owners

Owner is a data structure specifically designed for mt owned by statistical model owners. The ownership of an MT is set initially when an MT is minted and needs to be updated every time there's a transfer or when an MT is burned,defined as follows:

```go
// Owner of multi tokens
type Owner struct {
    Address       string            `json:"address"`
    IDCollections []IDCollection    `json:"id_collections"`
}
```

An `IDCollection` is similar to a `Collection` except instead of containing MTs it only contains an array of `MT` IDs. This saves storage by avoiding redundancy.

```go
// IDCollection of multi tokens
type IDCollection struct {
    DenomId string   `json:"denom_id"`
    TokenIds []string `json:"token_ids"`
}

```
