# iriscli gov deposit

## 描述
 
充值保证金以激活提议
 
## 使用方式
 
```
iriscli gov deposit [flags]
```

## 标志
 
| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] 用来签名交易的AccountNumber                                                                                                            |          |
| --async          |                            | 异步广播交易                                                                                                                |          |
| --chain-id       |                            | [string] tendermint节点的链ID                                                                                                                 | Yes      |
| --deposit        |                            | [string] 发起提议的保证金                                                                                                                         | Yes      |
| --dry-run        |                            | 忽略--gas标志并进行本地的交易仿真                                                              |          |
| --fee            |                            | [string] 支付的交易费用                                                                                                           | Yes      |
| --from           |                            | [string] 用来签名的私钥名                                                                                                      | Yes      |
| --from-addr      |                            | [string] 指定generate-only模式下的from address                                                                                                  |          |
| --gas            | 200000                     | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                                 |          |
| --gas-adjustment | 1                          | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only  |                            | 构建一个未签名交易并将其打印到标准输出                                                                                                 |          |
| --help, -h       |                            | 查询命令帮助                                                                                                                                     |          |
| --indent         |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --json           |                            | 输出将以json格式返回                                                                                                                         |          |
| --ledger         |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --memo           |                            | [string] 发送交易的备忘录                                                                                                         |          |
| --node           | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --print-response |                            | 返回交易响应 (当且仅当同步模式下使用))                                                                                                  |          |
| --proposal-id    |                            | [string] 充值保证金的提议ID                                                                                                        | Yes      |
| --sequence       |                            | [int] 用来签名交易的sequence number                                                                                                                 |          |
| --trust-node     | true                       | 关闭响应结果校验                                                                                                                    |          |

## 例子

### 充值保证金

```shell
iriscli gov deposit --chain-id=test --proposal-id=1 --deposit=50iris --from=node0 --fee=0.01iris
```

输入正确的密码后，你就充值了50个iris用以激活提议的投票状态。

```txt
Password to sign with 'node0':
Committed at block 473 (tx hash: 0309E969589F19A9D9E4BFC9479720487FBF929ED6A88824414C5E7E91709206, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:6710 Tags:[{Key:[97 99 116 105 111 110] Value:[100 101 112 111 115 105 116] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 112 111 115 105 116 101 114] Value:[102 97 97 49 52 113 53 114 102 57 115 108 50 100 113 100 50 117 120 114 120 121 107 97 102 120 113 51 110 117 51 108 106 50 102 112 57 108 55 112 103 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 51 51 53 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "deposit",
     "completeConsumedTxFee-iris-atto": "\"335500000000000\"",
     "depositor": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
     "proposal-id": "1"
   }
 }
```

如何查询保证金充值明细？

请点击下述链接：

[query-deposit](query-deposit.md)

[query-deposits](query-deposits.md)