# iriscli distribution delegation-distr-info

## introduction

Query a delegation distribution information

## Usage

```
iriscli distribution delegation-distr-info [flags]
```

Print help messages:
```
iriscli distribution delegation-distr-info --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator | string | true     | ""       | Bech address of the validator |
| --address-delegator | string | true     | ""       | Bech address of the delegator |

## Examples

```
iriscli distribution delegation-distr-info --address-delegator=<delegator address> --address-validator=<validator address>
```
Example response:
```json
{
  "delegator_addr": "iaa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
  "val_operator_addr": "iva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
  "del_pool_withdrawal_height": "4044"
}
```