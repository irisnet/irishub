---
order: 1
---

# 简介

`iris` 是IRIShub网络的命令行客户端。IRIShub用户可以使用`iris`发送交易和查询区块链数据。

## 工作目录

`iris` 默认工作目录是 `$HOME/.iris`，主要用于保存配置文件和数据。IRIShub “密钥”数据保存在`iris`的工作目录中。您还可以通过`--home`指定`iris`的工作目录。

## 连接全节点

`--node`用来指定所连接`iris`节点的rpc地址，交易和查询的消息都发送到监听这个端口的iris进程。默认值为`tcp://localhost:26657`。

## 全局标识

### GET 请求

所有GET请求都有以下全局标识:

| 名称，速记 | 类型   | 必需 | 默认值               | 描述                     |
| ---------- | ------ | ---- | -------------------- | ------------------------ |
| --chain-id | string |      | ""                   | tendermint节点的Chain ID |
| --home     | string |      | /Users/bianjie/.iris | 配置和数据的目录         |
| --trace    | string |      |                      | 打印出错时的完整堆栈跟踪 |

### POST 请求

所有POST请求都有以下全局标识:

| 名称，速记        | 类型   | 必需 | 默认值                | 描述                                                                                                                |
| ----------------- | ------ | ---- | --------------------- | ------------------------------------------------------------------------------------------------------------------- |
| --account-number  | int    |      | 0                     | 发起交易的账户的编号                                                                                                |
| --broadcast-mode  | string |      | sync                  | 广播交易的节点                                                                                                      |
| --dry-run         | bool   |      | false                 | 模拟执行交易，并返回消耗的`gas`。`--gas`指定的值会被忽略                                                            |
| --fees            | string |      |                       | 交易费（指定交易费的上限）                                                                                          |
| --from            | string |      |                       | 签名交易的私钥名称                                                                                                  |
| --gas             | string |      | 50000                 | 交易的gas上限；设置为"simulate"将自动计算相应的阈值                                                                 |
| --gas-adjustment  | float  |      | 1.5                   | gas调整因子，这个值降乘以模拟执行消耗的`gas`，计算的结果返回给用户；如果`--gas`的值不是`simulate`，这个标志将被忽略 |
| --gas-prices      | string |      |                       | gas以十进制格式确定交易费                                                                                           |
| --generate-only   | bool   |      | false                 | 是否仅仅构建一个未签名的交易便返回                                                                                  |
| --help, -h        | string |      |                       | 打印帮助信息                                                                                                        |
| --keyring-backend | string |      | os                    | 选择钥匙圈后端                                                                                                      |
| --ledger          | bool   |      | false                 | 使用ledger设备                                                                                                      |
| --memo            | string |      |                       | 指定交易的memo字段                                                                                                  |
| --node            | string |      | tcp://localhost:26657 | tendermint节点的rpc地址                                                                                             |
| --offline         | string |      |                       | 离线节点                                                                                                            |
| --sequence        | int    |      | 0                     | 发起交易的账户的sequence                                                                                            |
| --sign-mode       | string |      |                       | 选择签名节点，这是高级特性                                                                                          |
| --trust-node      | bool   |      | true                  | 是否信任全节点返回的数据，如果不信任，客户端会验证查询结果的正确性                                                  |
| --yes             | bool   |      | true                  | 跳过交易广播提示确认                                                                                                |
| --chain-id        | string |      |                       | tendermint节点的`Chain ID`                                                                                          |
| --home            | string |      |                       | 配置文件和数据文件目录 (默认 "~/.iris")                                                                             |
| --trace           | string |      |                       | 打印出完整的堆栈跟踪错误                                                                                            |

## 模块命令列表

| **子命令**                        | **描述**                           |
| --------------------------------- | ---------------------------------- |
| [bank](./bank.md)                 | 用于查询帐户和转账等的 Bank 子命令 |
| [debug](./debug.md)               | 调试子命令                         |
| [distribution](./distribution.md) | 收益管理子命令                     |
| [gov](./gov.md)                   | 治理和投票子命令                   |
| [htlc](./htlc.md)                 | HTLC 子命令                        |
| [keys](./keys.md)                 | 密钥管理子命令                     |
| [nft](./nft.md)                   | NFT 子命令                         |
| [oracle](./oracle.md)             | Oracle 子命令                      |
| [params](./params.md)             | 查询治理参数子命令                 |
| [random](./rand.md)               | 随机数子命令                       |
| [record](./record.md)             | Record 子命令                      |
| [slashing](./slashing.md)         | Slashing 子命令                    |
| [service](./service.md)           | 服务子命令                         |
| [staking](./staking.md)           | Staking 子命令                     |
| [status](./status.md)             | 查询远程节点的状态                 |
| [tendermint](./tendermint.md)     | Tendermint 状态查询子命令          |
| [token](./token.md)               | 资产子命令                         |
| [tx](./tx.md)                     | Tx 子命令                          |
| [upgrade](./upgrade.md)           | 软件升级子命令                     |
