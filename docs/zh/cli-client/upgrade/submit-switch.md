# iriscli upgrade submit-switch

## 描述

安装完新软件后，向这次升级相关的提议发送switch消息，表明自己已经安装新软件并把消息广播到全网。

## 用例

```
iriscli upgrade submit-switch [flags]
```

## 标志

| 名称, 速记       | 默认值    | 描述                                                         | 必需     |
| ---------------  | --------- | ------------------------------------------------------------ | -------- |
| --proposal-id    |           | 软件升级提议的ID                                             | 是       |
| --title          |           | switch消息对标题                                             |          |
| --account-number |           | [int] 用来签名交易的AccountNumber                            |          |
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
| --print-response |                            | 返回交易响应 (当且仅当同步模式下使用))                                                 |          |
| --sequence       |                            | [int] 用来签名交易的sequence number                                                  |          |
| --trust-node     | true                       | 关闭响应结果校验                                                                     |          |
## 用例

发送对软件升级提议（ID为5）switch消息

```
iriscli upgrade submit-switch --chain-id=IRISnet --from=x --fee=0.004iris --proposalID 5 --title="Run new verison"
```