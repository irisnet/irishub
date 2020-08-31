# iris bank

Bank module allows you to manage assets in your local accounts

## Available Commands

| Name                                             | Description                         |
| ------------------------------------------------ | ----------------------------------- |
| [balances](#iris-query-bank-balances)            | Query for account balances by address                     |
| [total](#iris-query-bank-total)                  | Query the total supply of coins of the chain                   |
| [send](#iris-tx-bank-send)                       | Create and/or sign and broadcast a MsgSend transaction     |

## Common Problems

### ERROR: decoding bech32 failed

```bash
iris bank account iaa1a0x4g8rqc90l3z9jh98x7mkd0w77e9q9r300h 
Error: decoding bech32 failed: checksum failed. Expected 9r300k, got 9r300h.
```

This means the account address is misspelled, please double check the address.

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
iris tx bank send --from=<key-name> --to=<address> --amount=10iris --fees=0.3iris --chain-id=irishub
```
