# Random

Rand模块允许你向IRIS Hub发送随机数请求，查询随机数或待处理的随机数请求队列。

## 可用命令

| 名称                                             | 描述                               |
| ------------------------------------------------ | ---------------------------------- |
| [request-random](#iris-tx-random-request-random) | 请求一个随机数                     |
| [query-random](#iris-query-random-random)        | 使用ID查询链上生成的随机数         |
| [query-queue](#iris-query-random-queue)          | 查询随机数请求队列，支持可选的高度 |

## iris tx random request-random

请求一个随机数。

```bash
iris tx random request-random [flags]
```

**标志：**

| 名称，速记        | 类型   | 必须 | 默认  | 描述                                       |
| ----------------- | ------ | ---- | ----- | ------------------------------------------ |
| --block-interval  | uint64 |      | 10    | 请求的随机数将在指定的区块间隔后生成       |
| --oracle          | bool   |      | false | 是否使用 Oracle 方式                       |
| --service-fee-cap | string |      | ""    | 最大服务费用（如果使用 Oracle 方式则必填） |

### 请求一个随机数

向 IRIS Hub 发送随机数请求，该随机数将在`--block-interval`指定块数后生成。

```bash
# without oracle
iris tx random request-random --block-interval=100 --from=<key-name> --chain-id=irishub --fee=0.3iris --commit

# with oracle
iris tx random request-random --block-interval=100 --oracle=true --service-fee-cap=1iris --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

:::tip
如果交易已被执行，你将获得一个唯一的请求ID，该ID可用于查询请求状态。你也可以通过[查询交易详情](./tendermint.md#iriscli-tendermint-tx)获取请求ID。
:::

## iris query random random

使用ID查询链上生成的随机数。

```bash
iris query random random <request-id> [flags]
```

### 查询随机数

查询已生成的随机数。

```bash
iris query random random <request-id>
```

## iris query random queue

查询随机数请求队列，支持可选的高度。

```bash
iris query random queue <gen-height> [flags]
```

### 查询随机数请求队列

查询尚未处理的随机数请求，可指定将要生成随机数（或请求 Service）的区块高度。

```bash
iris query random queue 100000
```
