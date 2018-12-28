# iriscli bank coin-type

## 描述

查询IRIS网络某一个特殊通证. IRIShub的原始通证 `iris`,其可用的单位如下： `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` 和 `iris-atto`. 

## 使用方式

```
 iriscli bank coin-type [coin_name] [flags]
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
iriscli bank coin-type iris
```

执行命令后我们获取到的iris通证详细信息如下：

```
[root@ce7da33d46c3 output]# iriscli bank coin-type iris
{
  "name": "iris",
  "min_unit": {
    "denom": "iris-atto",
    "decimal": "18"
  },
  "units": [
    {
      "denom": "iris",
      "decimal": "0"
    },
    {
      "denom": "iris-milli",
      "decimal": "3"
    },
    {
      "denom": "iris-micro",
      "decimal": "6"
    },
    {
      "denom": "iris-nano",
      "decimal": "9"
    },
    {
      "denom": "iris-pico",
      "decimal": "12"
    },
    {
      "denom": "iris-femto",
      "decimal": "15"
    },
    {
      "denom": "iris-atto",
      "decimal": "18"
    }
  ],
  "origin": 1,
  "desc": "IRIS Network"
}
```



## 扩展描述

查询iris网络中一个指定通证的信息

​    



​           
