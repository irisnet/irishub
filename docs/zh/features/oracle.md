# Oracle

## 简介

该模块结合`service`实现从 `Chainlink` 等可信Oracle向 IRISHub Oracle的去中心化注入。每项数据的收集任务称之为Feed，它的底层实现依赖于`service`模块。Feed的生命周期
同`service`的RequestContext基本一致(暂停，运行)，用于通过 Oracle 节点将链下数据存储在链上。另外，只能通过Profiler账户操作`Feed`，并且不能删除，只能暂停Feed。主要包括以下几种操作：

- 创建Feed
- 启动Feed
- 暂停Feed
- 编辑Feed

该模块除了通过创建Feed来收集数据，还预设了一些聚合函数，例如`avg`、`max`、`min`等，用于对收集来的数据进行加工处理以满足各种场景。每个`Feed`收集的数据，最多只保存最近的100条，其余将会被删除。

## 流程

该模块底层依赖于`service`模块，所以使用本模块的前提是执行`Service`的相关功能

   - 创建`Service`定义
   - 绑定`Service`服务。

具体说明参考[service](./service.md)。完成`service`相关操作后，开始`Oracle`流程：

1. **创建Feed**

```bash
iriscli oracle create --feed-name="test-feed" --latest-history=10 --service-name="test-service" --input={request-data} --providers="faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm,faa15rurzhkemsgfm42dnwhafjdv5s8e2pce0ku8ya" --service-fee-cap=1iris --timeout=2 --frequency=10 --threshold=1 --aggregate-func="avg" --value-json-path="high" --chain-id="irishub-test" --from=node0 --fee=0.3iris --commit
```

2. **启动Feed**

创建`Feed`之后，该收集任务处于`paused`状态，不会向服务提供者发起请求，可以通过`start`开启Feed的定时任务。

```bash
iriscli oracle start test-feed --chain-id="irishub-test" --from=node0 --fee=0.3iris --commit
```

3. **暂停Feed**

由于`Feed`一旦创建，不能够被删除，会一直消耗所有者账户的余额，直到余额耗尽，`Feed`才会进入`paused`状态。为了能够手动使`Feed`暂停，可以使用`pause`命令

```bash
iriscli oracle pause test-feed --chain-id="irishub-test" --from=node0 --fee=0.3iris --commit
```

4. **编辑Feed**

可以通过`edit`命令编辑已经存在的`feed`，改变`feed`的数据收集行为。

```bash
iriscli oracle edit test-feed --latest-history=5 --providers="faa1r3tyupskwlh07dmhjw70frxzaaaufta37y25yr,faa1ydahnhrhkjh9j9u0jn8p3s272l0ecqj40vra8h" --service-fee-cap=1iris --timeout=6 --threshold=5 --threshold=3 --chain-id="irishub-test" --from=node0 --fee=0.3iris --commit
```
需要注意的是，如果创建时的`latest-history`大于当前修改的值，oracle模块会删除多余的数据。