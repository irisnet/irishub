# iriscli stake redelegate

## 描述

重新委托一定的token从一个验证者到另一个验证者

## 用法

```
iriscli stake redelegate [flags]
```

## 标志

| 名称, 速记                    | 默认值                | 描述                                                                | 必需     |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] 用来签名交易的accountNumber                                    |          |
| --address-validator-dest     |                       | [string] 目标验证者bech地址                                          | Yes      |
| --address-validator-source   |                       | [string] 源验证者bech地址                                            | Yes      |
| --async                      |                       | 是否异步广播交易                                                     |          |
| --chain-id                   |                       | [string] Tendermint节点的链ID                                        | Yes      |
| --dry-run                    |                       | 忽略--gas标志并进行本地的交易仿真                                      |          |
| --fee                        |                       | [string] 交易费用                                                    | Yes      |
| --from                       |                       | [string] 用来签名的私钥名                                             | Yes      |
| --from-addr                  |                       | [string] 指定generate-only模式下的from address                       |          |
| --gas                        | 200000                | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值        |          |
| --gas-adjustment             | 1                     | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 ||
| --generate-only              |                       | 构建一个未签名交易并将其打印到标准输出                                 |          |
| --help, -h                   |                       | redelegate命令帮助                                                  |          |
| --indent                     |                       | 在JSON响应中添加缩进                                                 |          |
| --json                       |                       | 以json形式输出                                                       |          |
| --ledger                     |                       | 使用连接的硬件记账设备                                                |          |
| --memo                       |                       | [string] 发送交易时的备忘录                                           |          |
| --node                       | tcp://localhost:26657 | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                   |          |
| --print-response             |                       | 返回交易响应 (当且仅当同步模式下使用)                                  |          |
| --sequence                   |                       | [int] 用来签名交易的sequence                                          |          |
| --shares-amount              |                       | [string] 指定重新委托给其他验证者的股份                                |           |
| --shares-percent             |                       | [string] 指定重新委托给其他验证者的股份比例，位于0到1之间的正数          |          |
| --trust-node                 | true                  | 关闭响应结果校验                                                      |          |

## 例子

### 重新委托一定的token从一个验证者到另一个验证者

```shell
iriscli stake redelegate --chain-id=ChainID --from=KeyName --fee=Fee --address-validator-source=SourceValidatorAddress --address-validator-dest=DestinationValidatorAddress --shares-percent=SharesPercent
```

运行成功以后，返回的结果如下：

```txt
Committed at block 648 (tx hash: E59EE3C8F04D62DA0F5CFD89AC96402A92A56728692AEA47E8A126CDDA58E44B, response: {Code:0 Data:[11 8 185 204 185 223 5 16 247 169 147 42] Log:Msg 0:  Info: GasWanted:200000 GasUsed:29085 Tags:[{Key:[97 99 116 105 111 110] Value:[98 101 103 105 110 45 114 101 100 101 108 101 103 97 116 105 111 110] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 108 101 103 97 116 111 114] Value:[102 97 97 49 48 115 48 97 114 113 57 107 104 112 108 48 99 102 122 110 103 51 113 103 120 99 120 113 48 110 121 54 104 109 99 57 115 121 116 106 102 107] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[115 111 117 114 99 101 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 100 97 121 117 106 100 102 110 120 106 103 103 100 53 121 100 108 118 118 103 107 101 114 112 50 115 117 112 107 110 116 104 97 106 112 99 104 50] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 115 116 105 110 97 116 105 111 110 45 118 97 108 105 100 97 116 111 114] Value:[102 118 97 49 104 50 55 120 100 119 54 116 57 108 53 106 103 118 117 110 55 54 113 100 117 52 53 107 103 114 120 57 108 113 101 100 101 56 104 112 99 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[101 110 100 45 116 105 109 101] Value:[11 8 185 204 185 223 5 16 247 169 147 42] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 53 56 49 55 48 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "begin-redelegation",
     "completeConsumedTxFee-iris-atto": "\"5817000000000000\"",
     "delegator": "faa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9sytjfk",
     "destination-validator": "fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd",
     "end-time": "\u000b\u0008\ufffd̹\ufffd\u0005\u0010\ufffd\ufffd\ufffd*",
     "source-validator": "fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2"
   }
}
```
