# iriscli stake edit-validator

## Introduction

Edit existing validator, such as commission rate, name and other description message.

## Usage

```
iriscli stake edit-validator [flags]
```
Print help messages:
```
iriscli stake edit-validator --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --commission-rate   | string | float    | 0.0      | Commission rate percentage |
| --moniker           | string | false    | ""       | Validator name |
| --identity          | string | false    | ""       | Optional identity signature (ex. UPort or Keybase) |
| --website           | string | false    | ""       | Optional website  |
| --details           | string | false    | ""       | Optional details |


## Examples

```
iriscli stake edit-validator --from=<key name> --chain-id=<chain-id> --fee=0.004iris --commission-rate=0.15
```

```json
{
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3482,
   "codespace": "",
   "tags": {
     "action": "edit_validator",
     "destination-validator": "fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd",
     "identity": "",
     "moniker": "test2"
   }
 })
```