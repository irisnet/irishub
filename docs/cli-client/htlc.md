# iriscli htlc

HTLC allows you to manage local Hash Time Locked Contracts (HTLCs) for atomic swaps with other chains.

There are the following states involved in the lifecycle of an HTLC:
   - open: indicates the HTLC is claimable
   - completed: indicates the HTLC has been claimed
   - expired: indicates the HTLC is expired and refundable
   - refunded: indicates the HTLC has been refunded

## Available Commands

| Name                                  | Description                 |
| ------------------------------------- | --------------------------- |
| [create](#iriscli-htlc-create)        | Create an HTLC              |
| [claim](#iriscli-htlc-claim)          | Claim an opened HTLC        |
| [refund](#iriscli-htlc-refund)        | Refund from an expired HTLC |
| [query-htlc](#iriscli-htlc-query-htlc) | Query details of an HTLC    |

## iriscli htlc create

Create an HTLC

```bash
iriscli htlc create --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --to=<to> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --secret=<secret> --time-lock=<time-lock> --timestamp=<timestamp>
```

**Flags:**

| Name, shorthand           | Type     | Required | Default | Description                                                       |
| ------------------------- | -------- | -------- | ------- | ----------------------------------------------------------------- |
| --to                      | string   | Yes      |         | Bech32 encoding address to receive coins                          |
| --receiver-on-other-chain | string   |          |         | The claim receiving address on the other chain                 |
| --amount                  | string   | Yes      |         | Similar to the amount in the original transfer                    |
| --secret                  | bytesHex |          |         | The secret for generating the hash lock, omission will be randomly generated |
| --hash-lock               | bytesHex | Yes      |         | The sha256 hash generated from secret (and timestamp if provided) |
| --time-lock               | string   | Yes      |         | The number of blocks to wait before the asset may be returned to  |
| --timestamp               | uint     |          |         | The timestamp in seconds for generating hash lock if provided     |

### Create an HTLC

```bash
iriscli htlc create \
--from=node0 \
--to=faa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5 \
--receiver-on-other-chain=0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826 \
--amount=10iris \
--secret=382aa2863398a31474616f1498d7a9feba132c4bcf9903940b8a5c72a46e4a41 \
--time-lock=50 \
--timestamp=1580000000 \
--fee=0.3iris \
--chain-id=test \
--commit
```

## iriscli htlc claim

Claim an opened HTLC

```bash
iriscli htlc claim --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --hash-lock=<hash-lock> --secret=<secret>
```

**Flags:**

| Name, shorthand | Type     | Required | Default | Description                                      |
| --------------- | -------- | -------- | ------- | ------------------------------------------------ |
| --hash-lock     | bytesHex | Yes      |         | The hash lock identifying the HTLC to be claimed |
| --secret        | bytesHex | Yes      |         | The secret for generating hash lock              |

### Claim an opened HTLC

```bash
iriscli htlc claim \
--from=node0 \
--hash-lock=bae5acb11ad90a20cb07023f4bf0fcf4d38549feff486dd40a1fbe871b4aabdf \
--secret=382aa2863398a31474616f1498d7a9feba132c4bcf9903940b8a5c72a46e4a41 \
--fee=0.3iris \
--chain-id=test \
--commit
```

## iriscli htlc refund

Refund from an expired HTLC

```bash
iriscli htlc refund --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --hash-lock=<hash-lock>
```

**Flags:**

| Name, shorthand | Type     | Required | Default | Description                                       |
| --------------- | -------- | -------- | ------- | ------------------------------------------------- |
| --hash-lock     | bytesHex | Yes     |         | The hash lock identifying the HTLC to be refunded |

### Refund from an expired HTLC

```bash
iriscli htlc refund \
--from=node0 \
--hash-lock=bae5acb11ad90a20cb07023f4bf0fcf4d38549feff486dd40a1fbe871b4aabdf \
--fee=0.3iris \
--chain-id=test \
--commit
```

## iriscli htlc query-htlc

Query details of an HTLC

```bash
iriscli htlc query-htlc <hash-lock>
```

### Query details of an HTLC

```bash
iriscli htlc query-htlc bae5acb11ad90a20cb07023f4bf0fcf4d38549feff486dd40a1fbe871b4aabdf
```

After that, you will get the detailed info for the HTLC.

```bash
HTLC:
        Sender:               faa1a2g4k9w3v2d2l4c4q5rvvu7ggjcrfnynvrpqze
        To:                   faa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5
        ReceiverOnOtherChain: 0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826
        Amount:               10iris
        Secret:               382aa2863398a31474616f1498d7a9feba132c4bcf9903940b8a5c72a46e4a41
        Timestamp:            1580000000
        ExpireHeight:         59
        State:                completed
```
