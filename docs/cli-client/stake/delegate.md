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
## Flags

| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --fee            | String | True     |                       | Fee to pay along with transaction                            |
| --from           | String | True     |                       | Name of private key with which to sign                       |
| --gas            | String | False    | 20000                 | Gas limit to set per-transaction; set to "simulate" to calculate required gas automatically |
| --gas-adjustment | Float  |          | 1                     | Adjustment factor to be multiplied against the estimate returned by the tx simulation; if the gas limit is set manually this flag is ignored |
| --generate-only  |        |          |                       | Build an unsigned transaction and write it to STDOUT         |
| --commit         | String | False     | True                  |wait for transaction commit accomplishment, if true, --async will be ignored|


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
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```
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