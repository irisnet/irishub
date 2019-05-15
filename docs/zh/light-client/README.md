# IRISLCD用户文档

## 基本功能介绍

1.提供restful API以及swagger-ui

2.验证查询结果

## IRISLCD的用法

IRISLCD有两个子命令:

| 子命令      | 功能                 | 示例命令 |
| --------------- | --------------------------- | --------------- |
| version         | 打印版本信息   | `irislcd version` |
| start           | 启动一个IRISLCD节点  | `irislcd start --node=tcp://localhost:26657 --laddr=tcp://0.0.0.0:1317 --chain-id=<chain-id> --home=$HOME/.iriscli/ --trust-node` |

`start`子命令有如下参数可配置

| 参数名称        | 类型      | 默认值                 | 是否必填 | 功能描述                                          |
| --------------- | --------- | ----------------------- | -------- | ---------------------------------------------------- |
| chain-id        | string    | null                    | true     | Tendermint节点的chain ID |
| home            | string    | "$HOME/.irislcd"        | false    | 配置home目录，key和proof相关的信息都存于此 |
| node            | string    | "tcp://localhost:26657" | false    | 全节点的rpc地址 |
| laddr           | string    | "tcp://localhost:1317"  | false    | 侦听的地址和端口 |
| trust-node      | bool      | false                   | false    | 是否信任全节点 |
| max-open        | int       | 1000                    | false    | 最大连接数 |
| cors            | string    | ""                      | false    | 允许跨域访问的地址 |

## 示例命令

1. 默认情况下，IRISLCD不信任连接全节点。但是如果您确定连接的全节点是可信的，那么您应该在启动IRISLCD时加上`--trust-node`：
```bash
irislcd start --chain-id=<chain-id> --trust-node
```

2. 如果需要在其他机器上访问此IRISLCD节点，还需要配置`--laddr`参数，例如：
```bash
irislcd start --chain-id=<chain-id> --laddr=tcp://0.0.0.0:1317
```

## REST APIs

IRISLCD启动以后，您可以在浏览器中访问`localhost:1317/swagger-ui/`，然后你将看到所有的REST APIs。`swagger-ui`页面有关于API功能和所需参数的详细描述。在这里，我们只列出所有API并简要介绍它们的功能。

1. Tendermint相关APIs, 例如查询区块，交易和验证人集
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
    11. `POST /txs`: 广播交易

2. Key management模块的APIs

    1. `GET /keys`: 列出所有本地的秘钥
    2. `POST /keys`: 创建新的秘钥
    3. `GET /keys/seed`: 创建新的助记词
    4. `GET /keys/{name}`: 根据秘钥名称查询秘钥
    5. `PUT /keys/{name}`: 更新秘钥的密码
    6. `DELETE /keys/{name}`: 删除秘钥
    7. `POST /keys/{name}/recover`: 通过seed恢复一个账户

3. 签名和广播交易的APIs

    1. `POST /tx/sign`: 签名交易
    2. `POST /tx/broadcast`: 广播一个amino编码的交易
    
4. Bank模块的APIs
    1. `GET /bank/coins/{coin-type}`: 查询coin的类型信息
    2. `GET /bank/token-stats`: 查询token统计信息
    3. `GET /bank/accounts/{address}`: 查询秘钥对象账户的信息
    4. `POST /bank/accounts/{address}/send`: 发起转账交易
    5. `POST /bank/accounts/{address}/burn`: 销毁token

5. Stake模块的APIs

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

6. Governance模块的APIs

    1. `POST /gov/proposals`: 发起提交提议交易
    2. `GET /gov/proposals`: 查询提议
    3. `POST /gov/proposals/{proposalId}/deposits`: 发起缴纳押金的交易
    4. `GET /gov/proposals/{proposalId}/deposits`: 查询缴纳的押金
    5. `POST /gov/proposals/{proposalId}/votes`: 发起投票交易
    6. `GET /gov/proposals/{proposalId}/votes`: 查询投票
    7. `GET /gov/proposals/{proposalId}`: 查询某个提议
    8. `GET /gov/proposals/{proposalId}/deposits/{depositor}`:查询押金
    9. `GET /gov/proposals/{proposalId}/votes/{voter}`: 查询投票
    10. `GET /gov/params`: 查询可供治理的参数

