# iriscli gov submit-proposal

## 描述

提交区块链治理提议以及发起提议所涉及的初始保证金，其中提议的类型包括Text/ParameterChange/SoftwareUpgrade这三种类型。

## 使用方式

```
iriscli gov submit-proposal [flags]
```
打印帮助信息:

```
iriscli gov submit-proposal --help
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

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] 提议的保证金                                                                                                                         |          |
| --description    |                            | [string] 提议的描述                                                                                                           | Yes      |
| --key            |                            | 参数的键名称                                                                                                                        |          |
| --op             |                            | [string] 对参数的操作                                                                                                             |          |
| --param          |                            | [string] 提议参数,例如: [{key:key,value:value,op:update}]                                                                                 |          |
| --path           |                            | [string] param.json文件路径                                                                                                                      |          |
| --title          |                            | [string] 提议标题                                                                                                                           | Yes      |
| --type           |                            | [string] 提议类型,例如:ParameterChange/SoftwareUpgrade/SoftwareHalt/TxTaxUsage                                                                            | Yes      |
| --version           |            0                | [uint64] 新协议的版本信息                                                                           |       |
| --software           |           " "                 | [string] 新协议的软件地址                                                                       |       |
| --switch-height           |       0                     | [string] 新版本协议升级的高度                                                     |       |
| --threshold        | "0.85"   |  [string] 软件升级的阈值                                              |               |

## 例子

### 提交一个'ParameterChange'类型的提议

```shell
iriscli gov submit-proposal --chain-id=test --title="update MinDeposit proposal" --param='mint/Inflation=0.050' --type=ParameterChange --description="a new parameter change proposal" --from=node0 --fee=0.01iris
```

提交之后，您完成了提交新的“ParameterChange”提议。
更改参数的详细信息（通过查询参数获取参数，修改它，然后在“操作”上添加“更新”，使用方案中的更多详细信息）和其他类型的提议字段与文本提议类似。
注意：在这个例子中, --path 和 --param 不能同时为空。

### 提交一个'SoftwareUpgrade'类型的提议

```shell
iriscli gov submit-proposal --chain-id=test --title="irishub0.7.0 upgrade proposal" --type=SoftwareUpgrade --description="a new software upgrade proposal" --from=node0 --fee=0.01iris --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80 --threshold=0.9
```

在这种场景下，提议的 --title、--type 和--description参数必不可少，另外你也应该保留好提议ID，这是检索所提交提议的唯一方法。


如何查询提议详情？

请点击下述链接：

[query-proposal](query-proposal.md)

[query-proposals](query-proposals.md)
