# iriscli distribution validator-distr-info

## Description

Query a validator distribution information

## Usage

```
iriscli distribution validator-distr-info [validator address] [flags]
```

Print help messages:
```
iriscli distribution validator-distr-info --help
```

## Examples

```
iriscli distribution validator-distr-info [validator address]
```
Example response:
```json
[
  {
    "delegator_addr": "iaa1ezzh0humhy3329xg4avhcjtay985nll06lgq50",
    "val_operator_addr": "iva14a70gzu0v2w8dlfx462c9sldvja24qazzr2ds4",
    "del_pool_withdrawal_height": "10859"
  },
  {
    "delegator_addr": "iaa1ezzh0humhy3329xg4avhcjtay985nll06lgq50",
    "val_operator_addr": "iva1ezzh0humhy3329xg4avhcjtay985nll00wz0fg",
    "del_pool_withdrawal_height": "4044"
  }
]
```