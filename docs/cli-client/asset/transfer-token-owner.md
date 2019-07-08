# iriscli asset transfer-token-owner

## Description

Transfer the ownership of a token

## Usage

```bash
iriscli asset transfer-token-owner <token-id> [flags]
```

## Flags

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --to           | string | True | "" | the new owner address |

## Example

```bash
iriscli asset transfer-token-owner kitty --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```
