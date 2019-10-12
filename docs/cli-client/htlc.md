# iriscli htlc

HTLC allows you to manage local Hash Time Locked Contracts (HTLCs) for atomic swaps with other chains

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
iriscli htlc create --receiver=<receiver> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --hash-lock=<hash-lock> --time-lock=<time-lock> --timestamp=<timestamp>
```

**Flags:**

| Name, shorthand           | Type     | Required | Default | Description                                                       |
| ------------------------- | -------- | -------- | ------- | ----------------------------------------------------------------- |
| --receiver                | string   | Yes      |         | Bech32 encoding address to receive coins                          |
| --receiver-on-other-chain | bytesHex | Yes      |         | The receiver address on the other chain                           |
| --amount                  | string   | Yes      |         | Similar to the amount in the original transfer                    |
| --hash-lock               | bytesHex | Yes      |         | The sha256 hash generated from secret (and timestamp if provided) |
| --time-lock               | string   | Yes      |         | The number of blocks to wait before the asset may be returned to  |
| --timestamp               | uint     |          |         | The timestamp in seconds for generating hash lock if provided     |

### Create an HTLC

```bash
iriscli htlc create \
--from=userX \
--receiver=faa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5 \
--receiver-on-other-chain=bb9188215a6112a6f7eb93e3e929197b3d44004cb691f95babde84cc18789364 \
--hash-lock=e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561 \
--amount=10iris \
--time-lock=50 \
--timestamp=1580000000 \
--fee=0.3iris \
--chain-id=testNet \
--commit
```

## iriscli htlc claim

Claim an opened HTLC

```bash
iriscli htlc claim --hash-lock=<hash-lock> --secret=<secret>
```

**Flags:**

| Name, shorthand | Type     | Required | Default | Description                                      |
| --------------- | -------- | -------- | ------- | ------------------------------------------------ |
| --hash-lock     | bytesHex | Yes      |         | The hash lock identifying the HTLC to be claimed |
| --secret        | bytesHex | Yes      |         | The secret for generating hash lock              |

### Claim an opened HTLC

```bash
iriscli htlc claim \
--from=userX \
--hash-lock=f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20 \
--secret=5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f \
--fee=0.3iris \
--chain-id=testNet \
--commit
```

## iriscli htlc refund

Refund from an expired HTLC

```bash
iriscli htlc refund --hash-lock=<hash-lock>
```

**Flags:**

| Name, shorthand | Type     | Required | Default | Description                                       |
| --------------- | -------- | -------- | ------- | ------------------------------------------------- |
| --hash-lock     | bytesHex | Yes     |         | The hash lock identifying the HTLC to be refunded |

### Refund from an expired HTLC

```bash
iriscli htlc refund \
--from=userX \
--hash-lock=f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20 \
--fee=0.3iris \
--chain-id=testNet \
--commit
```

## iriscli htlc query-htlc

Query details of an HTLC

```bash
iriscli htlc query-htlc <hash-lock>
```

### Query details of an HTLC

```bash
iriscli htlc query-htlc f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20
```

After that, you will get the detail info for the account.

```bash
HTLC:
    Sender:               iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym
    Receiver:             iaa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5
    ReceiverOnOtherChain: 72656365697665724f6e4f74686572436861696e
    Amount:               10000000000000000000iris-atto
    Secret:               5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f
    Timestamp:            0
    ExpireHeight:         107
    State:                completed
```
