<!--
order: 1
-->

# State

## Token

Definition of data structure of FungibleToken

- Token: `0x1 -> amino(Token)`

```go
type Token struct {
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

## Params

Params is a module-wide configuration structure that stores system parameters
and defines overall functioning of the token module.

- Params: `Paramsspace("token") -> amino(params)`

```go
type Params struct {
  TokenTaxRate      sdk.Dec
  IssueTokenBaseFee sdk.Coin
  MintTokenFeeRatio sdk.Dec
}
```
