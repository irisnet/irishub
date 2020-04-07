---
order: 1
---

# 简介

IRIShub API服务器也称为LCD（Light Client Daemon）。IRISLCD实例是IRIShub的轻节点。与IRIShub全节点不同，它不会存储所有块并执行所有事务，这意味着它仅需要最少的带宽，计算和存储资源。在不信任模式下，它将追踪验证人集更改的过程，并要求全节点返回共识证明和Merkle证明。除非具有超过2/3投票权的验证者采取拜占庭式行为，否则IRISLCD证明验证算法可以检测到所有潜在的恶意数据，这意味着IRISLCD实例可以提供与全节点相同的安全性。

irislcd的默认主文件夹为`$HOME/.irislcd`。一旦启动IRISLCD，它将创建两个目录`keys`和`trust-base.db`，密钥存储db位于`keys`中。`trust-base.db`存储所有受信任的验证人集合以及其他与验证相关的文件。

当IRISLCD以非信任模式启动时，它将检查`trust-base.db`是否为空。如果为true，它将获取最新的块作为其信任基础，并将其保存在`trust-base.db`下。IRISLCD实例始终信任该基础。所有查询证明将在此信任的基础上进行验证，有关详细的证明验证算法，请参阅[tendermint lite](https://github.com/tendermint/tendermint/blob/master/docs/tendermint-core/light-client-protocol.md)。

## 基本功能

- 提供restful APIs并使用swagger-ui列出这些APIs。
- 验证查询证明

## 启动

IRISLCD有两个子命令:

| 子命令  | 描述                |
| ------- | ------------------- |
| version | 打印IRISLCD版本     |
| start   | 启动一个IRISLCD实例 |

`start`子命令具有以下标识:

| 标识       | 类型   | 默认值                  | 必须 | 描述                                       |
| ---------- | ------ | ----------------------- | ---- | ------------------------------------------ |
| chain-id   | string |                         | 是   | Tendermint节点的chain ID                   |
| home       | string | "$HOME/.irislcd"        |      | 配置home目录，key和proof相关的信息都存于此 |
| node       | string | "tcp://localhost:26657" |      | 全节点的rpc地址                            |
| laddr      | string | "tcp://localhost:1317"  |      | 侦听的地址和端口                           |
| trust-node | bool   | false                   |      | 是否信任全节点                             |
| max-open   | int    | 1000                    |      | 最大连接数                                 |
| cors       | string |                         |      | 允许跨域访问的地址                         |

默认情况下，IRISLCD不信任连接的完整节点。但是，如果确定所连接的完整节点是可信任的，则应使用`--trust-node`标识运行IRISLCD：

```bash
irislcd start --node=tcp://localhost:26657 --chain-id=irishub --trust-node
```

要公开访问你的IRIS LCD实例，您需要指定`--ladder`：

```bash
irislcd start --node=tcp://localhost:26657 --chain-id=irishub --laddr=tcp://0.0.0.0:1317 --trust-node
```

## REST APIs

一旦启动IRISLCD，就可以在浏览器中打开<http://localhost:1317/swagger-ui/>，然后可以浏览可用的restful APIs。swagger-ui页面包含有关APIs功能和所需参数的详细说明。在这里，我们仅列出所有API并简要介绍其功能。

:::tip
**注意**

`POST` API ([广播交易的API](#广播交易的API)除外) 只能用于生成未签名的交易，需要在[广播](#广播交易的API)之前使用其他方式其进行签名。
:::

### Tendermint相关的APIs

例如查询区块，交易和验证人集

1. `GET /node-info`: 查询所连接全节点的信息
2. `GET /syncing`: 查询所连接全节点是否处于追赶区块的状态
3. `GET /blocks/latest`: 获取最新区块
4. `GET /blocks/{height}`: 获取某一高度的区块
5. `GET /block-results/latest`: 获取最新区块交易结果
6. `GET /block-results/{height}`: 获取某一高度区块的交易结果
7. `GET /validatorsets/latest`: 获取最新的验证人集合
8. `GET /validatorsets/{height}`: 获取某一高度的验证人集合
9. `GET /txs/{hash}`: 通过交易hash查询交易
10. `GET /txs`: 搜索交易

### 广播交易的API

1. `POST /tx/broadcast`：广播一个amino或者json编码的交易

此api支持以下特殊参数。默认情况下，它们的值均为false。每个参数都有其唯一的优先级（这里的`0`是最高优先级）。如果将多个参数指定为true，则优先级较低的参数将被忽略。例如，如果`simulate`为真，则将忽略`commit`和`async`。

| 参数名称 | 类型 | 默认值 | 优先级 | 描述                            |
| -------- | ---- | ------ | ------ | ------------------------------- |
| simulate | bool | false  | 0      | 忽略gas并模拟交易，但不广播交易 |
| commit   | bool | false  | 1      | 等待交易被打包到区块中          |
| async    | bool | false  | 2      | 异步广播交易                    |

### Bank模块的APIs

1. `GET /bank/coins/{coin-type}`: 查询coin的类型信息
2. `GET /bank/token-stats`: 查询token统计信息
3. `GET /bank/token-stats/{symbol}`: 查询指定token统计信息
4. `GET /bank/accounts/{address}`: 查询链上账户信息
5. `POST /bank/accounts/{address}/send`: 发起转账交易
6. `POST /bank/accounts/{address}/burn`: 销毁token

### Stake模块的APIs

1. `POST /stake/delegators/{delegatorAddr}/delegations`: 发起委托交易
2. `POST /stake/delegators/{delegatorAddr}/redelegations`: 发起转委托交易
3. `POST /stake/delegators/{delegatorAddr}/unbonding-delegations`: 发起解委托交易
4. `GET /stake/delegators/{delegatorAddr}/delegations`: 查询委托人的所有委托记录
5. `GET /stake/delegators/{delegatorAddr}/unbonding-delegations`: 查询委托人的所有解委托记录
6. `GET /stake/delegators/{delegatorAddr}/redelegations`: 查询委托人的所有转委托记录
7. `GET /stake/delegators/{delegatorAddr}/validators`: 查询委托人的所委托的所有验证人
8. `GET /stake/delegators/{delegatorAddr}/validators/{validatorAddr}`: 查询某个被委托的验证人上信息
9. `GET /stake/delegators/{delegatorAddr}/txs`: 查询所有委托人相关的委托交易
10. `GET /stake/delegators/{delegatorAddr}/delegations/{validatorAddr}`: 查询委托人在某个验证人上的委托记录
11. `GET /stake/delegators/{delegatorAddr}/unbonding-delegations/{validatorAddr}`: 查询委托人在某个验证人上所有的解委托记录
12. `GET /stake/validators`: 获取所有验证人信息
13. `GET /stake/validators/{validatorAddr}`: 获取某个验证人信息
14. `GET /stake/validators/{validatorAddr}/delegations`:  获取某个验证人上的所有委托记录
15. `GET /stake/validators/{validatorAddr}/unbonding-delegations`: 获取某个验证人上的所有解委托记录
16. `GET /stake/validators/{validatorAddr}/redelegations`: 获取某个验证人上的所有转委托记录
17. `GET /stake/pool`: 获取权益池信息
18. `GET /stake/parameters`: 获取权益证明的参数

### Slashing模块的APIs

1. `GET /slashing/validators/{validatorPubKey}/signing-info`: 获取验证人的签名记录
2. `POST /slashing/validators/{validatorAddr}/unjail`: 解禁某个作恶的验证人节点

### Distribution模块的APIs

1. `POST /distribution/{delegatorAddr}/withdraw-address`: 设置收益取回地址
2. `GET /distribution/{delegatorAddr}/withdraw-address`: 查询收益取回地址
3. `POST /distribution/{delegatorAddr}/rewards/withdraw`: 取回收益
4. `GET /distribution/{address}/rewards`: 查询收益
5. `GET /distribution/community-tax`: 查询社区税金

### Governance模块的APIs

1. `POST /gov/proposals`: 发起提交提议交易
2. `GET /gov/proposals`: 查询提议
3. `POST /gov/proposals/{proposalId}/deposits`: 发起抵押押金的交易
4. `GET /gov/proposals/{proposalId}/deposits`: 查询抵押的押金
5. `POST /gov/proposals/{proposalId}/votes`: 发起投票交易
6. `GET /gov/proposals/{proposalId}/votes`: 查询投票
7. `GET /gov/proposals/{proposalId}`: 查询某个提议
8. `GET /gov/proposals/{proposalId}/deposits/{depositor}`:查询押金
9. `GET /gov/proposals/{proposalId}/votes/{voter}`: 查询投票

### Asset模块的APIs

1. `POST /asset/tokens`: 发行一个通证
2. `PUT /asset/tokens/{symbol}`: 编辑一个已存在的通证
3. `POST /asset/tokens/{symbol}/mint`: 增发通证到指定地址
4. `POST /asset/tokens/{symbol}/transfer`: 转让通证的所有权
5. `GET /asset/tokens/{symbol}`: 查询通证
6. `GET /asset/tokens`: 查询指定所有者的通证集合
7. `GET /asset/tokens/{symbol}/fee`: 查询发行和铸造指定通证的费用

### Coinswap模块的APIs

1. `POST /coinswap/liquidities/{voucher-coin-name}/deposit`: 增加流动性
2. `POST /coinswap/liquidities/{voucher-coin-name}/withdraw`: 提取流动性
3. `POST /coinswap/liquidities/buy`: 兑换代币(购买)
4. `POST /coinswap/liquidities/sell`: 兑换代币(出售)
5. `GET /coinswap/liquidities/{voucher-coin-name}`: 查询流动性

### HTLC模块的APIs

1. `POST /htlc/htlcs`: 创建一个HTLC
2. `GET /htlc/htlcs/{hash-lock}`: 通过hash-lock查询一个HTLC
3. `POST /htlc/htlcs/{hash-lock}/claim`: 将一个OPEN状态的HTLC中锁定的资金发放到收款人地址
4. `POST /htlc/htlcs/{hash-lock}/refund`: 从一个过期的HTLC中取回退款

### Service模块的APIs

1. `POST /service/definitions`: 定义一个新的服务
2. `GET /service/definitions/{service-name}`: 查询服务定义
3. `POST /service/bindings`: 绑定一个服务
4. `GET /service/bindings/{service-name}/{provider}`: 查询服务绑定
5. `GET /service/bindings{service-name}`: 查询服务绑定列表
6. `POST /service/providers/{provider}/withdraw-address`: 设置提取地址
7. `GET /service/providers/{provider}/withdraw-address`: 查询提取地址
8. `PUT /service/bindings/{service-name}/{provider}`: 更新一个存在的服务绑定
9. `POST /service/bindings/{service-name}/{provider}/disable`: 禁用一个可用的服务绑定
10. `POST /service/bindings/{service-name}/{provider}/enable`: 启用一个不可用的服务绑定
11. `POST /service/bindings/{service-name}/{provider}/refund-deposit`: 取回一个服务绑定的所有押金
12. `POST /service/contexts`: 发起服务调用
13. `GET /service/contexts/{request-context-id}`: 查询请求上下文
14. `PUT /service/contexts/{request-context-id}`: 更新请求上下文
15. `POST /service/contexts/{request-context-id}/pause`: 暂停一个正在进行的请求上下文
16. `POST /service/contexts/{request-context-id}/start`: 启动一个暂停的请求上下文
17. `POST /service/contexts/{request-context-id}/kill`: 终止请求上下文
18. `GET /service/requests/{request-id}`: 查询服务请求
19. `GET /service/requests/{service-name}/{provider}`: 查询一个服务绑定的活跃请求
20. `GET /service/requests/{request-context-id}/{batch-counter}`: 根据请求上下文ID和批次计数器查询请求列表
21. `POST /service/responses`: 响应服务请求
22. `GET /service/responses/{request-id}`: 查询服务响应
23. `GET /service/responses/{request-context-id}/{batch-counter}`: 根据请求上下文ID和批次计数器查询服务响应列表
24. `GET /service/fees/{provider}`: 查询服务提供者的收益
25. `POST /service/fees/{provider}/withdraw`: 提取服务提供者的收益

### Oracle模块的APIs

1. `POST /oracle/feeds`: 创建一个初始状态为`paused`的Feed。
2. `POST /oracle/feeds/<feed-name>/start`: 启动一个处于`paused`的Feed。
3. `POST /oracle/feeds/<feed-name>/pause`: 暂停一个处于`running`的Feed。
4. `PUT /oracle/feeds/<feed-name>`: 更新Feed的相关信息。
5. `GET /oracle/feeds/<feed-name>`: 通过名称查询Feed的相关信息。
6. `GET /oracle/feeds?state=<state>`: 通过状态查询符合条件的一组Feed。
7. `GET /oracle/feeds/<feed-name>/values`: 查询Feed的执行结果，按照时间倒序排列

### Rand模块的APIs

1. `POST /rand/rands`: 请求一个随机数
2. `GET /rand/rands/{request-id}`: 查询指定请求ID对应的随机数
3. `GET /rand/queue`: 查询请求队列，提供一个可选的高度参数

### Params模块的APIs

1. `GET /params`: 查询系统参数

### 查询版本

1. `GET /version`: 查询IRISLCD的版本
2. `GET /node-version`: 查询IRISLCD所连接全节点的版本
