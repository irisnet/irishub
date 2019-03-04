# iriscli stake delegations-to

## Description

Query all delegations to one validator

## Usage

```
iriscli stake delegations-to [validator-address] [flags]
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

```json
[
  {
    "delegator_addr": "iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm",
    "validator_addr": "iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh",
    "shares": "0.2000000000",
    "height": "290"
  }
]
```
