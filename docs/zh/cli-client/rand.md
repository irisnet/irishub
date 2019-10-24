# iriscli rand

Rand模块允许你向IRIS Hub发送随机数请求，查询随机数或待处理的随机数请求队列。

## 可用命令

| 名称                                       | 描述                               |
| ------------------------------------------ | ---------------------------------- |
| [request-rand](#iriscli-rand-request-rand) | 请求一个随机数                     |
| [query-rand](#iriscli-rand-query-rand)     | 使用ID查询链上生成的随机数         |
| [query-queue](#iriscli-rand-query-queue)   | 查询随机数请求队列，支持可选的高度 |

## iriscli rand request-rand

请求一个随机数。

```bash
iriscli rand request-rand <flags>
```

**标志：**

| 名称，速记       | 类型   | 必须 | 默认 | 描述                                 |
| ---------------- | ------ | ---- | ---- | ------------------------------------ |
| --block-interval | uint64 |      | 10   | 请求的随机数将在指定的区块间隔后生成 |

### 请求一个随机数

向 IRIS Hub 发送随机数请求，该随机数将在`--block-interval`指定块数后生成。

```bash
iriscli rand request-rand --block-interval=100 --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

:::tip
如果交易已被执行，你将获得一个唯一的请求ID，该ID可用于查询请求状态。你也可以通过[查询交易详情](./tendermint.md#iriscli-tendermint-tx)获取请求ID。
:::

## iriscli rand query-rand

使用ID查询链上生成的随机数。

```bash
iriscli rand query-rand <flags>
```

**标志：**

| 名称，速记   | 类型   | 必须 | 默认 | 描述                   |
| ------------ | ------ | ---- | ---- | ---------------------- |
| --request-id | string |      |      | 请求ID，由请求交易返回 |

## 查询随机数

查询已生成的随机数。

```bash
iriscli rand query-rand --request-id=035a8d4cf64fcd428b5c77b1ca85bfed172d3787be9bdf0887bbe8bbeec3932c
```

## iriscli rand query-queue

查询随机数请求队列，支持可选的高度。

```bash
iriscli rand query-queue <flags>
```

**标志：**

| 名称，速记     | 类型  | 必须 | 默认 | 描述           |
| -------------- | ----- | ---- | ---- | -------------- |
| --queue-height | int64 |      | 0    | 查询的目标高度 |

## Query random number request queue

查询尚未处理的随机数请求，可指定将要生成随机数的区块高度。

```bash
iriscli rand query-queue --queue-height=100000
```