7. Slashing模块的APIs

    1. `GET /slashing/validators/{validatorPubKey}/signing-info`: 获取验证人的签名记录
    2. `POST /slashing/validators/{validatorAddr}/unjail`: 赦免某个作恶的验证人节点

8. Distribution模块的APIs

    1. `POST /distribution/{delegatorAddr}/withdraw-address`: 设置收益取回地址
    2. `GET /distribution/{delegatorAddr}/withdraw-address`: 查询收益取回地址
    3. `POST /distribution/{delegatorAddr}/rewards/withdraw`: 取回收益
    4. `GET /distribution/{address}/rewards`: 查询收益
    5. `GET /distribution/community-tax`: 查询社区税金

9. Service模块的APIs

    1. `POST /service/definitions`: 添加服务定义
    2. `GET /service/definitions/{defChainId}/{serviceName}`: 查询服务定义
    3. `POST /service/bindings`: 添加服务绑定
    4. `GET /service/bindings/{defChainId}/{serviceName}/{bindChainId}/{provider}`: 查询服务绑定
    5. `GET /service/bindings/{defChainId}/{serviceName}`: 查询服务绑定列表
    6. `PUT /service/bindings/{defChainId}/{serviceName}/{provider}`: 更新服务绑定
    7. `PUT /service/bindings/{defChainId}/{serviceName}/{provider}/disable`: 使绑定失效
    8. `PUT /service/bindings/{defChainId}/{serviceName}/{provider}/enable`: 重新启用绑定
    9. `PUT /service/bindings/{defChainId}/{serviceName}/{provider}/deposit/refund`: 取回服务绑定的抵押
    10. `POST /service/requests`: 请求服务
    11. `GET /service/requests/{defChainId}/{serviceName}/{bindChainId}/{provider}`: 查询某服务提供者收到的服务请求
    12. `POST /service/responses`: 响应服务请求
    13. `GET /service/responses/{reqChainId}/{reqId}`: 查询服务响应
    14. `GET /service/fees/{address}`:  查询（某个地址的）服务费用
    15. `POST /service/fees/{address}/refund`: 消费者取回（未被响应的）服务费用
    16. `POST /service/fees/{address}/withdraw`: 服务提供者取回服务收益
    
10. 查询版本

    1. `GET /version`: 获取IRISHUB的版本
    2. `GET /node-version`: 查询全节点版本

## 特殊参数

这些是从部分挑选出来的可用于构建和广播交易的APIs：
1. `POST /bank/accounts/{address}/send`
2. `POST /stake/delegators/{delegatorAddr}/delegations`
3. `POST /stake/delegators/{delegatorAddr}/redelegations`
4. `POST /stake/delegators/{delegatorAddr}/unbonding-delegations`
5. `POST /gov/proposals`
6. `POST /gov/proposals/{proposalId}/deposits`
7. `POST /gov/proposals/{proposalId}/votes`
8. `POST /slashing/validators/{validatorAddr}/unjail`

上述的API都有四个特殊的查询参数，如下表所示。默认情况下，它们的值都是false。每个参数都有其唯一的优先级(这里`0`是最高优先级)。如果多个参数的值都是`true`，则将忽略优先级较低的。例如，如果`generate-only`为`true`，那么其他参数，例如`simulate`和`commit`将被忽略。

| 参数名字        | 类型 | 默认值 | 优先级 | 功能描述                 |
| --------------- | ---- | ------- |--------- |--------------------------- |
| generate-only   | bool | false | 0 | 构建一个未签名的交易并返回 |
| simulate        | bool | false | 1 | 用仿真的方式去执行交易 |
| commit          | bool | false | 2 | 等待交易被打包入块  |
| async           | bool | false | 3 | 用异步地方式广播交易  |