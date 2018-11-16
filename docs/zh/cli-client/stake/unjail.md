# iriscli stake unjail

## 描述

恢复之前由于宕机被惩罚的验证者的身份

## 用法

```
iriscli stake redelegate [flags]
```

## 标志

| 名称, 速记                    | 默认值                | 描述                                                                | 必需     |
| ---------------------------- | --------------------- | ------------------------------------------------------------------- | -------- |
| --account-number             |                       | [int] 用来签名交易的accountNumber                                    |          |
| --address-validator-dest     |                       | [string] 目标验证者bech地址                                          |          |
| --address-validator-source   |                       | [string] 源验证者bech地址                                            |          |
| --async                      |                       | 是否异步广播交易                                                     |          |
| --chain-id                   |                       | [string] Tendermint节点的链ID                                       | Yes      |
| --dry-run                    |                       | 忽略--gas标志并进行本地的交易仿真                                      |          |
| --fee                        |                       | [string] 交易费用                                                    | Yes      |
| --from                       |                       | [string] 用来签名的私钥名                                             | Yes      |
| --from-addr                  |                       | [string] 指定generate-only模式下的from address                       |          |
| --gas                        | 200000                | [string] 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值        |          |
| --gas-adjustment             | 1                     | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略 ||
| --generate-only              |                       | 构建一个未签名交易并将其打印到标准输出                                 |          |
| --help, -h                   |                       | unjail命令帮助                                                      |          |
| --indent                     |                       | 在JSON响应中添加缩进                                                 |          |
| --json                       |                       | 以json形式输出                                                       |          |
| --ledger                     |                       | 使用连接的硬件记账设备                                                |          |
| --memo                       |                       | [string] 发送交易时的备忘录                                          |          |
| --node                       | tcp://localhost:26657 | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                   |          |
| --print-response             |                       | 返回交易响应 (当且仅当同步模式下使用)                                  |          |
| --sequence int               |                       | [int] 用来签名交易的sequence                                         |          |
| --shares-amount              |                       | [string] 指定恢复验证者身份是的股份                                    |         |
| --shares-percent             |                       | [string] 指定恢复验证者身份是的股份比，为0到1之间的正数                 |          |
| --trust-node                 | true                  | 关闭响应结果校验                                                      |          |

## 例子

### 恢复之前由于宕机被惩罚的验证者的身份

```shell
iriscli stake unjail --from=KeyName --fee=Fee --chain-id=ChainID
```

执行完成以后，你成功地恢复了指定的验证者的身份。
