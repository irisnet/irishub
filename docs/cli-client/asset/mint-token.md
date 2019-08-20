# iriscli asset mint-token

## Description

The asset owner can directly mint tokens to a specified address

## Usage

```bash
iriscli asset mint-token <token-id> [flags]
```

## Flags

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --to    | string | No | "" | address of mint token to, default is your own address |
| --amount | uint64 | Yes | 0 | amount of the token to mint |

## Example

```bash
iriscli asset mint-token kitty --amount=1000000 --from=<key-name> --chain-id=irishub --fee=0.4iris
```
