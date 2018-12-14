# iriscli stake unbond

## Introduction

Unbond tokens from a validator

## Usage

```
iriscli stake unbond [flags]
```

Print help messages:
```
iriscli stake unbond --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator | string | true     | ""       | Bech address of the validator |
| --shares-amount     | float  | false    | 0.0      | Amount of source-shares to either unbond or redelegate as a positive integer or decimal |
| --shares-percent    | float  | false    | 0.0      | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

Users must specify the unbond amount. There two options can do this: `--shares-amount` or `--shares-percent`. Keep in mind, don't specify both of them.

## Examples

```
iriscli stake unbond --address-validator=<ValidatorAddress> --shares-percent=0.1 --from=<key name> --chain-id=<chain-id> --fee=0.004iris
```
```json
 {
   "code": 0,
   "data": "CwiAkrjDmP7///8B",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 18990,
   "codespace": "",
   "tags": {
     "action": "begin_unbonding",
     "delegator": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2",
     "end-time": "\u000b\u0008\ufffd\ufffd\ufffdØ\ufffd\ufffd\ufffd\ufffd\u0001",
     "source-validator": "fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll"
   }
 })
```