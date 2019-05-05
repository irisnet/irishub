# iriscli distribution validator-distr-info

## Description

Query a validator distribution information

## Usage

```
iriscli distribution validator-distr-info <validator_address> <flags>
```

Print help messages:
```
iriscli distribution validator-distr-info --help
```

## Examples

```
iriscli distribution validator-distr-info <validator_address>
```

Example response:
```json
{
  "operator_addr": "iva1e7wljxhz7u7xrh63xjlds8vcy047a47ejpnz7a",
  "fee_pool_withdrawal_height": "101290",
  "del_accum": {
    "update_height": "101290",
    "accum": "0.0000000000"
  },
  "del_pool": "0.0000000000000000000000000000iris",
  "val_commission": "12.8560369893449408111336573478iris"
}
```