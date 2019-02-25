# iriscli stake validator

## Description

Query a validator

## Usage

```
iriscli stake validator [validator-address] [flags]
```
Print help messages:
```
iriscli stake validator --help
```

## Examples

Query a validator
```
iriscli stake validator [validator-address]
```

After that, you will get the specified validator's info.

```txt
Validator
Operator Address: fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd
Validator Consensus Pubkey: fcp1zcjduepq8fnuxnceuy4n0fzfc6rvf0spx56waw67lqkrhxwsxgnf8zgk0nus2r55he
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
