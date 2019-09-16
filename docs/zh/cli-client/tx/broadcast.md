# iriscli tx broadcast

## 描述

使用此命令在网络中广播已在线下[签名](./sign.md)的交易

## 使用方式

```
iriscli tx broadcast <tx-file> <flags> 
```

## 标志

| 命令，速记     | 类型   | 是否必须  | 默认值                  | 描述                                   |
| ------------ | ------ | -------- | --------------------- | ------------------------------------- |
| -h, --help   |        |          |                       | 打印帮助信息                             |
| --chain-id   | string |          |                       | tendermint 节点网络ID                   |
| --height     | int    |          |                       | 查询的区块高度用于获取最新的区块            |
| --ledger     | string |          |                       | 使用ledger设备                          |
| --node       | string |          | tcp://localhost:26657 | `<主机>:<端口>` 链上的tendermint rpc 接口 |
| --trust-node | string |          | true                  | 不验证响应的证明                          |

## 例子

### 广播交易

```
iriscli tx broadcast sign.json --chain-id=<chain-id>
```
