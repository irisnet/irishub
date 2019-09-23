# iriscli asset

Asset module allows you to manage assets on IRIS Hub

## Available Commands

| Name                                                | Description                                     |
| --------------------------------------------------- | ----------------------------------------------- |
| [create-gateway](create-gateway.md)                 | Create a gateway                                |
| [edit-gateway](edit-gateway.md)                     | Edit a gateway                                  |
| [transfer-gateway-owner](transfer-gateway-owner.md) | Transfer the ownership of a gateway             |
| [issue-token](issue-token.md)                       | Issue a new token                               |
| [edit-token](edit-token.md)                         | Edit an existing token                          |
| [transfer-token-owner](transfer-token-owner.md)     | Transfer the ownership of a token               |
| [mint-token](mint-token.md)                         | Mint tokens to a specified address              |
| [query-token](query-token.md)                       | Query details of a token                        |
| [query-tokens](query-tokens.md)                     | Query details of a group of tokens              |
| [query-gateway](query-gateway.md)                   | Query details of a gateway by the given moniker |
| [query-gateways](query-gateways.md)                 | Query all gateways with an optional owner       |
| [query-fee](query-fee.md)                           | Query the asset related fees                    |

## iriscli asset create-gateway

Create a gateway which is used to peg external assets

```bash
iriscli asset create-gateway [flags]
```

**Unique Flags:**

| Name, shorthand     | type   | Required | Default  | Description                                                                                              |
| --------------------| -----  | -------- | -------- | -------------------------------------------------------------------------------------------------------- |
| --moniker           | string | true     |          | The unique name with a size between 3 and 8, beginning with a letter followed by alphanumeric characters |
| --identity          | string |          |          | Optional identity signature with a maximum length of 128 (ex. UPort or Keybase)                          |
| --details           | string |          |          | Optional details with a maximum length of 280                                                            |
| --website           | string |          |          | Optional website with a maximum length of 128|

**Examples:**

```bash
iriscli asset create-gateway --moniker=cats --identity=<pgp-id> --details="Cat Tokens" --website="www.example.com" --from=<key-name> --chain-id=irishub --fee=0.3iris
```

## iriscli asset edit-gateway

Edit a gateway with the given moniker

```bash
iriscli asset edit-gateway [flags]
```

**Unique Flags:**

| Name, shorthand     | type   | Required | Default  | Description                                                                                              |
| --------------------| -----  | -------- | -------- | -------------------------------------------------------------------------------------------------------- |
| --moniker           | string | true     |          | The unique name with a size between 3 and 8, beginning with a letter followed by alphanumeric characters |
| --identity          | string |          |          | Optional identity signature with a maximum length of 128                                                 |
| --details           | string |          |          | Optional details with a maximum length of 280                                                            |
| --website           | string |          |          | Optional website with a maximum length of 128                                                            |

**Examples:**

```bash
iriscli asset edit-gateway --moniker=cats --identity=<pgp-id> --details="Cat Tokens" --website="http://www.example.com" --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```

## iriscli asset transfer-gateway-owner

Transfer the ownership of a gateway to a new owner.

```bash
iriscli asset transfer-gateway-owner [flags]
```

**Unique Flags:**

| Name, shorthand     | type    | Required | Default  | Description                                            |
| --------------------| ------- | -------- | -------- |------------------------------------------------------- |
| --moniker           | string  | true     |          | The unique name of the gateway to be transferred       |
| --to                | Address | true     |          | The new owner to which the gateway will be transferred |

**Examples:**

```bash
iriscli asset transfer-gateway-owner --moniker=cats --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

## iriscli asset issue-token

This command is used to issue a new token on IRIS Hub.

```bash
iriscli asset issue-token [flags]
```

**Unique Flags:**

| Name, shorthand    | Type    | Required | Default       | Description                                                  |
| ------------------ | ------- | -------- | ------------- | ------------------------------------------------------------ |
| --family           | string  | true     | fungible      | The token type: fungible, non-fungible (unsupported) |
| --source           | string  |          | native        | The token source: native, gateway                              |
| --name             | string  | true     |               | Name of the newly issued token, limited to 32 unicode characters, e.g. "IRIS Network" |
| --gateway          | string  |          |               | The unique moniker of the gateway, required when the source is gateway |
| --symbol           | string  | true     |               | The symbol of the token, length between 3 and 8, alphanumeric characters, case insensitive |
| --canonical-symbol | string  |          |               | When the source is gateway, it is used to identify the symbol on its' original chain |
| --min-unit-alias   | string  |          |               | The alias of minimum uint                                      |
| --initial-supply   | uint64  | true     |               | The initial supply of this token. The amount before boosting should not exceed 100 billion. |
| --max-supply       | uint64  |          | 1000000000000 | The hard cap of this token, total supply can not exceed max supply. The amount before boosting should not exceed 1000 billion.|
| --decimal          | uint8   | true     |               | A token can have a maximum of 18 digits of decimal         |
| --mintable         | boolean |          | false         | Whether this token could be minted(increased) after the initial issuing |

**Examples:**

### Issue native token

```bash
iriscli asset issue-token --family=fungible --source=native --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true --fee=1iris --from=<key-name> --commit
```

### Issue gateway token

#### Create a gateway

A gateway named `cats` is required to be created before this example, [more details](#iriscli-asset-create-gateway)

```bash
iriscli asset create-gateway --moniker=cats --identity=<identity> --details=<details> --website=<website> --from=<key-name> --commit
```

#### Issue a gateway token

```bash
iriscli asset issue-token --family=fungible --source=gateway --gateway=cats --canonical-symbol=cat --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true  --fee=1iris --from=<key-name> --commit
```

### Send tokens

You can send any tokens you have just like [sending iris](../bank/send)

#### Send native tokens

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10kitty --fee=0.3iris --chain-id=irishub
```

