# iriscli asset create-gateway

## Introduction

Create a gateway which is used to map the external assets

## Usage

```
iriscli asset create-gateway [flags]
```

Print help messages:
```
iriscli asset create-gateway --help
```

## Unique Flags

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | true     | ""       | the unique name with a size between 3 and 8 letters|
| --identity          | string | false    | ""       | Optional identity signature (ex. UPort or Keybase) |
| --details           | string | false    | ""       | Optional website with a maximum length of 280|
| --website           | string | false    | ""       | Optional website with a maximum length of 128|


## Examples

```
iriscli asset create-gateway --moniker=<name> --identity=<identity> --details=<description> 
--website=<website> --chain-id=<chain-id> --from=<key_name> --fee=0.3iris
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
     "moniker": "testgw"
   }
 })
```