# iriscli stake unbonding-delegations

## Description

Query all unbonding-delegations records for one delegator

## Usage

```
iriscli stake unbonding-delegations <delegator-address> <flags>
```

Print help messages:
```
iriscli stake unbonding-delegations --help
```

## Examples

Query unbonding-delegations
```
iriscli stake unbonding-delegations <delegator-address>
```

After that, you will get all unbonding delegations' detailed info for one delegator

```
Unbonding Delegation
Delegator: iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
Validator: iva15grv3xg3ekxh9xrf79zd0w077krgv5xfzzunhs
Creation height: 1310
Min time to unbond (unix): 2018-11-15 06:24:22.754703377 +0000 UTC
Expected balance: 0.02iris
```
