# Farm

Farm 模块允许您在`irishub`上轻松创建`Farm`活动。

## 可用命令

| 名称                                       | 描述                         |
| ------------------------------------------ | ---------------------------- |
| [创建 Pool](#iris-tx-farm-create)          | 创建一个新的`Farm`池         |
| [调整 Pool 参数](#iris-tx-farm-adjust)     | 调整`Farm`池参数             |
| [销毁 Pool](#iris-tx-farm-destroy)         | 销毁`Farm`池并取回投入的奖金 |
| [抵押流动性](#iris-tx-farm-stake)          | 抵押流动性代币               |
| [取回奖励](#iris-tx-farm-harvest)          | 取回参与`Farm`池而获取的奖金 |
| [查询 Farmer](#iris-query-farm-farmer)     | 查询 Farmer 信息             |
| [查询`Farm`池](#iris-query-farm-pool)      | 查询`Farm`池的当前状态       |
| [分页查询`Farm`池](#iris-query-farm-pools) | 按页查询`Farm`池信息         |
| [查询可治理参数](#iris-query-farm-params)  | 查询`Farm`模块的可治理参数   |

## iris tx farm create

创建一个新的`Farm`池 并支付手续费和奖金。

```bash
iris tx farm create <Farm Pool Name> [flags]
```

**标志：**

| 名称，速记         | 是否必须 | 默认  | 描述                      |
| ------------------ | -------- | ----- | ------------------------- |
| --lp-token-denom   | 是       |       |`Farm`池可接受的流动性代币 |
| --reward-per-block | 是       |       | 每个区块的奖励            |
| --total-reward     | 是       |       | 总奖励                    |
| --description      | 否       | ""    |`Farm`池的简要描述         |
| --start-height     | 是       |       |`Farm`池的开始高度         |
| --editable         | 否       | false |`Farm`池是否可编辑         |

### iris tx farm adjust

在`Farm`池结束前调整池的参数，例如`reward-per-block`、`total-reward`。

```bash
iris tx farm adjust <Farm Pool Name> [flags]
```

**标志：**

| 名称，速记          | 是否必须                      | 默认 | 描述                 |
| ------------------- | ----------------------------- | ---- | -------------------- |
| --additional-reward | 和`--reward-per-block`二选一  | ""   | 向`Farm`池追加的奖金 |
| --reward-per-block  | 和`--additional-reward`二选一 | ""   | 每个区块的奖励       |

## iris tx farm destroy

销毁`Farm`池，取回投入的奖金。此时用户农场获得的奖励结束，需要用户手动取回奖励和流动性代币。

```bash
iris tx farm destroy <Farm Pool Name> [flags]
```

### iris tx farm stake

Farmer 通过抵押`Farm`池指定的流动性代币来参与`Farm`活动。参与活动获得的奖励与质押代币数量和`Farm`池参数有关。

```bash
iris tx farm stake <Farm Pool Name> <lp-token> [flags]
```

### iris tx farm harvest

Farmer 取回奖励。

```bash
iris tx farm harvest <Farm Pool Name> [flags]
```

### iris query farm farmer

查询 Farmer 信息，包括待领取的奖励、抵押的流动性等。

```bash
iris query farm farmer <Farmer Address> --pool-name <Farm Pool Name>
```

**标志：**

| 名称，简写  | 必填 | 默认 | 说明        |
| ----------- | ---- | ---- | ----------- |
| --pool-name | 否   | ""   |`Farm`池名称 |

### iris query farm pool

按名称查询`Farm`池的相关信息

```bash
iris query farm pool <Farm Pool Name>
```

### iris query farm pools

分页查询`Farm`池

```bash
iris query farm pools
```

### iris query farm params

查询`Farm`模块的可治理的参数

```bash
iris query farm params
```
