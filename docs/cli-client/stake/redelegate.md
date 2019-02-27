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
iriscli stake redelegate --chain-id=test-irishub --from=fuxi --fee=0.3iris --address-validator-source=iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms --address-validator-dest=iva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6s30mrlz  --shares-percent=0.1
```
After that, you will get the following ouptut:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

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
     "delegator": "iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh",
     "destination-validator": "iva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6s30mrlz",
     "end-time": "\u000c\u0008\ufffd\ufffd\ufffd\ufffd\u0005\u0010\ufffd\ufffd\ufffd\ufffd\u0001",
     "source-validator": "iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms"
   }
 })
```