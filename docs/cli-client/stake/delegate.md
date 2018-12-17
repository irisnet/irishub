# iriscli stake delegate

## Introduction

Delegate tokens to a validator

## Usage

```
iriscli stake delegate [flags]
```

Print help messages:
```
iriscli stake delegate --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-delegator | string | true     | ""       | Bech address of the delegator |
| --amount            | string | true     | ""       | Amount of coins to bond |

## Examples

```
iriscli stake delegate --chain-id=test-irishub --from=KeyName --fee=0.04iris  --amount=10iris --address-validator=fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd
```
Output:
```json
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 16462,
   "codespace": "",
   "tags": {
     "action": "delegate",
     "delegator": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2",
     "destination-validator": "fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd"
   }
 })
```