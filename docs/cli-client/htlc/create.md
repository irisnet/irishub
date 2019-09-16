# iriscli htlc create

## Description

Create an HTLC

## Usage

```bash
iriscli htlc create --receiver=<receiver> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --hash-lock=<hash-lock> --time-lock=<time-lock> --timestamp=<timestamp>
```

## Flags

| Name, shorthand           | Type     | Required | Default | Description                                                       |
| ------------------------- | -------- | -------- | ------- | ----------------------------------------------------------------- |
| --receiver                | string   | true     |         | Bech32 encoding address to receive coins                          |
| --receiver-on-other-chain | bytesHex | true     |         | The receiver address on the other chain                           |
| --amount                  | string   | true     |         | Similar to the amount in the original transfer                    |
| --hash-lock               | bytesHex | true     |         | The sha256 hash generated from secret (and timestamp if provided) |
| --time-lock               | string   | true     |         | The number of blocks to wait before the asset may be returned to  |
| --timestamp               | uint     |          |         | The timestamp in seconds for generating hashLock if provided      |

## Examples

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
