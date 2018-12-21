# iriscli service call 

## 描述

调用服务方法

## 用法

```
iriscli service call [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                          | Required |
| --------------------- | ----------------------- | ------------------------------------ | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID           | 是       |
| --service-name        |                         | [string] 服务名称                     | 是       |
| --method-id           |                         | [int] 调用的服务方法ID                 | 是       |
| --bind-chain-id       |                         | [string] 绑定该服务的区块链ID           | 是       |
| --provider            |                         | [string] bech32编码的服务提供商账户地址  | 是       |
| --service-fee         |                         | [string] 服务调用支付的服务费            |          |
| --request-data        |                         | [string] hex编码的服务调用请求数据        |          |
| -h, --help            |                         | 调用命令帮助                             |          |

## 示例

### 发起一个服务调用请求
```shell
iriscli service call --chain-id=test --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service --method-id=1 --bind-chain-id=test --provider=faa1qm54q9ta97kwqaedz9wzd90cacdsp6mq54cwda --service-fee=1iris --request-data=434355
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
     "consumer": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "provider": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz",
     "request-id": "64-54-0"
   }
 })
```

