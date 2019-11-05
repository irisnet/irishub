# iriscli bank

Bank模块用于管理本地帐户中的资产

## 可用命令

| 名称                                             | 描述                         |
| ------------------------------------------------ | ---------------------------- |
| [coin-type](#iriscli-bank-coin-type)             | 查询Coin的定义               |
| [token-stats](#iriscli-bank-token-stats)         | 查询通证统计信息             |
| [account](#iriscli-bank-account)                 | 查询账户余额                 |
| [send](#iriscli-bank-send)                       | 创建、签名、广播一个转账交易 |
| [burn](#iriscli-bank-burn)                       | 销毁通证                     |
| [set-memo-regexp](#iriscli-bank-set-memo-regexp) | 设置账户备注规则             |

## 常见问题

### ERROR: decoding bech32 failed

```bash
iriscli bank account iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f6d429zz
ERROR: decoding bech32 failed: checksum failed. Expected 46vaym, got d429zz.
```

这表示该帐户地址拼写错误，请仔细检查该地址。

### ERROR: account xxx does not exist

```bash
iriscli bank account iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr
ERROR: {"codespace":"sdk","code":9,"message":"account iaa1kenrwk5k4ng70e5s9zfsttxpnlesx5psh804vr does not exist"}
```

这通常是因为您要查询的帐户地址在链上没有任何交易。

## iriscli bank coin-type

在IRIShub上查询一种特殊的通证。 IRIShub上的本机通证是“ iris”，它具有以下可用单位：`iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` 和 `iris-atto`。

```bash
 iriscli bank coin-type <coin_name> <flags>
```

**标识：**

| 名称, 速记   | 类型   | 必须 | 默认                  | 描述                                  |
| ------------ | ------ | -------- | --------------------- | ------------------------------------- |
| -h, --help   |        |          |                       | `coin-type`子命令的提示信息           |
| --chain-id   | string |          |                       | tendermint节点的Chain ID              |
| --height     | int    |          |                       | 查询区块高度，忽略获取最新可证明区块  |
| --indent     | string |          |                       | 向JSON响应添加缩进                    |
| --ledger     | string |          |                       | 连接的Ledger设备                      |
| --node       | string |          | tcp://localhost:26657 | `<host>:<port>`，连接的tendermint节点 |
| --trust-node | string |          | true                  | 是否验证节点响应的结果                |

### 查询本地通证 `iris`

```bash
iriscli bank coin-type iris
```

之后，您将获得本地通证的详细信息 `iris`。

```bash
CoinType:
  Name:     iris
  MinUnit:  iris-atto: 18
  Units:    iris: 0,  iris-milli: 3,  iris-micro: 6,  iris-nano: 9,  iris-pico: 12,  iris-femto: 15,  iris-atto: 18
  Origin:   native
  Desc:     IRIS Network
```

## iriscli bank token-stats

查询通证统计信息，包括未质押通证总数，已销毁总数和已质押总数。

```bash
 iriscli bank token-stats <token-id> <flags>
```

**标识：**

| 名称, 速记   | 类型   | 必须 | 默认                  | 描述                                 |
| ------------ | ------ | -------- | --------------------- | ------------------------------------ |
| -h, --help   |        |          |                       | `token-stats`子命令的提示信息        |
| --chain-id   | string |          |                       | tendermint节点的Chain ID             |
| --height     | int    |          |                       | 查询区块高度，忽略获取最新可证明区块 |
| --indent     | string |          |                       | 向JSON响应添加缩进                   |
| --ledger     | string |          |                       | 连接的Ledger设备                     |
| --node       | string |          | tcp://localhost:26657 | `<host>:<port>` 连接的tendermint节点 |
| --trust-node | string |          | true                  | 是否验证节点响应的结果               |

### 查询通证统计信息

```bash
iriscli bank token-stats iris
```

输出:

```bash
TokenStats:
  Loose Tokens:             1404158512.076790096410637686iris
  Bonded Tokens:            609544925.59191727606475175iris
  Burned Tokens:            19205096.20000004iris
  Total Supply:             2013703437.668707372475389436iris
```

## iriscli bank account

该命令用于查询特定地址的余额信息。

```bash
iriscli bank account <address> <flags>
```

**标识：**

| 名称, 速记   | 类型   | 必须 | 默认                  | 描述                                 |
| ------------ | ------ | -------- | --------------------- | ------------------------------------ |
| -h, --help   |        |          |                       | `account`子命令的提示信息            |
| --chain-id   | string |          |                       | tendermint节点的Chain ID             |
| --height     | int    |          |                       | 查询区块高度，忽略获取最新可证明区块 |
| --ledger     | string |          |                       | 连接的Ledger设备                     |
| --node       | string |          | tcp://localhost:26657 | `<host>:<port>` 连接的tendermint节点 |
| --trust-node | string |          | true                  | 是否验证节点响应的结果               |

### 以信任模式查询您的帐户

```bash
iriscli bank account <address>
```

输出:

```bash
Account:
  Address:         iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym
  Pubkey:          iap1addwnpepqwnsrt9m8tevhy4fdqyarunzuzzgz8e5q8jlceyf7uwpw0q0ptp2cp3lmjt
  Coins:           50iris
  Account Number:  0
  Sequence:        2
  Memo Regexp:
```

## iriscli bank send

发送令牌到另一个地址，此命令包括 `generate`, `sign` 和 `broadcast` 这些步骤。

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=<amount> --fee=<native-fee> --chain-id=<chain-id>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                           |
| ---------- | ------ | -------- | ---- | ------------------------------ |
| --amount   | string | Yes      |      | 要发送的通证数量，例如：10iris |
| --to       | string |          |      | 接收通证的Bech32编码地址       |

### 将通证发送到另一个地址

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10iris --fee=0.3iris --chain-id=irishub
```

## iriscli bank burn

此命令用于从您自己的地址销毁通证。

```bash
iriscli bank burn --from=<key-name> --amount=<amount-to-burn> --fee=<native-fee> --chain-id=<chain-id>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                         |
| ---------- | ------ | -------- | ---- | ---------------------------- |
| --amount   | string | Yes      |      | 要销毁的通证数量，例如10iris |

### 销毁通证

```bash
 iriscli bank burn --from=<key-name> --amount=10iris --chain-id=irishub --fee=0.3iris
```

## iriscli bank set-memo-regexp

此命令用于为您自己的地址设置备注正则表达式，以便您只能从具有相应备注的交易中接收通证。

```bash
iriscli bank set-memo-regexp --regexp=<regular-expression> --from=<key-name> --fee=<native-fee> --chain-id=<chain-id>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                                          |
| ---------- | ------ | -------- | ---- | --------------------------------------------- |
| --regexp   | string | Yes      |      | 正则表达式，最大长度为50，例如^[A-ZA-Z0-9]+ $ |

### 为帐户地址设置备注正则表达式

```bash
iriscli bank set-memo-regexp --regexp=^[A-Za-z0-9]+$ --from=<key-name> --fee=0.3iris --chain-id=irishub
```
