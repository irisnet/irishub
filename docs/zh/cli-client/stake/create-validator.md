# iriscli stake create-validator

## 描述

以自委托的方式创建一个新的验证者 

## 用法

```
iriscli stake create-validator [flags]
```

## 标志

| 名称, 速记                    | 默认值                | 描述                                                                 | 必需     |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] 用来签名交易的accountNumber                                    |          |
| --address-delegator          |                       | [string] 委托者bech地址                                              |          |
| --amount                     |                       | [string] 绑定的代币数量                                               | Yes      |
| --async                      |                       | 是否异步广播交易                                                      |          |
| --chain-id                   |                       | [string] Tendermint节点的链ID                                        | Yes      |
| --commission-max-change-rate |                       | [string] 最大佣金变化率(每天)                                         | Yes      |
| --commission-max-rate        |                       | [string] 最大佣金率                                                  | Yes      |
| --commission-rate            |                       | [string] 初始佣金率                                                  | Yes      |
| --details                    |                       | [string] 可选details                                                 |          |
| --dry-run                    |                       | 忽略--gas标志并进行本地的交易仿真                                      |          |
| --fee                        |                       | [string] 交易费用                                                    | Yes      |
| --from                       |                       | [string] 用来签名的私钥名                                             | Yes      |
| --from-addr                  |                       | [string] 指定generate-only模式下的from address                        |          |
| --gas                        | 200000                | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值         |       |
| --gas-adjustment             | 1                     | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 |          |
| --generate-only              |                       | 构建一个未签名交易并将其打印到标准输出                                   |          |
| --genesis-format             |                       | 以gen-tx的格式导出交易; 其暗含 --generate-only                          |          |
| --help, -h                   |                       | create-validator命令帮助                                               |          |
| --identity                   |                       | [string] 可选身份签名 (ex. UPort or Keybase)                           |          |
| --indent                     |                       | 在JSON响应中添加缩进                                                   |          |
| --ip                         |                       | [string] N节点的公有IP，仅开启--genesis-format时生效                    |           |
| --json                       |                       | 以json形式输出                                                         |          |
| --ledger                     |                       | 使用连接的硬件记账设备                                                  |          |
| --memo                       |                       | [string] 发送交易时的备忘录                                             |          |
| --moniker                    |                       | [string] 验证者名字                                                    |          |
| --node                       | tcp://localhost:26657 | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                     |          |
| --node-id                    |                       | [string] 节点ID                                                        |          |
| --print-response             |                       | 返回交易响应 (当且仅当同步模式下使用)                                     |          |
| --pubkey                     |                       | [string] 验证者的Go-Amino编码的16进制公钥. 对于Ed25519，go-amino的16进制前缀为1624de6220 | Yes       |
| --sequence                   |                       | [int] 用来签名交易的sequence                                            |          |
| --trust-node                 | true                  | 关闭响应结果校验                                                        |          |
| --website                    |                       | [string] 选填的网站                                                    |          |

## 例子

### 创建新的验证者

```shell
iriscli stake create-validator --chain-id=ChainID --from=KeyName --fee=Fee --pubkey=ValidatorPublicKey --commission-max-change-rate=CommissionMaxChangeRate --commission-max-rate=CommissionMaxRate --commission-rate=CommissionRate --amount=Coins
```

运行上述命令之后，你便成功地创建了一个验证者。
