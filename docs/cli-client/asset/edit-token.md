# iriscli asset edit-token

## Description

Edit token informations

## Usage

```bash
iriscli asset edit-token [flags]
```

## Flags

| Name                   | Type   | Required | Default | Description                                        |
| ---------------------- | -----  | -------- | ------- | -------------------------------------------------- |
| --name                 | string |          |         | The token name, e.g. IRIS Network                  |
| --canonical-symbol     | string |          |         | The source symbol of a gateway or external token   |
| --min-unit-alias       | string |          |         | The token symbol minimum alias                     |
| --max-supply           | uint   |          | 0       | The max supply of the token                        |
| --mintable             | bool   |          | false   | Whether the token can be minted, default false     |

## Example

`max-supply` can only be reduced and no less than the current total supply

```bash
iriscli asset edit-token cat --name="Cat Token" --canonical-symbol="cat" --min-unit-alias=kitty --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```
