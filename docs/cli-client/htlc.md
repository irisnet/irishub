# HTLC

[HTLC module](../features/htlc.md) allows you to manage local Hash Time Locked Contracts (HTLCs) for atomic swaps with other chains.

There are the following states involved in the lifecycle of an HTLC:

- open: indicates the HTLC is claimable
- completed: indicates the HTLC has been claimed
- refunded: indicates the HTLC has been refunded

## Available Commands

| Name                           | Description              |
| ------------------------------ | ------------------------ |
| [create](#iris-tx-htlc-create) | Create an HTLC           |
| [claim](#iris-tx-htlc-claim)   | Claim an opened HTLC     |
| [htlc](#iris-query-htlc-htlc)  | Query details of an HTLC |

## iris tx htlc create

Create an HTLC

```bash
iris tx htlc create \
    --to=<recipient> \
    --receiver-on-other-chain=<receiver-on-other-chain> \
    --sender-on-other-chain=<sender-on-other-chain> \
    --amount=<amount> \
    --hash-lock=<hash-lock> \
    --secret=<secret> \
    --timestamp=<timestamp> \
    --time-lock=<time-lock> \
    --transfer=<true | false> \
    --from=mykey
```

**Flags:**

| Name, shorthand           | Type   | Required | Default | Description                                                                                           |
| ------------------------- | ------ | -------- | ------- | ----------------------------------------------------------------------------------------------------- |
| --to                      | string | Yes      |         | Bech32 encoding address to receive coins                                                              |
| --receiver-on-other-chain | string |          |         | The claim receiving address on the other chain                                                        |
| --sender-on-other-chain   | string |          |         | The counterparty creator address on the other chain                                                   |
| --amount                  | string | Yes      |         | Similar to the amount in the original transfer                                                        |
| --secret                  | string |          |         | The secret for generating the hash lock, generated randomly if omitted                                |
| --hash-lock               | string |          |         | The sha256 hash generated from secret (and timestamp if provided), generated from `secret` if omitted |
| --time-lock               | string | Yes      |         | The number of blocks to wait before the asset may be returned to                                      |
| --timestamp               | uint   |          |         | The timestamp in seconds for generating hash lock if provided                                         |
| --transfer                | bool   |          | false   | Whether it is an HTLT transaction                                                                     |

## iris tx htlc claim

Claim an opened HTLC

```bash
iris tx htlc claim [id] [secret] [flags] --from=mykey
```

## iris query htlc htlc

Query details of an HTLC

```bash
iris query htlc htlc <id>
```

## iris query htlc params

Query params of HTLC module

```bash
iris query htlc params
```

## iris query htlc supplies

Query supplies of all HTLT assets

```bash
iris query htlc supplies
```

## iris query htlc supply

Query supply of an HTLT asset

```bash
iris query htlc supply [denom]
```
