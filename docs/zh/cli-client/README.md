# 命令行客户端

## 查询命令的flags

| 名称, 速记       | 类型         |必需          |默认值                | 描述                                                        | 
| --------------- | ----   | -------- | --------------------- | -------------------------------------------------------------------- |
| --chain-id      | string | false    | ""                    | Tendermint节点的Chain ID |
| --height        | int    | false    | 0                     | 查询某个高度的区块链数据，如果是0，这返回最新的区块链数据 |
| --help, -h      | string | false    |                       | 打印帮助信息 |
| --indent        | bool   | false    | false                 | 格式化json字符串|
| --ledger        | bool   | false    | false                 | 是否使用硬件钱包 |
| --node          | string | false    | tcp://localhost:26657 | tendermint节点的rpc地址|
| --trust-node    | bool   | false    | true                  | 是否信任全节点返回的数据，如果不信任，客户端会验证查询结果的正确性 |

每个区块链状态查询命令都包含上表中的flags，同时不同查询命令还可能会有自己独有的flags。

## 发送交易命令的flags


| 名称, 速记        | 类型         |必需          |默认值                | 描述                         |
| -----------------| -----  | -------- | --------------------- | ------------------------------------------------------------------- |
| --account-number | int    | false    | 0                     | 发起交易的账户的编号 |
| --async          | bool   | false    | false                 | 是否异步广播交易（仅当commit为false时有效） |
| --commit         | bool   | false    | false                 | 广播交易并等到交易被打包再返回 |
| --chain-id       | string | true     | ""                    | Tendermint节点的`Chain ID` |
| --dry-run        | bool   | false    | false                 | 模拟执行交易，并返回消耗的`gas`。`--gas`指定的值会被忽略 |
| --fee            | string | true     | ""                    | 交易费 |
| --from           | string | false    | ""                    | 发送交易的账户名称 |
| --from-addr      | string | false    | ""                    | 签名地址，在`generate-only`为`true`的情况下有效 |
| --gas            | int    | false    | 200000                | 交易的gas上限; 设置为"simulate"将自动计算相应的阈值 |
| --gas-adjustment | int    | false    | 1                     | gas调整因子，这个值降乘以模拟执行消耗的`gas`，计算的结果返回给用户; 如果`--gas`的值不是`simulate`，这个标志将被忽略 |
| --generate-only  | bool   | false    | false                 | 是否仅仅构建一个未签名的交易便返回 |
| --help, -h       | string | false    |                       | 打印帮助信息 |
| --indent         | bool   | false    | false                 | 格式化json字符串 |
| --json           | string | false    | false                 | 指定返回结果的格式，`json`或者`text` |
| --ledger         | bool   | false    | false                 | 是否使用硬件钱包|
| --memo           | string | false    | ""                    | 指定交易的memo字段 |
| --node           | string | false    | tcp://localhost:26657 | tendermint节点的rpc地址 |
| --print-response | bool   | false    | false                 | 是否打印交易返回结果，仅在`async`为true的情况下有效|
| --sequence int   | int    | false    | 0                     | 发起交易的账户的sequence |
| --trust-node     | bool   | false    | true                  | 是否信任全节点返回的数据，如果不信任，客户端会验证查询结果的正确性 | 

每个发送交易的命令都包含上表中的flags，同时不同交易的命令还可能会有自己独有的flags。

## 模块列表

1. [bank command](./bank/README.md)
2. [distribution command](./distribution/README.md)
3. [gov command](./gov/README.md)
4. [keys command](./keys/README.md)
5. [lcd command](./lcd/README.md)
6. [record command](./record/README.md)
7. [service command](./service/README.md)
8. [stake command](./stake/README.md)
9. [status command](./status/README.md)
10. [tendermint command](./tendermint/README.md)
11. [upgrade command](./upgrade/README.md)

## iriscli config命令

iriscli config命令可以交互式地配置一些默认参数，例如chain-id，home，fee 和 node。
