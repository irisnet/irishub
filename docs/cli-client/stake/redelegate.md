# iriscli stake redelegate

## Introduction

Transfer delegation from one validator to another one.

## Usage

```
iriscli stake redelegate [flags]
```

Print help messages:

```
iriscli stake redelegate --help
```

## Unique Flags

| Name, shorthand            | type   | Required | Default  | Description                                                         |
| -------------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-validator-dest   | string | true     | ""       | Bech address of the destination validator |
| --address-validator-source | string | true     | ""       | Bech address of the source validator |
| --shares-amount            | float  | false    | 0.0      | Amount of source-shares to either unbond or redelegate as a positive integer or decimal |
| --shares-percent           | float  | false    | 0.0      | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

Users must specify the redeleagtion token amount. There two options can do this: `--shares-amount` or `--shares-percent`. Keep in mind, don't specify both of them.

## Examples

If you want to redelegte 10% of your share:
```
iriscli stake redelegate --chain-id=test-chain-kvGwYI --from=fuxi --fee=0.004iris --address-validator-source=fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd --address-validator-dest=fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll  --shares-percent=0.1
```
After that, you will get the following ouptut:


```json
 {
   "code": 0,
   "data": "DAiX2MzgBRCAtK+tAQ==",
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 27011,
   "codespace": "",
   "tags": {
     "action": "begin_redelegate",
     "delegator": "faa106nhdckyf996q69v3qdxwe6y7408pvyvufy0x2",
     "destination-validator": "fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll",
     "end-time": "\u000c\u0008\ufffd\ufffd\ufffd\ufffd\u0005\u0010\ufffd\ufffd\ufffd\ufffd\u0001",
     "source-validator": "fva106nhdckyf996q69v3qdxwe6y7408pvyvfcwqmd"
   }
 })
```