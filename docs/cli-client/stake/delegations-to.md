# iriscli stake delegations-to

## Description

Query all delegations to one validator

## Usage

```
iriscli stake delegations-to <validator-address> <flags>
```

Print help messages:
```
iriscli stake delegations-to --help
```

## Examples

Query all delegations to one validator
```
iriscli stake delegations-to iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
```

After that, you will get all detailed info of delegations from the specified delegator address.

```
Delegation:
  Delegator:  iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
  Validator:  iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
  Shares:     100.0000000000000000000000000000
  Height:     0
Delegation:
  Delegator:  iaa1td4xnefkthfs6jg469x33shzf578fed6n7k7ua
  Validator:  iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
  Shares:     1.0000000000000000000000000000
  Height:     26
```
