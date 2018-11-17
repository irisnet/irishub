# iriscli stake edit-validator

## 描述

编辑已存在的验证者信息

## 用法

```
iriscli stake edit-validator [flags]
```

## 标志

| 名称, 速记                    | 默认值                | 描述                                                                | 必需     |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] 用来签名交易的accountNumber                                    |          |
| --address-delegator          |                       | [string] 委托者bech地址                                              |          |
| --amount                     |                       | [string] 绑定的代币数量                                     |          |
| --async                      |                       | 是否异步广播交易                                                     |          |
| --chain-id                   |                       | [string] Tendermint节点的链ID                                       | Yes      |
| --commission-max-change-rate |                       | [string] 最大佣金变化率(每天)                                        |          |
| --commission-max-rate        |                       | [string] 最大佣金率                                                 |          |
| --commission-rate            |                       | [string] [string] 初始佣金率                                         |          |
| --details                    |                       | [string] 可选details                                                 |          |
| --dry-run                    |                       | 忽略--gas标志并进行本地的交易仿真                                      |          |
| --fee                        |                       | [string] 交易费用                                                     | Yes      |
| --from                       |                       | [string] 用来签名的私钥名                                              | Yes      |
| --from-addr                  |                       | [string] 指定generate-only模式下的from address                        |          |
| --gas                        | 200000                | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值        |           |
| --gas-adjustment             | 1                     | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only              |                       | 构建一个未签名交易并将其打印到标准输出                                  |          |
| --genesis-format             |                       | 以gen-tx的格式导出交易; 其暗含 --generate-only                        |          |
| --help, -h                   |                       | edit-validator命令帮助                                               |          |
| --identity                   |                       | [string] 可选身份签名 (ex. UPort or Keybase)                         |          |
| --indent                     |                       | 在JSON响应中添加缩进                                                 |          |
| --ip                         |                       | [string] N节点的公有IP，仅开启--genesis-format时生效                  |           |
| --json                       |                       | 以json形式输出                                                       |          |
| --ledger                     |                       | 使用连接的硬件记账设备                                                |          |
| --memo                       |                       | [string] 发送交易时的备忘录                                          |          |
| --moniker                    |                       | [string] 验证者名字                                                 |          |
| --node                       | tcp://localhost:26657 | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                  |          |
| --node-id                    |                       | [string] 节点ID                                                     |          |
| --print-response             |                       | 返回交易响应 (当且仅当同步模式下使用)                                  |          |
| --pubkey                     |                       | [string] 验证者的Go-Amino编码的16进制公钥. 对于Ed25519，go-amino的16进制前缀为1624de6220 |           |
| --sequence                   |                       | [int] 用来签名交易的sequence                                         |          |
| --trust-node                 | true                  | 关闭响应结果校验                                                     |          |
| --website                    |                       | [string] 选填网站                                                    |          |

## 例子

### 编辑已存在的验证者信息

```shell
iriscli stake edit-validator --from=KeyName --chain-id=ChainID --fee=Fee --memo=YourMemo
```

运行成功以后，返回的结果如下：

```txt
Committed at block 2160 (tx hash: C48CABDA1183B5319003433EB1FDEBE5A626E00BD319F1A84D84B6247E9224D1, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3540 Tags:[{Key:[97 99 116 105 111 110] Value:[101 100 105 116 45 118 97 108 105 100 97 116 111 114] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 115 116 105 110 97 116 105 111 110 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 54 100 54 116 104 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[109 111 110 105 107 101 114] Value:[117 98 117 110 116 117 49 56] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[105 100 101 110 116 105 116 121] Value:[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 55 55 48 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "edit-validator",
     "completeConsumedTxFee-iris-atto": "\"177000000000000\"",
     "destination-validator": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd",
     "identity": "",
     "moniker": "ubuntu18"
   }
}
```
