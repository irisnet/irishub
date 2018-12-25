# iriscli bank token-stats

## 描述

查询token统计数据，包括全部流通的token，全部绑定的token，以及销毁的token

## 使用方式

```
 iriscli bank token-stats [flags]
```

 

## 标志

| 命令，速记   | 类型   | 是否必须 | 默认值                | 描述                                      |
| ------------ | ------ | -------- | --------------------- | ----------------------------------------- |
| -h, --help   |        | 否       |                       | 帮助                                      |
| --chain-id   | String | 否       |                       | tendermint 节点网络ID                     |
| --height     | Int    | 否       |                       | 查询的区块高度用于获取最新的区块。        |
| --indent     | String | 否       |                       | 在JSON响应中增加缩进                      |
| --ledger     | String | 否       |                       | 使用一个联网的分账设备                    |
| --node       | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。 |
| --trust-node | String | 否       | True                  | 不验证响应的证明                          |



## 全局标志

| 命令，速记            | 默认值         | 描述                                | 是否必须 |
| --------------------- | -------------- | ----------------------------------- | -------- |
| -e, --encoding string | hex            | 字符串二进制编码 (hex \|b64 \|btc ) | 否       |
| --home string         | /root/.iriscli | 配置和数据存储目录                  | 否       |
| -o, --output string   | text           | 输出格式 (text \|json)              | 否       |
| --trace               |                | 出错时打印完整栈信息                | 否       |

## 

## 例子

### 查询本地通证iris

```
iriscli bank token-stats
```

示例返回：

```json
{
  "loosen_token": [
    "1864477.596384156921391687iris"
  ],
  "burned_token": [
    "177.59638iris"
  ],
  "bonded_token": "425182.329615843078608313iris"
}
```
