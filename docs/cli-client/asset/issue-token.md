# iriscli asset issue-token

## Description

This command is used to issue a new token on IRIS Hub.

## Usage

```bash
iriscli asset issue-token [flags]
```


## Flags

| Name,shorthand     | Type    | Required | Default       | Description                                                  |
| ------------------ | ------- | -------- | ------------- | ------------------------------------------------------------ |
| --family           | string  | true     | fungible      | The token type: fungible, non-fungible (unsupported) |
| --source           | string  | false    | native        | The token source: native, gateway                              |
| --name             | string  | true     |               | Name of the newly issued token, limited to 32 unicode characters, e.g. "IRIS Network" |
| --gateway          | string  | false    |               | The unique moniker of the gateway, required when the source is gateway |
| --symbol           | string  | true     |               | The symbol of the token, length between 3 and 6, alphanumeric characters, case insensitive |
| --symbol-at-source | string  | false    |               | When the source is gateway, it is used to identify the symbol on its' original chain |
| --symbol-min-alias | string  | false    |               | The alias of minimum uint                                      |
| --initial-supply   | uint64  | true     |               | The initial supply of this token. The amount before boosting should not exceed 100 billion. |
| --max-supply       | uint64  | false    | 1000000000000 | The hard cap of this token, total supply can not exceed max supply. The amount before boosting should not exceed 1000 billion.|
| --decimal          | uint8   | false    | 0             | A token can have a maximum of 18 digits of decimal         |
| --mintable         | boolean | false    | false         | Whether this token could be minted(increased) after the initial issuing |

## Examples

### Issue a native token

```bash
iriscli asset issue-token --family=fungible --source=native --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true --fee=1iris --from=<key-name> --commit
```

### Issue a gateway token

#### Create a gateway

A gateway named `cats` is required to be created before this example, [more details](./create-gateway.md)

```bash
iriscli asset create-gateway --moniker=cats --identity=<identity> --details=<details> --website=<website> --from=<key-name> --commit
```

#### Issue a gateway token

```bash
iriscli asset issue-token --family=fungible --source=gateway --gateway=cats --symbol-at-source=cat --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true  --fee=1iris --from=<key-name> --commit
```
