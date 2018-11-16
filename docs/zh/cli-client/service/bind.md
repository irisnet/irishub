# iriscli service bind 

## 描述

创建一个新的服务绑定

## 用法

```
iriscli service bind [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int]  服务平均返回时间的毫秒数表示                                                     |  Yes     |
| --bind-type           |                         | [string] 对服务是本地还是全局的设置，可选值Local/Global                                  |  Yes     |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                          |  Yes     |
| --deposit             |                         | [string] 服务提供者的保证金                                                            |  Yes     |
| --prices              |                         | [strings] 服务定价,按照服务方法排序的定价列表                                             |          |
| --service-name        |                         | [string] 服务名称                                                                    |  Yes     |
| --usable-time         |                         | [int] 每一万次服务调用的可用性的整数表示                                                  |  Yes     |
| -h, --help            |                         | 绑定命令帮助                                                                          |          |
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

### 添加服务绑定到已存在的服务定义
```shell
iriscli service bind --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

运行成功以后，返回的结果如下:

```txt
Password to sign with 'node0':
Committed at block 6 (tx hash: 87A477AEA41B22F7294084B4794837211C43A297D73EABA2F42F6436F3D975DD, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:5568 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 98 105 110 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 49 49 51 54 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-bind",
     "completeConsumedTxFee-iris-atto": "\"111360000000000\""
   }
 }
```

