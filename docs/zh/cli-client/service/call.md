# iriscli service call 

## 描述

调用服务方法

## 用法

```
iriscli service call [flags]
```
## 标志

| 命令，缩写       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | 否       |                       | 打印帮助                                                     |
| --chain-id       | String | 否       |                       | tendermint 节点网络ID                                        |
| --account-number | int    | 否       |                       | 账户数字用于签名通证发送                                     |
| --async          |        | 否       | True                  | 异步广播传输信息                                             |
| --dry-run        |        | 否       |                       | 忽略--gas 标志 ，执行仿真传输，但不广播。                    |
| --fee            | String | 是       |                       | 设置传输需要的手续费                                         |
| --from           | String | 是       |                       | 用于签名的私钥名称                                           |
| --gas            | String | 否       | 20000                 | 每笔交易设定的gas限额; 设置为“simulate”以自动计算所需气体    |
| --gas-adjustment | Float  | 否       | 1                     | 调整因子乘以传输模拟返回的估计值; 如果手动设置气体限制，则忽略该标志 |
| --generate-only  |        | 否       |                       | 创建一个未签名的传输并写到标准输出中。                       |
| --indent         |        | 否       |                       | 在JSON响应中增加缩进                                         |
| --json           |        | 否       |                       | 以json格式返回输出                                           |
| --memo           | String | 否       |                       | 传输中的备注信息                                             |
| --print-response |        | 否       |                       | 返回传输响应 (仅仅当 async = false时有效)                    |
| --sequence       | Int    | 否       |                       | 等待签名传输的序列号。                                       |
| --ledger         | String | 否       |                       | 使用一个联网的分账设备                                       |
| --node           | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。                    |
| --trust-node     | String | 否       | True                  | 不验证响应的证明                                             |
| --commit         | String | 否     | True                  |是否等到交易有明确的返回值，如果是True，则忽略--async的内容|

## 特殊标志

| Name, shorthand       | Default                 | Description                          | Required |
| --------------------- | ----------------------- | ------------------------------------ | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID           | 是       |
| --service-name        |                         | [string] 服务名称                     | 是       |
| --method-id           |                         | [int] 调用的服务方法ID                 | 是       |
| --bind-chain-id       |                         | [string] 绑定该服务的区块链ID           | 是       |
| --provider            |                         | [string] bech32编码的服务提供商账户地址  | 是       |
| --service-fee         |                         | [string] 服务调用支付的服务费            |          |
| --request-data        |                         | [string] hex编码的服务调用请求数据        |          |

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

