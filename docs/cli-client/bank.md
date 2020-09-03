# Bank

Bank module allows you to manage assets in your local accounts

## Available Commands

| Name                                             | Description                         |
| ------------------------------------------------ | ----------------------------------- |
| [balances](#iris-query-bank-balances)            | Query for account balances by address                     |
| [total](#iris-query-bank-total)                  | Query the total supply of coins of the chain                   |
| [send](#iris-tx-bank-send)                       | Create and/or sign and broadcast a MsgSend transaction     |

## iris query bank balances

Query the total balance of an account or of a specific denomination.

```bash
iris query bank balances [address] [flags]
```

**Flags:**

| Name, shorthand | Type   | Required | Default               | Description                                                   |
| --------------- | ------ | -------- | --------------------- | ------------------------------------------------------------- |
| -h, --help      |        |          |                       | Help for coin-type                                            |
| --denom         | string |          |                       | The specific balance denomination to query for                |

### iris query bank total

Query total supply of coins that are held by accounts in the chain.

```bash
iris query bank total [flags]
```
**Flags:**

| Name, shorthand | Type   | Required | Default               | Description                                                   |
| --------------- | ------ | -------- | --------------------- | ------------------------------------------------------------- |
| -h, --help      |        |          |                       | Help for coin-type                                            |
| --denom         | string |          |                       | The specific balance denomination to query for                |


## iris tx bank send

Sending tokens to another address, this command includes `generate`, `sign` and `broadcast` steps.

```bash
iris tx bank send [from_key_or_address] [to_address] [amount] [flags]
```

**Flags:**

| Name, shorthand | Type   | Required | Default | Description                                   |
| --------------- | ------ | -------- | ------- | --------------------------------------------- |
| --amount        | string | true     |         | Amount of coins to send, for instance: 10iris |
| --to            | string |          |         | Bech32 encoding address to receive coins      |

### Send tokens to another address

```bash
iris tx bank send [from_key_or_address] [to_address] [amount] [flags]
```
