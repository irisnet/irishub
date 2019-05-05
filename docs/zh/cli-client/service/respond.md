# iriscli service respond 

## 描述

响应服务调用

## 用法

```
iriscli service respond <flags>
```

## 特有标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | 发起该服务调用的区块链ID                                                                                              | 是       |
| --request-id          |                         | 该服务调用的ID                                                                                                                                | 是       |
| --response-data       |                         | hex编码的服务调用响应数据                                                                       |        |

## 示例

### 响应一个服务调用 

```shell
iriscli service respond --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --request-chain-id=<request_chain_id> --request-id=<request-id> --response-data=<response-data>
```
>  `request-id` 可以从[service call](call.md)的返回中得到。

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
     "consumer": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "provider": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "request-id": "78-68-0"
   }
 })
```

