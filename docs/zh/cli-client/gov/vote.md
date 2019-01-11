# iriscli gov vote

## 描述

给VotingPeriod状态的提议投票, 可选项包括: Yes/No/NoWithVeto/Abstain

## 使用方式

```
iriscli gov vote [flags]
```

打印帮助信息:

```
iriscli gov vote --help
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

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须 |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --option         |                            | [string] 投票选项 {Yes, No, NoWithVeto, Abstain}                                                                                                  | Yes      |
| --proposal-id    |                            | [string] 投票的提议ID                                                                                                            | Yes      |

## 例子

### 给提议投票

```shell
iriscli gov vote --chain-id=test --proposal-id=1 --option=Yes --from node0 --fee=0.01iris
```

输入正确的密码之后，你就完成了对于所指定的提议投票。
注意：验证人和委托人才能对指定提议投票，并且可投票的提议必须是'VotingPeriod'状态。

```txt
Committed at block 43 (tx hash: 01C4C3B00C6048A12AE2CF2294F63C55A69011381B819C35F11B04C921DB81CC, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 2048,
   "codespace": "",
   "tags": {
     "action": "vote",
     "proposal-id": "2",
     "voter": "faa1x25y3ltr4jvp89upymegvfx7n0uduz5kmh5xuz"
   }
 })
```

如何查询投票详情？

请点击下述链接：

[query-vote](query-vote.md)

[query-votes](query-votes.md)

[query-tally](query-tally.md)
