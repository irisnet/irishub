# iriscli stake redelegations-from

## Description

Query all outgoing redelegatations from a validator

## Usage

```
iriscli stake redelegations-from [validator-address] [flags]
```
Print help messages:
```
iriscli stake redelegations-from --help
```

## Examples

Query all outgoing redelegatations
```
iriscli stake redelegations-from [validator-address]
```

After that, you will get all outgoing redelegatations' from specified validator

```json
[
  {
    "delegator_addr": "iaa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9gtd2ft",
    "validator_src_addr": "fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2",
    "validator_dst_addr": "fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd",
    "creation_height": "1130",
    "min_time": "2018-11-16T07:22:48.740311064Z",
    "initial_balance": {
      "denom": "iris-atto",
      "amount": "100000000000000000"
    },
    "balance": {
      "denom": "iris-atto",
      "amount": "100000000000000000"
    },
    "shares_src": "100000000000000000.0000000000",
    "shares_dst": "100000000000000000.0000000000"
  }
]
```
