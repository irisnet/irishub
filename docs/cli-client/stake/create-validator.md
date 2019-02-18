# iriscli stake create-validator

## Introduction

Send transaction to apply to be validator and delegate a certain amount tokens on it

## Usage

```
iriscli stake create-validator [flags]
```

Print help messages:
```
iriscli stake create-validator --help
```

## Unique Flags

| Name, shorthand              | type   | Required | Default  | Description                                                         |
| ---------------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --amount                     | string | true     | ""       | Amount of coins to bond |
| --commission-rate            | float  | true     | 0.0      | The initial commission rate percentage |
| --details                    | string | false    | ""       | Optional details |
| --genesis-format             | bool   | false    | false    | Export the transaction in gen-tx format; it implies --generate-only |
| --identity                   | string | false    | ""       | Optional identity signature (ex. UPort or Keybase) |
| --ip                         | string | false    | ""       | Node's public IP. It takes effect only when used in combination with |
| --moniker                    | string | true     | ""       | Validator name |
| --pubkey                     | string | true     | ""       | Go-Amino encoded hex PubKey of the validator. For Ed25519 the go-amino prepend hex is 1624de6220 |
| --website                    | string | false    | ""       | Optional website |

## Examples

```
iriscli stake create-validator --chain-id=test-irishub --from=<key name> --fee=0.4iris --pubkey=<Validator PubKey> --commission-rate=0.1 --amount=100iris --moniker=<validator name>
```
Sample Output:
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
   "gas_used": 12050,
   "codespace": "",
   "tags": {
     "action": "create_validator",
     "destination-validator": "fva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll",
     "identity": "",
     "moniker": "test"
   }
 })
```