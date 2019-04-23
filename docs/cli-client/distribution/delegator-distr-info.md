# iriscli distribution delegator-distr-info

## Description

Query a delegator distribution information

## Usage

```
iriscli distribution delegator-distr-info <delegator_address> <flags>
```

Print help messages:
```
iriscli distribution delegator-distr-info --help
```

## Examples

```
iriscli distribution delegator-distr-info <delegator_address> <flags>

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