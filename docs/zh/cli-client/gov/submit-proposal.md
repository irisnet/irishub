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

| 名称, 速记        | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] 提议的保证金                                                                                                                         |          |
| --description    |                            | [string] 提议的描述                                                                                                           | Yes      |
| --key            |                            | 参数的键名称                                                                                                                        |          |
| --op             |                            | [string] 对参数的操作                                                                                                             |          |
| --param          |                            | [string] 提议参数,例如: [{key:key,value:value,op:update}]                                                                                 |          |
| --path           |                            | [string] param.json文件路径                                                                                                                      |          |
| --title          |                            | [string] 提议标题                                                                                                                           | Yes      |
| --type           |                            | [string] 提议类型,例如:Text/ParameterChange/SoftwareUpgrade                                                                            | Yes      |
| --version           |            0                | [uint64] 新协议的版本信息                                                                           |       |
| --software           |           " "                 | [string] 新协议的软件地址                                                                       |       |
| --switch-height           |       0                     | [string] 新版本协议升级的高度                                                     |       |

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
iriscli gov submit-proposal --chain-id=test --title="irishub0.7.0 upgrade proposal" --type=SoftwareUpgrade --description="a new software upgrade proposal" --from=node0 --fee=0.01iris --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80
```

在这种场景下，提议的 --title、--type 和--description参数必不可少，另外你也应该保留好提议ID，这是检索所提交提议的唯一方法。


如何查询提议详情？

请点击下述链接：

[query-proposal](query-proposal.md)

[query-proposals](query-proposals.md)
