# iriscli service bind 

## 描述

创建一个新的服务绑定

## 用法

```
iriscli service bind [flags]
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

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int]  服务平均返回时间的毫秒数表示                                                     | 是       |
| --bind-type           |                         | [string] 对服务是本地还是全局的设置，可选值Local/Global                                  | 是       |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                          | 是       |
| --deposit             |                         | [string] 服务提供者的保证金                                                            | 是       |
| --prices              |                         | [strings] 服务定价,按照服务方法排序的定价列表                                             |          |
| --service-name        |                         | [string] 服务名称                                                                    | 是       |
| --usable-time         |                         | [int] 每一万次服务调用的可用性的整数表示                                                  | 是       |

## 示例

### 添加服务绑定到已存在的服务定义
在服务绑定中, 你需要抵押`deposit`指定数量的IRIS, 最小的抵押金额为 该服务的服务费价格`price` * genesis 中定义的`MinDepositMultiple`倍数

```shell
iriscli service bind --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1000iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

运行成功以后，返回的结果如下:

```txt
Committed at block 567 (tx hash: A48DBD217CBB843E72CC47B40F90CE7DEEEDD6437C86A74A2976ADC9F449A034, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 5893,
   "codespace": "",
   "tags": {
     "action": "service_bind"
   }
 })
```

