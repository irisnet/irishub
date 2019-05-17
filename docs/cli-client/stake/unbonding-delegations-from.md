# iriscli stake unbonding-delegations-from

## Description

Query all unbonding delegatations from a validator

## Usage
```
iriscli stake unbonding-delegations-from <validator-address> <flags>
```

Print help messages:
```
iriscli stake unbonding-delegations-from --help
```

## Examples

Query all unbonding delegatations from a validator
```
iriscli stake unbonding-delegations-from <validator-address> 
```

After that, you will get unbonding delegation's detailed info from specified validator.

```
Unbonding Delegation
Delegator: iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
Validator: iva15grv3xg3ekxh9xrf79zd0w077krgv5xfzzunhs
Creation height: 1310
Min time to unbond (unix): 2018-11-15 06:24:22.754703377 +0000 UTC
Expected balance: 0.02iris
```
