# iriscli bank coin-type

## 描述

查询IRISnet中某一个特殊通证. IRIShub的原始通证`iris`,其可用的单位如下： `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` 和 `iris-atto`. 

## 使用方式

```
 iriscli bank coin-type <coin_name> <flags>
``` 

## 标志

| 命令，速记   | 类型   | 是否必须 | 默认值                | 描述                                      |
| ------------ | ------ | -------- | --------------------- | ----------------------------------------- |
| -h, --help   |        | 否       |                       | 帮助                                      |
| --chain-id   | String | 否       |                       | tendermint 节点Chain ID                     |
| --height     | Int    | 否       |                       | 查询的区块高度用于获取最新的区块。        |
| --indent     | String | 否       |                       | 在JSON响应中增加缩进                      |
| --ledger     | String | 否       |                       | 使用ledger设备                    |
| --node       | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。 |
| --trust-node | String | 否       | True                  | 不验证响应的证明                          |

## 例子

### 查询本地通证`iris`

```
iriscli bank coin-type iris
```

执行命令后我们获取到的iris通证详细信息如下：

```
CoinType:
  Name:     iris
  MinUnit:  iris-atto: 18
  Units:    iris: 0,  iris-milli: 3,  iris-micro: 6,  iris-nano: 9,  iris-pico: 12,  iris-femto: 15,  iris-atto: 18
  Origin:   native
  Desc:     IRIS Network
```



## 扩展描述

查询IRISnet中一个指定通证的信息

​    



​           
