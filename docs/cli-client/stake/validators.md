# iriscli stake validators

## Description

Query for all validators

## Usage

```
iriscli stake validators [flags]
```
Print help messages:
```
iriscli stake validators --help
```

## Examples

Query a validator
```
iriscli stake validators
```

After that, you will get all validators' info.

```txt
Validator
Operator Address: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Validator Consensus Pubkey: fvp1zcjduepq47906n2zvq597vwyqdc0h35ve0fcl64hwqs9xw5fg67zj4g658aqyuhepj
Jailed: false
Status: Bonded
Tokens: 100.0000000000
Delegator Shares: 100.0000000000
Description: {node0   }
Bond Height: 0
Unbonding Height: 0
Minimum Unbonding Time: 1970-01-01 00:00:00 +0000 UTC
Commission: {{0.0000000000 0.0000000000 0.0000000000 0001-01-01 00:00:00 +0000 UTC}}
```
