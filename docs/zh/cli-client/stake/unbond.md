# iriscli stake unbond

## 描述

从指定的验证者解绑一定的股份

## 用法

```
iriscli stake unbond [flags]
```

## 标志

| 名称, 速记             | 默认值                | 描述                                                                                                           | 必需     |
| --------------------- | --------------------- | ------------------------------------------------------------------------------------------------------------- | -------- |
| --account-number      |                       | [int] 用来签名交易的accountNumber                                                                              |          |
| --address-validator   |                       | [string] 验证者bech地址                                                                                        | Yes      |
| --async               |                       | 是否异步广播交易                                                                                                |          |
| --chain-id            |                       | [string] Tendermint节点的链ID                                                                                  | Yes      |
| --dry-run             |                       | 忽略--gas标志并进行本地的交易仿真                                                                                |          |
| --fee                 |                       | [string] 交易费用                                                                                              | Yes      |
| --from                |                       | [string] 用来签名的私钥名                                                                                      | Yes      |
| --from-addr           |                       | [string] [string] 指定generate-only模式下的from address                                                        |          |
| --gas                 | 200000                | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                                            |          |
| --gas-adjustment      | 1                     | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略                            |          |
| --generate-only       |                       | 构建一个未签名交易并将其打印到标准输出                                                                            |          |
| --help, -h            |                       | unbond命令帮助                                                                                                 |          |
| --indent              |                       | 在JSON响应中添加缩进                                                                                           |          |
| --json                |                       | 以json形式输出                                                                                                 |          |
| --ledger              |                       | 使用连接的硬件记账设                                                                                           |          |
| --memo                |                       | [string] 发送交易时的备忘录                                                                                    |          |
| --node                | tcp://localhost:26657 | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                                                            |          |
| --print-response      |                       | 返回交易响应 (当且仅当同步模式下使用)                                                                            |          |
| --sequence            |                       | [int] 用来签名交易的sequence                                                                                   |          |
| --shares-amount       |                       | [string] 指定解绑的股份数                                                                                      |          |
| --shares-percent      |                       | [string] 指定解绑的股份比例，为0到1之间的正数                                                                    |          |
| --trust-node          | true                  | 关闭响应结果校验                                                                                               |          |

## 例子

### 从指定的验证者解绑一定的股份

```shell
iriscli stake unbond --address-validator=ValidatorAddress --shares-percent=SharePercent --from=UnbondInitiator --chain-id=ChainID --fee=Fee
```

运行成功以后，返回的结果如下：

```txt
Committed at block 851 (tx hash: A82833DE51A4127BD5D60E7F9E4CD5895F97B1B54241BCE272B68698518D9D2B, response: {Code:0 Data:[11 8 230 225 179 223 5 16 249 233 245 21] Log:Msg 0:  Info: GasWanted:200000 GasUsed:16547 Tags:[{Key:[97 99 116 105 111 110] Value:[98 101 103 105 110 45 117 110 98 111 110 100 105 110 103] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 108 101 103 97 116 111 114] Value:[102 97 97 49 51 108 99 119 110 120 112 121 110 50 101 97 51 115 107 122 109 101 107 54 52 118 118 110 112 57 55 106 115 107 56 113 109 104 108 54 118 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[115 111 117 114 99 101 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 54 100 54 116 104 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[101 110 100 45 116 105 109 101] Value:[11 8 230 225 179 223 5 16 249 233 245 21] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 56 50 55 51 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "begin-unbonding",
     "completeConsumedTxFee-iris-atto": "\"8273500000000000\"",
     "delegator": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
     "end-time": "\u000b\u0008\ufffd\ufffd\ufffd\ufffd\u0005\u0010\ufffd\ufffd\ufffd\u0015",
     "source-validator": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd"
   }
 }

```
