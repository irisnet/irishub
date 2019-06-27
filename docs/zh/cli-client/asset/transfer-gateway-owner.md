# iriscli asset transfer-gateway-owner

## 描述

转让指定Gateway的所有权到另一个地址。此命令仅用于产生未签名交易。当前的所有者和新的所有者需要依次对此交易进行签名，并将双重签名后的交易进行广播。签名和广播命令分别为"iriscli tx sign" 和 "iriscli tx broadcast"

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
| --------------------| -----  | -------- | --------  -------------------------------------------------------- |
| --moniker           | string  | 是     | ""       | the unique name of the gateway to be transferred       |
| --to                | Address | 是     |          | the new owner to which the gateway will be transferred |


## 示例

```
iriscli asset transfer-gateway-owner --moniker=tgw --to=faa1anyffkfjhyvdawv2trlv4eq00c3xrjmqvldwal --chain-id=irishub 
--from=node0 --fee=0.3iris
```

输出信息:
```txt
{"type":"irishub/bank/StdTx","value":{"msg":[{"type":"irishub/asset/MsgTransferGatewayOwner","value":{"owner":"faa1an4wfvsnxrp97lug5fngct6melhgcuvdv2qye3","moniker":"testgw","to":"faa1anyffkfjhyvdawv2trlv4eq00c3xrjmqvldwal"}}],"fee":{"amount":[{"denom":"iris-atto","amount":"600000000000000000"}],"gas":"50000"},"signatures":null,"memo":""}}
```
