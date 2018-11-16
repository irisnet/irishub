# iriscli gov vote

## 描述

给VotingPeriod状态的提议投票, 可选项包括: Yes/No/NoWithVeto/Abstain

## 使用方式

```
iriscli gov vote [flags]
```

## 标志

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须 |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] 用来签名交易的AccountNumber                                                                                                            |          |
| --async          |                            | 异步广播交易                                                                                                                |          |
| --chain-id       |                            | [string] tendermint节点的链ID                                                                                                                 | Yes      |
| --dry-run        |                            | 忽略--gas标志并进行本地的交易仿真                                                              |          |
| --fee            |                            | [string] 支付的交易费用                                                                                                           | Yes      |
| --from           |                            | [string] 用来签名的私钥名                                                                                                      | Yes      |
| --from-addr      |                            | [string] 指定generate-only模式下的from address                                                                                                  |          |
| --gas            | 200000                     | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                                 |          |
| --gas-adjustment | 1                          | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only  |                            | 构建一个未签名交易并将其打印到标准输出                                                                                                 |          |
| --help, -h       |                            | 查询命令帮助                                                                                                                                        |          |
| --indent         |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --json           |                            | 输出将以json格式返回                                                                                                                         |          |
| --ledger         |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --memo           |                            | [string] 发送交易的备忘录                                                                                                         |          |
| --node           | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --option         |                            | [string] 投票选项 {Yes, No, NoWithVeto, Abstain}                                                                                                  | Yes      |
| --print-response |                            | 返回交易响应 (当且仅当同步模式下使用))                                                                                                   |          |
| --proposal-id    |                            | [string] 投票的提议ID                                                                                                            | Yes      |
| --sequence       |                            | [int] 用来签名交易的sequence number                                                                                                                 |          |
| --trust-node     | true                       | 关闭响应结果校验                                                                                                                    |          |

## 例子

### 给提议投票

```shell
iriscli gov vote --chain-id=test --proposal-id=1 --option=Yes --from node0 --fee=0.01iris
```

输入正确的密码之后，你就完成了对于所指定的提议投票。
注意：验证人和委托人才能对指定提议投票，并且可投票的提议必须是'VotingPeriod'状态。

```txt
Vote[Voter:faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd,ProposalID:1,Option:Yes]Password to sign with 'node0':
Committed at block 989 (tx hash: ABDD88AA3CA8F365AC499427477CCE83ADF09D7FC2D62643D0217107E489A483, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:2408 Tags:[{Key:[97 99 116 105 111 110] Value:[118 111 116 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[118 111 116 101 114] Value:[102 97 97 49 52 113 53 114 102 57 115 108 50 100 113 100 50 117 120 114 120 121 107 97 102 120 113 51 110 117 51 108 106 50 102 112 57 108 55 112 103 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 50 48 52 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "vote",
     "completeConsumedTxFee-iris-atto": "\"120400000000000\"",
     "proposal-id": "1",
     "voter": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd"
   }
 }
```

如何查询投票详情？

请点击下述链接：

[query-vote](query-vote.md)

[query-votes](query-votes.md)

[query-tally](query-tally.md)