# iriscli asset mint-token

## 描述

资产的所有者增发一定数量的代币到指定账户地址

## 使用方式

```shell
iriscli asset mint-token <token-id> [flags]
```

## 特有的标志

| Name | Type | Required | Default | Description                                              |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --to    | string | 否 | "" | address of mint token to |
| --amount | uint64 | 是 | 0 | amount of mint token |


## 示例

```shell
iriscli asset mint-token btc --amount=1000000 --from=node0 --chain-id=irishub-test --fee=0.4iris
```

输出信息:
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
