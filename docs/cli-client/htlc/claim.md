# iriscli htlc claim

## Description

Claim an opened HTLC

## Usage

```bash
iriscli htlc claim --hash-lock=<hash-lock> --secret=<secret>
```

## Flags

| Name, shorthand | Type     | Required | Default | Description                                      |
| --------------- | -------- | -------- | ------- | ------------------------------------------------ |
| --hash-lock     | bytesHex | true     |         | The hash lock identifying the HTLC to be claimed |
| --secret        | string   | true     |         | The secret for generating hash lock              |

## Examples

### Claim an opened HTLC

```bash
iriscli htlc claim \
--from=userX \
--hash-lock=f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20 \
--secret=___abcdefghijklmnopqrstuvwxyz___ \
--fee=0.3iris \
--chain-id=testNet \
--commit
```
