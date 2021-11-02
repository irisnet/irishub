# TIBC

`TIBC` todo
## Client

### Available Commands

| Name                                     | Description                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [create](#iris-tx-gov-submit-proposal-client-create)              | 使用指定的客户端状态和共识状态创建一个新的 TIBC 客户端。            |
| [update](#iris-tx-tibc-update)                |  使用header更新现有的 TIBC 客户端。                 |
| [upgrade](#iris-tx-gov-submit-proposal-client-upgrade)                | 升级具有指定客户端状态和共识状态的 TIBC 客户端。                                      |
| [regesiter](#iris-tx-gov-submit-proposal-relayer-register)        | 提交指定客户端的 relayer 注册提案。                                 |
| [header](#iris-query-tibc-client-header)        | 查询正在运行的链的最新 Tendermint 头。                    |
| [node-state](#iris-query-tibc-client-node-state)        | 查询节点共识状态。 此结果将提供给客户端创建交易。                                 |
| [relayers](#iris-query-tibc-client-relayers)        | 查询一个客户端所有注册的relayer。                               |
| [client-state](#iris-query-tibc-client-state)        | 查询客户端状态。                              |
| [client-states](#iris-query-tibc-client-states)        | 查询所有可用的轻客户端。                              |
| [consensus-state](#iris-query-consensus-state)        | 查询给定高度的特定轻客户端的共识状态。                                 |
| [consensus-states](#iris-query-consensus-states)        | 查询一个客户端的所有共识状态。                                |


### iris tx gov submit-proposal client-create
提交客户端创建提案。

```bash 
iris tx gov submit-proposal client-create [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]
```

**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          提议的标题  |
| --description        |          |  |          提议的描述   |
| --deposit        |          |  |          提议的押金  |

### iris tx tibc update

使用header更新现有客户端。

```bash
iris tx tibc update [chain-name] [path/to/header.json]
```

## iris tx gov submit-proposal client-upgrade

Submit a client upgrade proposal.

```bash
iris tx gov submit-proposal client-upgrade [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]
```
**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          提议的标题  |
| --description        |          |  |          提议的描述   |
| --deposit        |          |  |          提议的押金  |


## iris tx gov submit-proposal relayer-register

提交注册relayer提案

```bash
iris tx gov submit-proposal relayer-register [chain-name] [relayers-address] [flags]
```

**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          提议的标题  |
| --description        |          |  |          提议的描述   |
| --deposit        |          |  |          提议的押金  |


### iris query consensus state

查询给定高度的特定轻客户端的共识状态。

```bash
iris query tibc client consensus-state [chain-name] [{revision}-{height}] [flags]
```
**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | --------------------------- | ----------- |
| --height          |   int       |  |       使用特定高度查询状态（如果节点处于修剪状态，这可能会出错）      |
| --latest-height          |          |             |  返回最新存储的共识状态，格式：{revision}-{height}           |
| --node          |   string       |       tcp://localhost:26657      | Host:port 到此链的 Tendermint RPC 接口          |
| --prove          |          |         true    |  显示查询结果的证明       |

### iris query consensus states

查询一个客户端的所有共识状态。
```bash
iris query tibc client consensus-states [chain-name] [flags]
```

**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | --------------------------- | ----------- |
| --count-total          |          |  |       计算要查询的共识状态中的记录总数      |
| --height          |          |             |  返回最新存储的共识状态，格式：{revision}-{height}           |

### iris query tibc client header
查询正在运行的链的最新 Tendermint 头。
```bash
iris query tibc client header [flags]
```
| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | --------------------------- | ----------- |
| --node          |   string       |       tcp://localhost:26657      | Host:port 到此链的 Tendermint RPC 接口          |
| --height          |          |             |   返回最新存储的共识状态，格式：{revision}-{height}           |


### iris query tibc client node-state
查询节点共识状态。 此结果将提供给客户端创建交易。
```bash
iris query tibc client node-state [flags]
```

### iris query tibc client relayers
查询一个客户端所有注册的relayer。
```bash
iris query tibc client relayers [chain-name] [flags]
```

### iris query tibc client state
查询客户端状态。
```bash
iris query tibc client state [chain-name] [flags]
```

### iris query tibc client states 
查询所有可用的轻客户端。
```bash
iris query tibc client states [flags]
```

## Packet

### Available Commands

| Name                                     | Description                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [send-clean-packet](#iris-tx-tibc-packet-send-clean-packet)              | 发送 clean 数据包。            |
| [clean-packet-commitment](#iris-query-tibc-packet-clean-packet-commitment)                |  查询 clean 数据包承诺。                 |
| [packet-ack](#iris-query-tibc-packet-packet-ack)                | 查询数据包确认。                                      |
| [packet-commitment](#iris-query-tibc-packet-packet-commitment)        | 查询数据包承诺。                                                               |
| [packet-commitments](#iris-query-tibc-packet-packet-commitments)                | 查询与源关联的所有数据包承诺。                                      |
| [packet-receipt](#iris-query-tibc-packet-packet-receipt)        | 查询数据包收据。                                                             |
| [unreceived-acks](#iris-query-tibc-packet-unreceived-acks)        | 查询所有与源链名和目的链名关联的未接收数据包。                                                          |


### iris tx tibc packet send-clean-packet
发送 clean 数据包。
```bash
iris tx tibc packet send-clean-packet [dest-chain-name] [sequence] [flags]
```

**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | --------------------------- | ----------- |
| --relay-chain-name          |   string       |            | 中继链的名字         |


### iris query tibc packet clean-packet-commitment
查询clean 数据包承诺。
```bash
iris tx tibc packet send-clean-packet [dest-chain-name] [sequence] [flags]
```

### iris query tibc packet packet-ack
查询数据包确认。
```bash
iris query tibc packet packet-ack [source-chain] [dest-chain] [sequence] [flags]
```
### iris query tibc packet packet-commitment
查询数据包承诺。
```bash
iris query tibc packet packet-commitment [source-chain] [dest-chain] [sequence] [flags]
```

### iris query tibc packet packet-commitments
查询与源关联的所有数据包承诺。
```bash
iris query tibc packet packet-commitments [source-chain] [dest-chain] [flags]
```

### iris query tibc packet packet-receipt
查询数据包收据。
```bash
iris query tibc packet packet-receipt [source-chain] [dest-chain] [sequence] [flags]
```

### iris query tibc packet unreceived-acks
查询与源链名称和目标链名称相关联的所有未接收的 ack。
```bash
iris query tibc packet unreceived-acks [source-chain] [dest-chain]] [flags]
```

**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | --------------------------- | ----------- |
| --sequences          |  int64Slice        |            | 包序列号的逗号分隔列表（默认 []）         |


### iris query tibc packet unreceived-packets
查询所有与源链名和目的链名关联的未接收数据包。
```bash
iris query tibc packet unreceived-packets [source-chain] [dest-chain] [flags]
```

**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | --------------------------- | ----------- |
| --sequences          |  int64Slice        |            | 包序列号的逗号分隔列表（默认 []）          |


## Routing

### Available Commands

| 名称                                     | 描述                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [set-rules](#iris-tx-gov-submit-proposal-set-rules)              | 提交设置规则提案。          |
| [routing-rules](#iris-query-tibc-routing-routing-rules)                |  查询路由规则承诺。                 |

### iris tx gov submit-proposal set-rules
提交设置规则提案。

```bash
iris tx gov submit-proposal set-rules [path/to/routing_rules.json] [flags]
```
**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | ------------------------------- | ----------- |
| --title          |          |            |          提议的标题  |
| --description        |          |  |          提议的描述   |
| --deposit        |          |  |          提议的押金  |

### iris query tibc routing routing-rules
查询路由规则承诺。
```bash
iris query tibc routing routing-rules
```

## Nft-Transfer

### Available Commands

| Name                                     | Description                                                                            |
| ---------------------------------------- | -------------------------------------------------------------------------------------- |
| [nft-transfer](#iris-tx-tibc-nft-transfer-transfer)              | 通过 TIBC 转移NFT。            |
| [class-trace](#iris-query-tibc-nft-transfer-class-trace)                |  Query the class trace info from a given trace hash.                |
| [class-traces](#iris-query-tibc-nft-transfer-class-traces)                |  Query the trace info for all nft classes.                |


### iris tx tibc-nft-transfer transfer
通过 TIBC 转移NFT。
```bash
iris tx tibc-nft-transfer transfer [dest-chain] [receiver] [class] [id] [flags]
```
**Flags:**

| 名称，速记  | 必须 | 默认                         | 描述 |
| --------------- | -------- | --------------------------- | ----------- |
| --relay-chain          |  string        |            | 跨链 NFT 使用的中继链         |
| --dest-contract          |  string        |            | 接收 nft 的目标合约地址         |


### iris query tibc-nft-transfer class-trace
从给定的 trace 哈希查询class trace 信息。
```bash
iris query tibc-nft-transfer class-trace [hash]
```

### iris query tibc-nft-transfer class-traces
查询所有 nft 类的 trace 信息。
```bash
iris query tibc-nft-transfer class-traces
```