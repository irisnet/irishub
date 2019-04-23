# iriscli bank broadcast

## 描述

使用此命令在网络中广播已在线下[签名](./sign.md)的交易

## 使用方式

```
iriscli bank broadcast <tx-file> <flags> 
```

## 标志

| 命令，速记   | 类型   | 是否必须 | 默认值                | 描述                                      |
| ------------ | ------ | -------- | --------------------- | ----------------------------------------- |
| -h, --help   |        | 否       |                       | 打印帮助信息                              |
| --chain-id   | String | 否       |                       | tendermint 节点网络ID                     |
| --height     | Int    | 否       |                       | 查询的区块高度用于获取最新的区块。        |
| --ledger     | String | 否       |                       | 使用一个联网的分账设备                    |
| --node       | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。 |
| --trust-node | String | 否       | True                  | 不验证响应的证明                          |


## 全局标志

| 命令，速记            | 默认值         | 描述                                | 是否必须 |
| --------------------- | -------------- | ----------------------------------- | -------- |
| -e, --encoding string | hex            | 字符串二进制编码 (hex\b64\btc ) | 否       |
| --home string         | $HOME/.iriscli | 配置和数据存储目录                  | 否       |
| -o, --output string   | text           | 输出格式 (text\json)              | 否       |
| --trace               |                | 出错时打印完整栈信息                | 否       |


## 例子

### 广播交易

```
iriscli bank broadcast sign.json --chain-id=<chain-id>
```
