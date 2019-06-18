# iriscli asset issue-token

## Description

This command is used to issue token in IRIS Hub.

## Usage

issue  10000000000 kitty
```
iriscli asset issue-token [flags]
```


## Flags

| Name,shorthand | Type   | Required | Default               | Description                                                                                                                                                                                                                                                                           |
| ----------------  | --------- | ------- | ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------  |
| --family          | byte      | false       |  fungible          |fungible or non-fungible                                                                                                                                                                                                                                                  |
| --name            | string    | true        |                    | Name of the newly issued asset, limited to 32 unicode(english character , number, _) characters,  e.g. "IRISnet"                                                                                                                                                                 |
| --gateway         | string    | false       |                    | The symbol of gateway; when the source is gateway, this field is required                                                                                                                                                                                                         |
| --symbol          | string    | true        |                    | The length of the string for representing this asset is between 3 and 6 alphanumeric characters and is case insensitive                                                                                                                                                          |
| --symbol-at-source| string    | false       |                    | When the source is external or gateway, it is used as the identify of symbol in its source chain                                                                                                                                                                                   |
| --symbol-min-alias| string    | false       |                    | The alias of minim uint                                                                                                                                                                                    |
| --source          | string    | false       |  native            | native, external or Gateway IDs                                                                                                                                                                                                                               |
| --initial-supply  | uint64    | true        |                    | The initial supply for this asset. The amount before boosting should not exceed 100 billion. The amount should be positive integer                                                                                                                                               |
| --max-supply      | uint64    | true        |  1000000000000     | The hard cap of this asset, total supply can not exceed max supply. The amount should be positive integer                                                                                                                                                                        |
| --decimal         | uint8     | false       |  0                 | The asset can have a maximum of 18 digits of decimal                                                                                                                                                                                                                             |
| --mintable        | boolean   | false       |  false             | Whether this asset could be minted(increased) after the initial issuing                                                                                                                                                                                                          |




## Examples

### Issue native token

```
iriscli asset issue-token --family=fungible --name=kittyToken --symbol=kitty --source=native --initial-supply=10000000000 --max-supply=1000000000000 --decimal=0 --mintable=<true/false>    --fee=1iris
```


### Issue gateway token

```
iriscli asset issue-token --family=fungible --symbol-at-source=cat --name=kittyToken --symbol=kitty --source=gateway --gateway=gtty --initial-supply=10000000000 --max-supply=1000000000000 --decimal=0 --mintable=<true/false>  --fee=1iris
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
    "gas_wanted": 50000,
    "gas_used": 7008,
    "codespace": "",
    "tags": [
      {
        "key": "action",
        "value": "issue_token"
      },
      {
        "key": "action",
        "value": "issue-token"
      },
      {
        "key": "token-id",
        "value": "kitty"
      },
      {
        "key": "token-denom",
        "value": "kitty-min"
      },
      {
        "key": "token-source",
        "value": "native"
      },
      {
        "key": "token-gateway",
        "value": ""
      },
      {
        "key": "token-owner",
        "value": "faa1j8mlkem7s9a0jjhkd39zd24xh9gdj6zus77v7q"
      }
    ]
  })
 

```