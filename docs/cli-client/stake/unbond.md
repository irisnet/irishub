# iriscli stake unbond

## Description

Unbond shares from a validator

## Usage

```
iriscli stake unbond [flags]
```

Print all help messages:

```shell
iriscli stake unbond --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --address-delegator | string | true     | ""       | Bech address of the delegator |
| --shares-amount     | float  | false    | 0.0      | Amount of source-shares to either unbond or redelegate as a positive integer or decimal |
| --shares-percent    | float  | false    | 0.0      | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

Users must specify the unbond amount. There two options can do this: `--shares-amount` or `--shares-percent`. Keep in mind, don't specify them both.

## Examples

### Unbond shares from a validator

```shell
iriscli stake unbond --address-validator=ValidatorAddress --shares-percent=SharePercent --from=UnbondInitiator --chain-id=ChainID --fee=Fee
```

After that, you're done with unbonding shares from specified validator.

```txt
Committed at block 851 (tx hash: A82833DE51A4127BD5D60E7F9E4CD5895F97B1B54241BCE272B68698518D9D2B, response: {Code:0 Data:[11 8 230 225 179 223 5 16 249 233 245 21] Log:Msg 0:  Info: GasWanted:200000 GasUsed:16547 Tags:[{Key:[97 99 116 105 111 110] Value:[98 101 103 105 110 45 117 110 98 111 110 100 105 110 103] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 108 101 103 97 116 111 114] Value:[102 97 97 49 51 108 99 119 110 120 112 121 110 50 101 97 51 115 107 122 109 101 107 54 52 118 118 110 112 57 55 106 115 107 56 113 109 104 108 54 118 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[115 111 117 114 99 101 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 54 100 54 116 104 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[101 110 100 45 116 105 109 101] Value:[11 8 230 225 179 223 5 16 249 233 245 21] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 56 50 55 51 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "begin-unbonding",
     "completeConsumedTxFee-iris-atto": "\"8273500000000000\"",
     "delegator": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
     "end-time": "\u000b\u0008\ufffd\ufffd\ufffd\ufffd\u0005\u0010\ufffd\ufffd\ufffd\u0015",
     "source-validator": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd"
   }
 }

```
