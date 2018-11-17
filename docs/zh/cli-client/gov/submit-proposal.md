# iriscli gov submit-proposal

## 描述

提交区块链治理提议以及发起提议所涉及的初始保证金，其中提议的类型包括Text/ParameterChange/SoftwareUpgrade这三种类型。

## 使用方式

```
iriscli gov submit-proposal [flags]
```

## 标志

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] 用来签名交易的AccountNumber                                                                                                            |          |
| --async          |                            | 异步广播交易                                                                                                                |          |
| --chain-id       |                            | [string] tendermint节点的链ID                                                                                                                 | Yes      |
| --deposit        |                            | [string] 提议的保证金                                                                                                                         |          |
| --description    |                            | [string] 提议的描述                                                                                                                     | Yes      |
| --dry-run        |                            | 忽略--gas标志并进行本地的交易仿真                                                              |          |
| --fee            |                            | [string] 支付的交易费用                                                                                                           | Yes      |
| --from           |                            | [string] 用来签名的私钥名                                                                                                      | Yes      |
| --from-addr      |                            | [string] 指定generate-only模式下的from address                                                                                                  |          |
| --gas            | 200000                     | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                                 |          |
| --gas-adjustment | 1                          | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only  |                            | 构建一个未签名交易并将其打印到标准输出                                                                                                 |          |
| --help, -h       |                            | 查询命令帮助                                                                                                                             |          |
| --indent         |                            | 在JSON响应中添加缩进                                                                                                                          |          |
| --json           |                            | 输出将以json格式返回                                                                                                                         |          |
| --key            |                            | 参数的键名称                                                                                                                                 |          |
| --ledger         |                            | 使用连接的硬件记账设备                                                                                                                        |          |
| --memo           |                            | [string] 发送交易的备忘录                                                                                                         |          |
| --node           | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                                                                                  |          |
| --op             |                            | [string] 对参数的操作                                                                                                                  |          |
| --param          |                            | [string] 提议参数,例如: [{key:key,value:value,op:update}]                                                                                 |          |
| --path           |                            | [string] param.json文件路径                                                                                                                      |          |
| --print-response |                            | 返回交易响应 (当且仅当同步模式下使用))                                                                                                   |          |
| --sequence       |                            | [int] 用来签名交易的sequence number                                                                                                                 |          |
| --title          |                            | [string] 提议标题                                                                                                                           | Yes      |
| --trust-node     | true                       | 关闭响应结果校验                                                                                                                    |          |
| --type           |                            | [string] 提议类型,例如:Text/ParameterChange/SoftwareUpgrade                                                                            | Yes      |

## 例子

### 提交一个'Text'类型的提议

```shell
iriscli gov submit-proposal --chain-id=test --title="notice proposal" --type=Text --description="a new text proposal" --from=node0 --fee=0.01iris
```

输入正确的密码之后，你就完成提交了一个提议，需要注意的是要记下你的提议ID，这是可以检索你的提议的唯一要素。

```txt
Password to sign with 'node0':
Committed at block 44 (tx hash: 2C28A87A6262CACEDDB4EBBC60FE989D0DB2B7DEB1EC6795D2F4707DA32C7CBF, response: {Code:0 Data:[49] Log:Msg 0:  Info: GasWanted:200000 GasUsed:8091 Tags:[{Key:[97 99 116 105 111 110] Value:[115 117 98 109 105 116 45 112 114 111 112 111 115 97 108] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 101 114] Value:[102 97 97 49 115 108 116 106 120 100 103 107 48 48 115 56 54 50 57 50 122 48 99 110 55 97 53 100 106 99 99 116 54 101 115 115 110 97 118 100 121 122] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 97 114 97 109] Value:[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 52 48 52 53 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "submit-proposal",
     "completeConsumedTxFee-iris-atto": "\"4045500000000000\"",
     "param": "",
     "proposal-id": "1",
     "proposer": "faa1sltjxdgk00s86292z0cn7a5djcct6essnavdyz"
   }
 }
```

### 提交一个'ParameterChange'类型的提议

```shell
iriscli gov submit-proposal --chain-id=test --title="update MinDeposit proposal" --param='{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":20}","op":"update"}' --type=ParameterChange --description="a new parameter change proposal" --from=node0 --fee=0.01iris
```

提交之后，您完成了提交新的“ParameterChange”提议。
更改参数的详细信息（通过查询参数获取参数，修改它，然后在“操作”上添加“更新”，使用方案中的更多详细信息）和其他类型的提议字段与文本提议类似。
注意：在这个例子中, --path 和 --param 不能同时为空。

### 提交一个'SoftwareUpgrade'类型的提议

```shell
iriscli gov submit-proposal --chain-id=test --title="irishub0.7.0 upgrade proposal" --type=SoftwareUpgrade --description="a new software upgrade proposal" --from=node0 --fee=0.01iris
```

在这种场景下，提议的 --title、--type 和--description参数必不可少，另外你也应该保留好提议ID，这是检索所提交提议的唯一方法。


如何查询提议详情？

请点击下述链接：

[query-proposal](query-proposal.md)

[query-proposals](query-proposals.md)