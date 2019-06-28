# iriscli asset mint-token

## Description

The asset owner and operator can directly mint tokens to a specified address

## Usage

```shell
iriscli asset mint-token <token-id> [flags]
```

## Flags

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --to    | string | No | "" | address of mint token to |
| --amount | uint64 | Yes | 0 | amount of mint token |


## Example

```shell
iriscli asset mint-token btc --amount=1000000 --from=node0 --chain-id=irishub-test --fee=0.4iris
```

Output:
```json
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 50000,
   "gas_used": 6632,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "mint_token"
     },
     {
       "key": "recipient",
       "value": "faa1sf4xrfq3p45hlelp5pnw5rf4llsfg4st07mhjc"
     }
   ]
 }
```
