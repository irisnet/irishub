# iriscli stake delegate

## 描述

委托一定代币到某个验证者

## 用法

```
iriscli stake delegate [flags]
```

## 标志

| 名称, 速记                    | 默认值                | 描述                                                                | 必需     |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] 用来签名交易的accountNumber                                    |          |
| --address-validator          |                       | [string] 验证者bech地址                                              | Yes      |
| --amount                     |                       | [string] 绑定的代币数量                                              | Yes      |
| --async                      |                       | 是否异步广播交易                                                     |          |
| --chain-id                   |                       | [string] Tendermint节点的链ID                                        | Yes      |
| --dry-run                    |                       | 忽略--gas标志并进行本地的交易仿真                                      |          |
| --fee                        |                       | [string] 交易费用                                                    | Yes      |
| --from                       |                       | [string] 用来签名的私钥名                                             | Yes      |
| --from-addr                  |                       | [string] 指定generate-only模式下的from address                        |          |
| --gas                        | 200000                | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值         |          |
| --gas-adjustment             | 1                     | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only              |                       | 构建一个未签名交易并将其打印到标准输出                                 |          |
| --help, -h                   |                       | delegate命令帮助                                                    |          |
| --indent                     |                       | 在JSON响应中添加缩进                                                 |          |
| --json                       |                       | 以json形式输出                                                       |          |
| --ledger                     |                       | 使用连接的硬件记账设备                                               |          |
| --memo                       |                       | [string] 发送交易时的备忘录                                           |          |
| --node                       | tcp://localhost:26657 | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                   |          |
| --print-response             |                       | 返回交易响应 (当且仅当同步模式下使用)                                  |          |
| --sequence int               |                       | [int] 用来签名交易的sequence                                         |          |
| --trust-node                 | true                  | 关闭响应结果校验                                                     |          |

## 例子

### 委托一定代币到某个验证者

```shell
iriscli stake delegate --chain-id=ChainID --from=KeyName --fee=Fee --amount=CoinstoBond --address-validator=ValidatorOwnerAddress
```

运行成功以后，返回的结果如下：

```txt
Committed at block 2352 (tx hash: 2069F0453619637EE67EFB0CFC53713AF28A0BB89137DEB4574D8B13E723999B, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:15993 Tags:[{Key:[97 99 116 105 111 110] Value:[100 101 108 101 103 97 116 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 108 101 103 97 116 111 114] Value:[102 97 97 49 51 108 99 119 110 120 112 121 110 50 101 97 51 115 107 122 109 101 107 54 52 118 118 110 112 57 55 106 115 107 56 113 109 104 108 54 118 120] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 115 116 105 110 97 116 105 111 110 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 53 103 114 118 51 120 103 51 101 107 120 104 57 120 114 102 55 57 122 100 48 119 48 55 55 107 114 103 118 53 120 102 54 100 54 116 104 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 55 57 57 54 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "delegate",
     "completeConsumedTxFee-iris-atto": "\"7996500000000000\"",
     "delegator": "faa13lcwnxpyn2ea3skzmek64vvnp97jsk8qmhl6vx",
     "destination-validator": "fva15grv3xg3ekxh9xrf79zd0w077krgv5xf6d6thd"
   }
 }
```
