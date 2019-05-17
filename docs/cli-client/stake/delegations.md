# iriscli stake delegations

## Description

Query all delegations delegated from one delegator

## Usage
```
iriscli stake delegations <delegator-address> [flags]
```

Print help messages:
```
iriscli stake delegations --help
```

## Examples

Query all delegations delegated from one delegator
```
iriscli stake delegations iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh
```

After that, you will get all detailed info of delegations from the specified delegator address.

```
Delegation:
  Delegator:  faa1td4xnefkthfs6jg469x33shzf578fed6n7k7ua
  Validator:  fva1zkevgrasr5txhgsyqd7l9javln9et2d7k5yycy
  Shares:     1.0000000000000000000000000000
  Height:     26
```
