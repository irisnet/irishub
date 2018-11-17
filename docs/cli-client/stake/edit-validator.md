# iriscli stake edit-validator

## Description

Edit existing validator

## Usage

```
iriscli stake edit-validator [flags]
```
Print help messages:
```shell
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

### Edit existing validator account

```shell
iriscli stake edit-validator --from=<key name> --chain-id=<chain-id> --fee=0.004iris --commission-rate=0.15
```

After that, you're done with editing a new validator.

```txt
Committed at block 2160 (tx hash: C48CABDA1183B5319003433EB1FDEBE5A626E00BD319F1A84D84B6247E9224D1, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3540 Tags:[{Key:[97 99 116 105 111 110] Value:[101 100 105 116 45 118 97 108 105 100 97 116 111 114] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 115 116 105 110 97 116 105 111 110 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 54 100 54 116 104 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[109 111 110 105 107 101 114] Value:[117 98 117 110 116 117 49 56] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[105 100 101 110 116 105 116 121] Value:[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 55 55 48 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "edit-validator",
     "completeConsumedTxFee-iris-atto": "\"177000000000000\"",
     "destination-validator": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
     "identity": "",
     "moniker": "ubuntu18"
   }
}
```
