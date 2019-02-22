# iriscli tendermint validator-set

## 描述
根据指定高度在验证器上查询

## 用法

```
  iriscli tendermint validator-set [height] [flags]

```

## 标志

| 名称, 速记 | 默认值                    | 描述                                                             | 必需      |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    | 无 | [string] tendermint节点的链ID   | 是       |
| --node string     |   tcp://localhost:26657                         | 要连接的节点  |                                     
| --help, -h      |           无| 	下载命令帮助|
| --trust-node    | true                       | 信任连接的完整节点，关闭响应结果校验                                            |          |

## 例子 
### 查询某高度上区块中的validator-set


```shell
iriscli tendermint validator-set 114360 --chain-id=irishub-test
```
之后你会在验证器上查询到该高度的信息
### 查询最新区块中的validator-set

```shell
 iriscli tendermint validator-set --chain-id=irishub-test --trust-node=true

```

示例结果：

```json
{
  "block_height": "113",
  "validators": [
    {
      "address": "fca1q9zpqvm7cadx5walcg5jkdxklayr8c2ucya6mm",
      "pub_key": "fcp1zcjduepq8fnuxnceuy4n0fzfc6rvf0spx56waw67lqkrhxwsxgnf8zgk0nus2r55he",
      "proposer_priority": "-300",
      "voting_power": "100"
    },
    {
      "address": "fca1qxavppd679lyxxu9fdu0zxxfv59r7e0w38mejr",
      "pub_key": "fcp1zcjduepquvkj9qa9mgyhudkhsqxelr0k4zf45ehw4sv4m5wktzhke4zvskasy6p8nv",
      "proposer_priority": "100",
      "voting_power": "100"
    },
    {
      "address": "fca1grd8wp7vezr4czen2nujpejvt6597fmrw0kxhj",
      "pub_key": "fcp1zcjduepqnudzfngr6aq4hk47w6p9jx5w97fxmwj2vwwvpkd3sez3dzrm359sjpqvmn",
      "proposer_priority": "100",
      "voting_power": "100"
    },
    {
      "address": "fca15rg635p4j3xpxcs53dwl6nl2u7gjjsvsx5nesn",
      "pub_key": "fcp1zcjduepqxhc5c0fyfwta05tax036jmrr2x6aea2smnce9zhmravt7gwpm0qqzwy8vw",
      "proposer_priority": "100",
      "voting_power": "100"
    }
  ]
}
```
