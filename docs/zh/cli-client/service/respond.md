# iriscli service respond 

## 描述

响应服务调用

## 用法

```
iriscli service respond [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] 发起该服务调用的区块链ID                                                                                              | 是       |
| --request-id          |                         | [string] 该服务调用的ID                                                                                                                                | 是       |
| --response-data       |                         | [string] hex编码的服务调用响应数据                                                                       |        |
| -h, --help            |                         | 响应命令帮助                                                                                                                                         |          |

## 示例

### 响应一个服务调用 
```shell
iriscli service respond --chain-id=test --from=node0 --fee=0.004iris --request-chain-id=test --request-id=230-130-0 --response-data=abcd
```

运行成功以后，返回的结果如下:

```txt
Committed at block 71 (tx hash: C02BC5F4D6E74ED13D8D5A31F040B0FED0D3805AF1C546544A112DB2EFF3D9D5, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3784,
   "codespace": "",
   "tags": {
     "action": "service_respond",
     "consumer": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "provider": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "request-id": "78-68-0"
   }
 })
```

