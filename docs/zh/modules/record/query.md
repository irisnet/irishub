# iriscli record query

## 描述

查询存证信息

## 用法

```
iriscli record query [flags]
```

## 标志

| 名称, 速记       | 默认值                     | 描述                                                        | 必需     |
| --------------- | -------------------------- | ---------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] tendermint节点的链ID                               | 是       |
| --height        | 最近可证明区块高度           | [int] 查询的区块高度                                              |          |
| --help, -h      |                            | 查询命令帮助                                                |          |
| --indent        |                            | 在JSON格式的应答中添加缩进                                   |          |
| --ledger        |                            | 使用连接的硬件记账设备                                       |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口> |          |
| --record-id     |                            | [string] 存证ID                                            | 是       |
| --trust-node    | true                       | 关闭响应结果校验                                            |          |

## 例子

### 查询存证

```shell
iriscli record query --chain-id=test --record-id=MyRecordID
```

在这之后，你将得到存证ID指定的存证的详细信息。

```json
{
  "submit_time": "2018-11-13 15:31:36",
  "owner_addr": "faa122uzzpugtrzs09nf3uh8xfjaza59xvf9rvthdl",
  "record_id": "record:ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6",
  "description": "description",
  "data_hash": "ab5602bac13f11737e8798dd57869c468194efad2db37625795f1efd8d9d63c6",
  "data_size": "24",
  "data": "this is my on chain data"
}
```
