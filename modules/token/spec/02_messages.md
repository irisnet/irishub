<!--
order: 2
-->

# Messages

## MsgIssueToken

A token is created using the `MsgIssueToken` message.

```go
type MsgIssueToken struct {
    Symbol        string
    Name          string
    Scale         uint32
    MinUnit       string
    InitialSupply uint64
    MaxSupply     uint64
    Mintable      bool
    Owner         string
}
```

## MsgEditToken

The `MaxSupply`, `Mintable` , `Name` of a token can be updated using the `MsgEditToken`

```go
type MsgEditToken struct {
    Symbol    string
    Owner     string
    MaxSupply uint64
    Mintable  Bool
    Name      string
}
```

## MsgMintToken

The owner of the token can mint some tokens to the specified account

```go
type MsgMintToken struct {
    Symbol string
    Owner  string
    To     string
    Amount uint64
```

## MsgBurnToken

The owner of the token can mint some tokens to the specified account

```go
type MsgBurnToken struct {
    Symbol string
    Sender string
    Amount uint64
```

## MsgTransferTokenOwner

The ownership of the `token` can be transferred to others

```go
type MsgTransferTokenOwner struct {
    SrcOwner string
    DstOwner string
    Symbol   string
}
```
