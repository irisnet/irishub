# iriscli tendermint validator-set

## 描述
查询指定高度的验证人信息，包括地址公钥，投票权重以及当前出块优先级。如果未指定高度，则查询最新高度的验证人信息。

## 用法

```
  iriscli tendermint validator-set <height> <flags>
```
或者
```
  iriscli tendermint validator-set
```

## 标志

| 名称, 速记 | 默认值                    | 描述                                                             | 必需      |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    | 无 | tendermint节点的Chain ID   | 是       |
| --node string     |   tcp://localhost:26657                         | 要连接的节点  |                                     
| --help, -h      |           无| 	下载命令帮助|
| --trust-node    | true                       | 信任连接的完整节点，关闭响应结果校验                                            |          |

## 例子 
### 查询114360高度上区块中的validator-set

```shell
iriscli tendermint validator-set 114360 --chain-id=<chain-id>
```

### 查询最新区块中的validator-set

```shell
 iriscli tendermint validator-set --chain-id=<chain-id> --trust-node
```

示例结果：
```json
{
  "block_height": "114360",
  "validators": [
    {
      "address": "ica1q9zpqvm7cadx5walcg5jkdxklayr8c2uqtmzmx",
      "pub_key": "icp1zcjduepq8fnuxnceuy4n0fzfc6rvf0spx56waw67lqkrhxwsxgnf8zgk0nus66rkg4",
      "proposer_priority": "-300",
      "voting_power": "100"
    },
    {
      "address": "ica1qxavppd679lyxxu9fdu0zxxfv59r7e0wfgapj7",
      "pub_key": "icp1zcjduepquvkj9qa9mgyhudkhsqxelr0k4zf45ehw4sv4m5wktzhke4zvskas5rk9vq",
      "proposer_priority": "100",
      "voting_power": "100"
    },
    {
      "address": "ica1grd8wp7vezr4czen2nujpejvt6597fmrkqs7h0",
      "pub_key": "icp1zcjduepqnudzfngr6aq4hk47w6p9jx5w97fxmwj2vwwvpkd3sez3dzrm359szchwyl",
      "proposer_priority": "100",
      "voting_power": "100"
    },
    {
      "address": "ica15rg635p4j3xpxcs53dwl6nl2u7gjjsvs7m4psw",
      "pub_key": "icp1zcjduepqxhc5c0fyfwta05tax036jmrr2x6aea2smnce9zhmravt7gwpm0qqjhn9nz",
      "proposer_priority": "100",
      "voting_power": "100"
    }
  ]
}
```
