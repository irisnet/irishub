# iriscli service respond 

## 描述

响应服务调用

## 用法

```
iriscli service respond [flags]
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

## 特有标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --request-chain-id    |                         | [string] 发起该服务调用的区块链ID                                                                                              | 是       |
| --request-id          |                         | [string] 该服务调用的ID                                                                                                                                | 是       |
| --response-data       |                         | [string] hex编码的服务调用响应数据                                                                       |        |

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