#### Send gateway tokens

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10cats.kitty --fee=0.3iris --chain-id=irishub
```

## iriscli asset edit-token

Edit token informations

```bash
iriscli asset edit-token <token-id> [flags]
```

**Unique Flags:**

| Name                   | Type   | Required | Default | Description                                        |
| ---------------------- | -----  | -------- | ------- | -------------------------------------------------- |
| --name                 | string |          |         | The token name, e.g. IRIS Network                  |
| --canonical-symbol     | string |          |         | The source symbol of a gateway or external token   |
| --min-unit-alias       | string |          |         | The token symbol minimum alias                     |
| --max-supply           | uint   |          | 0       | The max supply of the token                        |
| --mintable             | bool   |          | false   | Whether the token can be minted, default false     |

**Examples:**

`max-supply` can only be reduced and no less than the current total supply

```bash
iriscli asset edit-token cat --name="Cat Token" --canonical-symbol="cat" --min-unit-alias=kitty --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```

## iriscli asset transfer-token-owner

Transfer the ownership of a token

```bash
iriscli asset transfer-token-owner <token-id> [flags]
```

**Unique Flags:**

| Name | Type   | Required | Default | Description           |
| ---- | ------ | -------- | ------- | --------------------- |
| --to | string | true     |         | The new owner address |

**Examples:**

```bash
iriscli asset transfer-token-owner kitty --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```

## iriscli asset mint-token

The asset owner can directly mint tokens to a specified address

```bash
iriscli asset mint-token <token-id> [flags]
```

**Unique Flags:**

| Name     | Type   | Required | Default | Description                                           |
| -------- | ------ | -------- | ------- | ----------------------------------------------------- |
| --to     | string |          |         | Address of mint token to, default is your own address |
| --amount | uint64 | true     | 0       | Amount of the token to mint                           |

**Examples:**

```bash
iriscli asset mint-token kitty --amount=1000000 --from=<key-name> --chain-id=irishub --fee=0.4iris
```

## iriscli asset query-token

Query a token issued on IRIS Hub.

```bash
iriscli asset query-token <token-id>
```

### Global Unique Token ID Generation Rule

- When Source is native: ID = [Symbol], e.g. iris

- When Source is external: ID = x.[Symbol], e.g. x.btc

- When Source is gateway: ID = [Gateway].[Symbol], e.g. cats.kitty

**Examples:**

### Query the native token named "kitty"

```bash
iriscli asset query-token kitty
```

### Query the token of gateway "cats" named "kitty"

```bash
iriscli asset query-token cats.kitty
```

### Query the external token named "btc"

```bash
iriscli asset query-token x.btc
```

## iriscli asset query-tokens

Query the collection of tokens issued on IRIS Hub based on criteria.

```bash
iriscli asset query-tokens [flags]
```

**Unique Flags:**

| Name      | Type   | Required | Default | Description                                                        |
| --------- | ------ | -------- | ------- | ------------------------------------------------------------------ |
| --source  | string |          | all     | Token Source: native / gateway / external                          |
| --gateway | string |          |         | The unique moniker of the gateway, required when source is gateway |
| --owner   | string |          |         | The owner of the tokens                                            |

### Query Rules

- when source is native
  - gateway will be ignored
  - owner optional
- When source is gateway
  - gateway required
  - owner will be ignored (because gateway tokens are all owned by the gateway)
- when source is external
  - gateway and owner are ignored
- when the gateway is not empty
  - source optional

**Examples:**

### Query all tokens

```bash
iriscli asset query-tokens
```

### Query all native tokens

```bash
iriscli asset query-tokens --source=native
```

### Query all tokens of the gateway named "cats"

```bash
iriscli asset query-tokens --gateway=cats
```

### Query all tokens of the specified owner

```bash
iriscli asset query-tokens --owner=<address>
```

## iriscli asset query-gateway

Query a gateway by moniker

```bash
iriscli asset query-gateway [flags]
```

**Unique Flags:**

| Name, shorthand     | type   | Required | Default  | Description                                                         |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --moniker           | string | true     |          | The unique name with a size between 3 and 8, beginning with a letter followed by alphanumeric characters |

**Examples:**

```bash
iriscli asset query-gateway --moniker cats
```

## iriscli asset query-gateways

Query all the gateways by its' owner

```bash
iriscli asset query-gateways [flags]
```

**Unique Flags:**

| Name, shorthand  | type    | Required | Default  | Description                        |
| ---------------- | ------- | -------- | -------- | ---------------------------------- |
| --owner          | Address |          |          | The owner address to be queried by |

**Examples:**

```bash
iriscli asset query-gateways --owner=<owner-address>
```

## iriscli asset query-fee

Query the asset related fees, including gateway creation and token issuance and minting

```bash
iriscli asset query-fee [flags]
```

**Unique Flags:**

| Name, shorthand     | type   | Required | Default  | Description                                            |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------ |
| --gateway           | string |          |          | The gateway moniker, required for querying gateway fee |
| --token             | string |          |          | The token id, required for querying token fees         |

**Examples:**

### Query fee of creating a gateway

```bash
iriscli asset query-fee --gateway=cats
```

### Query fee of issuing and minting a native token

```bash
iriscli asset query-fee --token=kitty
```

### Query fee of issuing and minting a gateway token

```bash
iriscli asset query-fee --token=cats.kitty
```
