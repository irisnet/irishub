# iriscli stake redelegations

## 描述

基于委托者地址的所有重新委托记录查询

## 用法

```
iriscli stake redelegations [delegator-addr] [flags]
```

## 标志

| 名称, 速记                  | 默认值                      | 描述                                                               | 必需     |
| -------------------------- | -------------------------- | ------------------------------------------------------------------- | -------- |
| --chain-id                 |                            | [string] Tendermint节点的链ID                                        |          |
| --height                   | 最新的可证明区块高度         | 查询的区块高度                                                       |          |
| --help, -h                 |                            | redelegations命令帮助                                                |          |
| --indent                   |                            | 在JSON响应中添加缩进                                                  |          |
| --ledger                   |                            | 使用连接的硬件记账设备                                                |          |
| --node                     | tcp://localhost:26657      | [string] Tendermint远程过程调用的接口\<主机>:\<端口>                   |          |
| --trust-node               | true                       | 关闭响应结果校验                                                      |          |

## 例子

### 基于委托者地址的所有重新委托记录查询

```shell
iriscli stake redelegations DelegatorAddress
```

运行成功以后，返回的结果如下：

```json
[
  {
    "delegator_addr": "faa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9sytjfk",
    "validator_src_addr": "fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2",
    "validator_dst_addr": "fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd",
    "creation_height": "1130",
    "min_time": "2018-11-16T07:22:48.740311064Z",
    "initial_balance": "0.1iris",
    "balance": "0.1iris",
    "shares_src": "0.1000000000",
    "shares_dst": "0.1000000000"
  }
]
```
