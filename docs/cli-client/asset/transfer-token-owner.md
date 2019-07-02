# iriscli asset transfer-token-owner

## Description

Transfer control of the token

## Usage

```shell
iriscli asset transfer-token-owner <token-id> --to=<new owner>
```

## Flags

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --to           | string | Yes | "" | Destination address |

## Example

```shell
iriscli asset transfer-token-owner btc --to=faa1ze4nx2k4ehsu83hk3texpktrt07gsff24mjq8d --from=node0 --chain-id=irishub-test --fee=0.4iris --commit
```

The output:

```json
{
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 50000,
   "gas_used": 5120,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "transfer_token_owner"
     },
     {
       "key": "token-id",
       "value": "btc"
     }
   ]
 })
```



