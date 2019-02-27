# iriscli stake redelegations

## Description

Query all redelegations records for one delegator

## Usage

```
iriscli stake redelegations [delegator-address] [flags]
```
Print help messages:
```
iriscli stake redelegations --help
```

## Examples

Query all redelegations records
```
iriscli stake redelegations [delegator-address]
```

After that, you will get all redelegations records' info for specified delegator

```json
[
  {
    "delegator_addr": "iaa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9sytjfk",
    "validator_src_addr": "fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2",
    "validator_dst_addr": "fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd",
    "creation_height": "1130",
    "min_time": "2018-11-16T07:22:48.740311064Z",
    "initial_balance": "0.1iris",
    "balance": "0.1iris",
    "shares_src": "0.1000000000",
    "shares_dst": "0.1000000000"
  }
]
```
