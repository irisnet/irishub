# iriscli record submit

## 描述

向链上提交一个存证

## 用法

```
iriscli record submit [flags]
```

## 标志

| 名称, 速记        | 默认值                     | 描述                                                                               | 必需      |
| ---------------  | -------------------------- | ---------------------------------------------------------------------------------- | -------- |
| --account-number |                            | [int] 用来签名交易的AccountNumber                                                   |          |
| --async          |                            | 异步广播交易                                                                        |          |
| --chain-id       |                            | [string] tendermint节点的链ID                                                       | 是       |
| --description    | description                | [string] 上传文件的描述信息                                                          |          |
| --dry-run        |                            | 忽略--gas标志并进行本地的交易仿真                                                     |          |
| --fee            |                            | [string] 支付的交易费用                                                              | 是       |
| --from           |                            | [string] 用来签名的私钥名                                                            | 是       |
| --from-addr      |                            | [string] 指定generate-only模式下的from address                                      |          |
| --gas string     | 200000                     | 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                |          |
| --gas-adjustment | 1                          | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only  |                            | 构建一个未签名交易并将其打印到标准输出                                                 |          |
| -h, --help       |                            | 提交命令帮助                                                                         |          |
| --indent         |                            | 在JSON响应中添加缩进                                                                 |          |
| --json           |                            | 输出将以json格式返回                                                                 |          |
| --ledger         |                            | 使用连接的硬件记账设备                                                                |          |
| --memo           |                            | [string] 发送交易的备忘录                                                            |          |
| --node           | tcp://localhost:26657      | [string] [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                  |          |
| --onchain-data   |                            | [string] 上链数据                                                                    | 是      |
| --print-response |                            | 返回交易响应 (当且仅当同步模式下使用))                                                 |          |
| --sequence       |                            | [int] 用来签名交易的sequence number                                                  |          |
| --trust-node     | true                       | 关闭响应结果校验                                                                     |          |

## 例子

### 提交存证

```shell
iriscli record submit --chain-id="test" --onchain-data="this is my on chain data" --from=node0 --fee=0.1iris
```

运行成功以后，返回的结果如下：

```txt
Password to sign with 'node0':
Committed at block 72 (tx hash: 7CCC8B4018D4447E6A496923944870E350A1A3AF9E15DB15B8943DAD7B5D782B, response: {Code:0 Data:[114 101 99 111 114 100 58 97 98 53 54 48 50 98 97 99 49 51 102 49 49 55 51 55 101 56 55 57 56 100 100 53 55 56 54 57 99 52 54 56 49 57 52 101 102 97 100 50 100 98 51 55 54 50 53 55 57 53 102 49 101 102 100 56 100 57 100 54 51 99 54] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4090 Tags:[{Key:[97 99 116 105 111 110] Value:[115 117 98 109 105 116 45 114 101 99 111 114 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[111 119 110 101 114 65 100 100 114 101 115 115] Value:[102 97 97 49 50 50 117 122 122 112 117 103 116 114 122 115 48 57 110 102 51 117 104 56 120 102 106 97 122 97 53 57 120 118 102 57 114 118 116 104 100 108] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[114 101 99 111 114 100 45 105 100] Value:[114 101 99 111 114 100 58 97 98 53 54 48 50 98 97 99 49 51 102 49 49 55 51 55 101 56 55 57 56 100 100 53 55 56 54 57 99 52 54 56 49 57 52 101 102 97 100 50 100 98 51 55 54 50 53 55 57 53 102 49 101 102 100 56 100 57 100 54 51 99 54] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 50 48 52 53 48 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "submit-record",
     "completeConsumedTxFee-iris-atto": "\"2045000000000000\"",
     "ownerAddress": "faa122uzzpugtrzs09nf3uh8xfjaza59xvf9rvthdl",
     "record-id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6"
   }
 }
```

本次存证操作的record-id如下:

```txt
"record-id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6"
```

请务必备份record-id，以备将来查询本次存证。若丢失record-id，本次存证再也无法查询到。
