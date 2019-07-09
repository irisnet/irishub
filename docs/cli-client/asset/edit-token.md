# iriscli asset edit-token

## Description

Edit token informations

## Usage

```bash
iriscli asset edit-token [flags]
```

## Flags

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --name           | string | No | "" | the token name, e.g. IRIS Network |
| --symbol-at-source | string | No | "" | the source symbol of a gateway or external token |
| --symbol-min-alias | string | No | "" | the token symbol minimum alias |
| --max-supply | uint | No | 0 | the max supply of the token |
| --mintable | bool | No | false | whether the token can be minted, default false |

## Example

`max-supply` can only be reduced and no less than the current total supply

```bash
iriscli asset edit-token cat --name="Cat Token" --symbol-at-source="cat" --symbol-min-alias=kitty --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```
