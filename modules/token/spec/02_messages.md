<!--
order: 2
-->

# Messages

In this section we describe the processing of the token messages and the corresponding updates to the state.

## MsgIssueToken

A token is created using the `MsgIssueToken` message.

```go
type MsgIssueToken struct {
  Symbol        string
  Name          string
  Scale         uint8
  MinUnit       string
  InitialSupply uint64
  MaxSupply     uint64
  Mintable      bool
  Owner         sdk.AccAddress
}
```

This message is expected to fail if:

- the `Symbol` of the token is faulty, namely:
  - is not begin with `[a-zA-Z]`
  - contains characters other than letters and numbers
  - character length is greater than 8 bits or less than 3 bits
  - this symbol is already registered
- the `Name` of the token is faulty, namely:
  - is not begin with `[a-zA-Z]`
  - contains characters other than letters and numbers
  - character length exceeds 32 bits
- the `Scale` > 18
- the `MinUnit` of the token is faulty, namely:
  - is not begin with `[a-zA-Z]`
  - contains characters other than letters and numbers
  - the length is not between 3 and 10
  - this minUnit is already registered
- the `InitialSupply` is greater than `100000000000`
- the `MaxSupply` > `1000000000000` or `MaxSupply` < `InitialSupply`

This message creates and stores the `Token` object at appropriate indexes.

## MsgEditToken

The `MaxSupply`, `Mintable` , `Name` of a token can be updated using the
`MsgEditToken`.  

```go
type MsgEditToken struct {
  Symbol    string
  Owner     sdk.AccAddress
  MaxSupply uint64
  Mintable  Bool
  Name      string
}
```

This message is expected to fail if:

- the `Symbol` is not existed
- the `MaxSupply` > `1000000000000`
- the `Owner` is not the token owner
- the `Name` of the token is faulty, namely:
  - is not begin with `[a-zA-Z]`
  - contains characters other than letters and numbers
  - character length exceeds 32 bits

This message stores the updated `Token` object.

## MsgMintToken

The owner of the token can mint some tokens to the specified account

```go
type MsgMintToken struct {
  Symbol string
  Owner  sdk.AccAddress
  To     sdk.AccAddress
  Amount uint64
```

This message is expected to fail if:

- the `Symbol` is not existed
- the `Mintable` of the token is false
- the `Owner` is not the token owner
- the `Amount` `Coin` has exceeded the number of additional issuances（**MaxSupply - Issued**）

## MsgTransferTokenOwner

The ownership of the `token` can be transferred to others

```go
type MsgTransferTokenOwner struct {
  SrcOwner sdk.AccAddress
  DstOwner sdk.AccAddress
  Symbol   string
}
```

This message is expected to fail if:

- the token is not existed
- the `Owner` is not the token owner
