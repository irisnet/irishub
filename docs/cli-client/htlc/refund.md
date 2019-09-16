# iriscli htlc refund

## Description

Refund from an expired HTLC

## Usage

```bash
iriscli htlc refund --hash-lock=<hash-lock>
```

## Flags

| Name, shorthand | Type     | Required | Default | Description                                       |
| --------------- | -------- | -------- | ------- | ------------------------------------------------- |
| --hash-lock     | bytesHex | true     |         | The hash lock identifying the HTLC to be refunded |

## Examples

### Refund from an expired HTLC

```bash
iriscli htlc refund \
--from=userX \
--hash-lock=f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20 \
--fee=0.3iris \
--chain-id=testNet \
--commit
```
