# iriscli bank

Bank module allows you to manage assets in your local accounts

## Available Commands

| Name                                             | Description                         |
| ------------------------------------------------ | ----------------------------------- |
| [coin-type](#iriscli-bank-coin-type)             | Query coin type                     |
| [token-stats](#iriscli-bank-token-stats)         | Query token stats                   |
| [account](#iriscli-bank-account)                 | Query account balance               |
| [send](#iriscli-bank-send)                       | Create/sign/broadcast a send tx     |
| [burn](#iriscli-bank-burn)                       | Burn tokens                         |
| [set-memo-regexp](#iriscli-bank-set-memo-regexp) | Set memo regexp                     |

## Common Problems

### ERROR: decoding bech32 failed

```bash
iriscli bank account iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429zz
ERROR: decoding bech32 failed: checksum failed. Expected 46vaym, got d429zz.
```

This means the account address is misspelled, please double check the address.

### ERROR: account xxx does not exist

```bash
iriscli bank account iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr
ERROR: {"codespace":"sdk","code":9,"message":"account iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr does not exist"}
```

This is usually because the account address you are querying has no transactions on the chain.

## iriscli bank coin-type

Query a special kind of token on IRIShub. The native token on IRIShub is `iris`, which has the following available units: `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` and `iris-atto`.

```bash
 iriscli bank coin-type <coin_name> <flags>
```

**Flags:**

| Name, shorthand | Type   | Required | Default               | Description                                                   |
| --------------- | ------ | -------- | --------------------- | ------------------------------------------------------------- |
| -h, --help      |        |          |                       | Help for coin-type                                            |
| --chain-id      | string |          |                       | Chain ID of tendermint node                                   |
| --height        | int    |          |                       | Block height to query, omit to get most recent provable block |
| --indent        | string |          |                       | Add indent to JSON response                                   |
| --ledger        | string |          |                       | Use a connected Ledger device                                 |
| --node          | string |          | tcp://localhost:26657 | `<host>:<port>` to tendermint rpc interface for this chain    |
| --trust-node    | string |          | true                  | Don't verify proofs for responses                             |

### Query native token `iris`

```bash
iriscli bank coin-type iris
```

After that, you will get the detail info for the native token `iris`

```bash
CoinType:
  Name:     iris
  MinUnit:  iris-atto: 18
  Units:    iris: 0,  iris-milli: 3,  iris-micro: 6,  iris-nano: 9,  iris-pico: 12,  iris-femto: 15,  iris-atto: 18
  Origin:   native
  Desc:     IRIS Network
```

## iriscli bank token-stats

Query the token statistic, including total loose tokens, total burned tokens and total bonded tokens.

```bash
 iriscli bank token-stats <token-id> <flags>
```

**Flags:**

| Name, shorthand | Type   | Required | Default               | Description                                                  |
| --------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help      |        |          |                       | Help for token-stats                                           |
| --chain-id      | string |          |                       | Chain ID of tendermint node                                 |
| --height        | int    |          |                       | Block height to query, omit to get most recent provable block|
| --indent        | string |          |                       | Add indent to JSON response                                  |
| --ledger        | string |          |                       | Use a connected Ledger device                                |
| --node          | string |          | tcp://localhost:26657 | `<host>:<port>` to tendermint rpc interface for this chain    |
| --trust-node    | string |          | true                  | Don't verify proofs for responses                            |

### Query the token statistic

```bash
iriscli bank token-stats iris
```

Output:

```bash
TokenStats:
  Loose Tokens:             1404158512.076790096410637686iris
  Bonded Tokens:            609544925.59191727606475175iris
  Burned Tokens:            19205096.20000004iris
  Total Supply:             2013703437.668707372475389436iris
```

## iriscli bank account

This command is used for querying balance information of certain address.

```bash
iriscli bank account <address> <flags>
```

**Flags:**

| Name, shorthand | Type   | Required | Default               | Description                                                   |
| --------------- | ------ | -------- | --------------------- | ------------------------------------------------------------- |
| -h, --help      |        |          |                       | Help for account                                              |
| --chain-id      | string |          |                       | Chain ID of tendermint node                                   |
| --height        | int    |          |                       | Block height to query, omit to get most recent provable block |
| --ledger        | string |          |                       | Use a connected Ledger device                                 |
| --node          | string |          | tcp://localhost:26657 | `<host>:<port>` to tendermint rpc interface for this chain    |
| --trust-node    | string |          | true                  | Don't verify proofs for responses                             |

### Query your account in trust-mode

```bash
iriscli bank account <address>
```

Output:

```bash
Account:
  Address:         iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym
  Pubkey:          iap1addwnpepqwnsrt9m8tevhy4fdqyarunzuzzgz8e5q8jlceyf7uwpw0q0ptp2cp3lmjt
  Coins:           50iris
  Account Number:  0
  Sequence:        2
  Memo Regexp:
```

## iriscli bank send

Sending tokens to another address, this command includes `generate`, `sign` and `broadcast` steps.

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=<amount> --fee=<native-fee> --chain-id=<chain-id>
```

**Flags:**

| Name, shorthand | Type   | Required | Default | Description                                   |
| --------------- | ------ | -------- | ------- | --------------------------------------------- |
| --amount        | string | true     |         | Amount of coins to send, for instance: 10iris |
| --to            | string |          |         | Bech32 encoding address to receive coins      |

### Send tokens to another address

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10iris --fee=0.3iris --chain-id=irishub
```

## iriscli bank burn

This command is used to burn tokens from your own address.

```bash
iriscli bank burn --from=<key-name> --amount=<amount-to-burn> --fee=<native-fee> --chain-id=<chain-id>
```

**Flags:**

| Name, shorthand  | Type   | Required | Default | Description                          |
| ---------------- | ------ | -------- | ------- | ------------------------------------ |
| --amount         | string | true     |         | Amount of coins to burn, e.g. 10iris |

### Burn Token

```bash
 iriscli bank burn --from=<key-name> --amount=10iris --chain-id=irishub --fee=0.3iris
```

## iriscli bank set-memo-regexp

This command is used to set memo regexp for your own address, so that you can only receive coins from transactions with the corresponding memo.

```bash
iriscli bank set-memo-regexp --regexp=<regular-expression> --from=<key-name> --fee=<native-fee> --chain-id=<chain-id>
```

**Flags:**

| Name, shorthand | Type   | Required | Default | Description                                                |
| --------------- | ------ | -------- | ------- | ---------------------------------------------------------- |
| --regexp        | string | true     |         | Regular expression, maximum length 50, e.g. ^[A-Za-z0-9]+$ |

### Set memo regexp for an account address

```bash
iriscli bank set-memo-regexp --regexp=^[A-Za-z0-9]+$ --from=<key-name> --fee=0.3iris --chain-id=irishub
```
