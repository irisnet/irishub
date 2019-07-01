# iriscli asset transfer-gateway-owner

## 描述

转让指定Gateway的所有权到另一个地址。

## 使用方式

```
iriscli asset transfer-gateway-owner [flags]
```

打印帮助信息:
```
iriscli asset transfer-gateway-owner --help
```

## 特定的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值   | 描述                                                       |
| --------------------| -----  | -------- | --------|-------------------------------------------------------- |
| --moniker           | string  | 是     | ""       | the unique name of the gateway to be transferred       |
| --to                | Address | 是     |          | the new owner to which the gateway will be transferred |


## 示例

```
iriscli asset transfer-gateway-owner --moniker=tgw --to=faa1anyffkfjhyvdawv2trlv4eq00c3xrjmqvldwal --chain-id=irishub 
--from=node0 --fee=0.3iris
```

输出信息:
```json
{
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 50000,
   "gas_used": 4654,
   "codespace": "",
   "tags": [
     {
       "key": "action",
       "value": "transfer_gateway_owner"
     },
     {
       "key": "moniker",
       "value": "tgw"
     }
   ]
}
```
