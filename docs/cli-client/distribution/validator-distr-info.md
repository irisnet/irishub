# iriscli distribution validator-distr-info

## Description

Query a validator distribution information

## Usage

```
iriscli distribution validator-distr-info [flags]
```

Print help messages:
```
iriscli distribution validator-distr-info --help
```

## Unique Flags

There is no unique flag. But it requires an argument: validator address


## Examples

```
iriscli distribution validator-distr-info <validator address>
```
Example response:
```json
[
  {
    "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
    "val_operator_addr": "fva14a70gzu0v2w8dlfx462c9sldvja24qaz6vv4sg",
    "del_pool_withdrawal_height": "10859"
  },
  {
    "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
    "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
    "del_pool_withdrawal_height": "4044"
  }
]
```