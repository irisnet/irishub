# Oracle

Oracle模块负责管理你在IRIS Hub上创建的feed。

## 可用命令

| 名称                                       | 描述                                 |
| ------------------------------------------ | ------------------------------------ |
| [create](#iris-tx-oracle-create)           | 创建一个新的feed，初始状态为"paused" |
| [start](#iris-tx-oracle-start)             | 启动一个处于"paused"状态的feed       |
| [pause](#iris-tx-oracle-pause)             | 暂停一个处于"running"状态的feed      |
| [edit](#iris-tx-oracle-edit)               | feed的所有者编辑一个feed的相关信息并更新服务调用参数   |
| [feed](#iris-query-oracle-feed)            | 通过名称查询一个feed信息             |
| [feeds](#iris-query-oracle-feeds)          | 查询一组feed信息                     |
| [value](#iris-query-oracle-value)          | 通过名称查询feed的执行结果           |

## iris tx oracle create

该命令用于创建一个新的feed

```bash
iris tx oracle create [flags]
```

**标识：**

| 名称, 速记        | 类型     | 必须 | 默认 | 描述                                                                              |
| ----------------- | -------- | ---- | ---- | --------------------------------------------------------------------------------- |
| --feed-name       | string   | 是   |      | feed的名称，唯一标识                                                              |
| --description     | string   |      |      | feed的描述                                                                        |
| --latest-history  | uint64   | 是   |      | feed执行结果保留的最大数目(按照时间降序保留)，范围取值为： [1, 100]               |
| --service-name    | string   | 是   |      | feed调用的服务名称.                                                               |
| --input           | string   | 是   |      | 调用服务所需要的参数，必须满足JSON格式.                                           |
| --providers       | []string | 是   |      | 服务提供者的地址列表                                                              |
| --service-fee-cap | string   | 是   |      | 单个请求愿意支付的服务费上限                                                      |
| --timeout         | int64    |      |      | 请求等待响应的最大区块数, 响应超过这个时间，请求将被忽略                          |
| --frequency       | uint64   |      |      | 重复性请求的调用频率                                                              |
| --threshold       | uint16   |      | 1    | 期待服务的最小响应数量，取值范围[1,服务提供者数量]                                      |
| --aggregate-func  | string   | 是   |      | 对 Service 响应结果进行处理的 IRISHub 预定义方法，目前支持：avg/max/min/          |
| --value-json-path | string   | 是   |      | Service响应结果中的字段名称或路径，用于从响应结果中获取调用 aggregate-func 的参数 |

### 创建一个新的feed

```bash
iris tx oracle create --chain-id=irishub --from=node0 --fees=0.3iris --feed-name="test-feed" --latest-history=10 --service-name="test-service" --input=<request-data> --providers=<provide1_address>,<provider2_address> --service-fee-cap=1iris --timeout=2 --frequency=10 --total=10 --threshold=1 --aggregate-func="avg" --value-json-path="high" -b block -y
```

## iris tx oracle start

该命令用于启动一个处于`暂停`状态的feed

```bash
iris tx oracle start <feed-name>
```

### 启动一个处于`暂停`状态的feed

```bash
iris tx oracle start test-feed --chain-id=irishub --from=node0 --fees=0.3iris -b block -y
```

## iris tx oracle pause

该命令用于暂停一个处于`运行`状态的feed

```bash
iris tx oracle pause [feed-name] [flags]
```

### 暂停一个处于`运行`状态的feed

```bash
iris tx oracle pause test-feed --chain-id=irishub --from=node0 --fees=0.3iris -b block -y
```

## iris tx oracle edit

该命令用于编辑一个已经存在的feed

```bash
iris tx oracle edit [feed-name] [flags]
```

**标识：**

| 名称, 速记        | 类型     | 必须 | 默认 | 描述                                                                |
| ----------------- | -------- | ---- | ---- | ------------------------------------------------------------------- |
| --description     | string   |      |      | feed的描述                                                          |
| --latest-history  | uint64   | 是   |      | feed执行结果保留的最大数目(按照时间降序保留)，范围取值为： [1, 100] |
| --providers       | []string | 是   |      | 服务提供者的地址列表                                                |
| --service-fee-cap | string   | 是   |      | 单个请求愿意支付的服务费上限                                        |
| --timeout         | int64    |      |      | 请求等待响应的最大区块数, 响应超过这个时间，请求将被忽略            |
| --frequency       | uint64   |      |      | 重复性请求的调用频率                                                |
| --threshold       | uint16   |      | 1    | 期待的最小响应数，取值范围[1,服务提供者数量]                        |

### 编辑feed

```bash
iris tx oracle edit test-feed --chain-id=irishub --from=node0 --fees=0.3iris --latest-history=5  -b block -y
```

## iris query oracle feed

该命令用于查询一个已存在的feed的信息

```bash
iris query oracle feed [feed-name] [flags]
```

### 查询一个已存在的feed的信息

```bash
iris query oracle feed test-feed
```

## iris query oracle feeds

该命令用于查询一组feed的信息

```bash
iris query oracle feeds [flags]
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                                   |
| ---------- | ------ | ---- | ---- | -------------------------------------- |
| --state    | string |      |      | feede状态，取值有：`paused`、`running` |

### 查询一组feed的信息

```bash
iris query oracle feeds --state=running
```

## iris query oracle value

该命令用于查询指定feed的执行结果

```bash
iris query oracle value test-feed
```

