# iriscli service call 

## 描述

调用服务方法

## 用法

```
iriscli service call <flags>
```

## 特殊标志

| Name, shorthand       | Default                 | Description                          | Required |
| --------------------- | ----------------------- | ------------------------------------ | -------- |
| --def-chain-id        |                         | 定义该服务的区块链ID           | 是       |
| --service-name        |                         | 服务名称                     | 是       |
| --method-id           |                         | 调用的服务方法ID                 | 是       |
| --bind-chain-id       |                         | 绑定该服务的区块链ID           | 是       |
| --provider            |                         | bech32编码的服务提供商账户地址  | 是       |
| --service-fee         |                         | 服务调用支付的服务费            |          |
| --request-data        |                         | hex编码的服务调用请求数据        |          |

## 示例

### 发起一个服务调用请求

```shell
iriscli service call --chain-id=<chain-id> --from=<key_name> --fee=0.3iris --def-chain-id=<service_define_chain_id> --service-name=<service_name> --method-id=1 --bind-chain-id=<service_bind_chain_id> --provider=<provider_address> --service-fee=1iris --request-data=<request-data>
```

运行成功以后，返回的结果如下:

```txt
Committed at block 54 (tx hash: F972ACA7DF74A6C076DFB01E7DD49D8694BF5AA1BA25A1F1B875113DFC8857C3, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 7614,
   "codespace": "",
   "tags": {
     "action": "service_call",
     "consumer": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "provider": "iaa1x25y3ltr4jvp89upymegvfx7n0uduz5krcj7ul",
     "request-id": "64-54-0"
   }
 })
```

