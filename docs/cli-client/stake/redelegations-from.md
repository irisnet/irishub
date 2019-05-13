# iriscli stake redelegations-from

## Description

Query all outgoing redelegations from a validator

## Usage

```
iriscli stake redelegations-from <validator-address> <flags>
```

Print help messages:
```
iriscli stake redelegations-from --help
```

## Examples

Query all outgoing redelegatations
```
iriscli stake redelegations-from <validator-address> 
```

After that, you will get all outgoing redelegatations' from specified validator
```
Redelegation
Delegator: iaa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9gtd2ft
Source Validator: iva1dayujdfnxjggd5ydlvvgkerp2supknth9a8qhh
Destination Validator: iva1h27xdw6t9l5jgvun76qdu45kgrx9lqedpg3ecs
Creation height: 1130
Min time to unbond (unix): 2018-11-16 07:22:48.740311064 +0000 UTC
Source shares: 0.1000000000
Destination shares: 0.1000000000
```
