# iriscli service refund-deposit 

## 描述

取回所有押金

## 用法

```
iriscli service refund-deposit [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                         |  Yes     |
| --service-name        |                         | [string] 服务名称                                                                   |  Yes     |
| -h, --help            |                         | 取回押金命令帮助                                                                      |          |
| --account-number      |                         | [int] 用来签名交易的AccountNumber                                                     |          |
| --async               |                         | 异步广播交易                                                                         |          |
| --chain-id            |                         | [string] tendermint节点的链ID                                                       | 是       |
| --dry-run             |                         | 忽略--gas标志并进行本地的交易仿真                                                       |          |
| --fee                 |                         | [string] 支付的交易费用                                                              | 是       |
| --from                |                         | [string] 用来签名的私钥名                                                            | 是       |
| --from-addr           |                         | [string] 指定generate-only模式下的from address                                       |          |
| --gas string          | 200000                  | 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                   |          |
| --gas-adjustment      | 1                       | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略       |          |
| --generate-only       |                         | 构建一个未签名交易并将其打印到标准输出                                                   |          |
| --indent              |                         | 在JSON响应中添加缩进                                                                 |          |
| --json                |                         | 输出将以json格式返回                                                                 |          |
| --ledger              |                         | 使用连接的硬件记账设备                                                                |          |
| --memo                |                         | [string] 发送交易的备忘录                                                            |          |
| --node                | tcp://localhost:26657   | [string] [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                   |          |
| --print-response      |                         | 返回交易响应 (当且仅当同步模式下使用))                                                  |          |
| --sequence            |                         | [int] 用来签名交易的sequence number                                                 |          |
| --trust-node          | true                    | 关闭响应结果校验                                                                    |          |

## 例子

### 取回所有押金
```shell
iriscli service refund-deposit --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

运行成功以后，返回的结果如下:

```txt
Password to sign with 'node0':
Committed at block 991 (tx hash: 8A7F0EA61AB73A8B241945C8942EC8593774346B36BB70E36E138A23E7A473EE, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:4614 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 114 101 102 117 110 100 45 100 101 112 111 115 105 116] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 57 50 50 56 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-refund-deposit",
     "completeConsumedTxFee-iris-atto": "\"92280000000000\""
   }
 }
```