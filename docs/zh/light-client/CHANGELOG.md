# IRISLCD 更新日志

## v0.16.0

*Nov 20th, 2019*

### BREAKING CHANGES

- 删除 `TokenAddition` 类型提案 API 中的 `initial_supply` 属性

### NON-BREAKING CHANGES

- 增加 Coinswap 模块
- 增加 HTLC 模块

## v0.15.0

*Aug 20th, 2019*

### 不兼容修改

在这个版本中，除'/tx/broadcast'以外的所有POST方法都只生成未签名交易，并不会对这些交易进行广播。用户需要在本地签名交易，并使用'/tx/broadcast'进行广播。

- 删除 Key 模块
- 删除 POST /tx/sign
- 删除 GET /distribution/community-tax
- 删除 GET /gov/params/{module}

### 兼容修改

- 增加 Asset 模块
- 增加 Rand 模块
- 增加 Params 模块
- 增加 GET /bank/token-stats/{id}
- 增加 POST /bank/accounts/{address}/set-memo-regexp

#### Bank 模块

| [v0.14.1]      | [v0.15.0]        | 输入改变 | 输出改变 | 备注 |
| --------------- | --------------- | --------------- | --------------- | ----- |
| GET /bank/accounts/{address} | GET /bank/accounts/{address} | No | Yes | 1. 添加 `memo_regexp` 字段; <br> 2. 添加多资产支持 |

#### Tendermint 模块

| [v0.14.1]      | [v0.15.0]        | 输入改变 | 输出改变 | 备注 |
| --------------- | --------------- | --------------- | --------------- | ---- |
| /txs/{hash} | /txs/{hash} | No | Yes | 添加 `timestamp` 字段 |
| /txs | /txs | No | Yes | 添加 `timestamp` 字段 |

## v0.14.1

*May 31th, 2019*

### 不兼容修改

- 修改 一些IRISLCD的API接口

#### Bank 模块

| [v0.14.0]      | [v0.14.1]        | 输入改变 | 输出改变 |
| --------------- | --------------- | --------------- | --------------- |
| GET /bank/accounts/{address} | GET /bank/accounts/{address} | No | Yes |

## v0.14.0

*May 20th, 2019*

### 不兼容修改

- 修改/删除 一些`IRISLCD`的API接口
- 添加了bank模块的API
- 在distribution模块中添加了 `community-tax` 和 `rewards` API接口

#### Tendermint 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- | --------------- |--------------- |
| GET /node_info  | GET /node-info   | No | No |
| GET /blocks-result/latest | GET /block-results/latest | No | No |
| GET /blocks-result/{height}  | GET /block-results/{height}  | No | No |
| POST /txs  | N/A | / | / |

#### Key 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- | --------------- | --------------- |
| GET /auth/accounts/{address}  |  N/A | / | / |

#### Sign and broadcast 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- | --------------- | --------------- |
| POST /txs/send  |  N/A | / | / |

#### Bank 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- | --------------- | --------------- |
| GET /bank/coin/{coin-type} | GET /bank/coins/{type} | No | No |
| GET /bank/token-stats | GET /bank/token-stats | No | Yes |
| GET /bank/balances/{address} | GET /bank/accounts/{address} | No | Yes |
| POST /bank/accounts/{address}/transfers | POST /bank/accounts/{address}/send | Yes | No |
| POST /bank/burn | POST /bank/accounts/{address}/burn | Yes | No |

#### Stake 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- |--------------- | --------------- |
| POST /stake/delegators/{delegatorAddr}/delegate | POST /stake/delegators/{delegatorAddr}/delegations | No | No |
| POST /stake/delegators/{delegatorAddr}/redelegate | POST /stake/delegators/{delegatorAddr}/redelegations | No | No |
| POST /stake/delegators/{delegatorAddr}/unbond | POST /stake/delegators/{delegatorAddr}/unbonding-delegations | No | No |
| GET /stake/delegators/{delegatorAddr}/unbonding_delegations | GET /stake/delegators/{delegatorAddr}/unbonding-delegations | No | No |
| GET /stake/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr} | GET /stake/delegators/{delegatorAddr}/unbonding-delegations/{validatorAddr} | No | No |
| GET /stake/validators/{validatorAddr}/unbonding_delegations | GET /stake/validators/{validatorAddr}/unbonding-delegations | No | No |

#### Slash 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- |--------------- | --------------- |
| GET /slashing/validators/{validatorPubKey}/signing_info | GET /slashing/validators/{validatorPubKey}/signing-info | No | No |
  
#### Distribution 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- |--------------- | --------------- |
| POST /distribution/{delegatorAddr}/withdrawAddress | POST /distribution/{delegatorAddr}/withdraw-address | No | No |
| GET /distribution/{delegatorAddr}/withdrawAddress | GET /distribution/{delegatorAddr}/withdraw-address | No | No |
| POST /distribution/{delegatorAddr}/withdrawReward | POST /distribution/{delegatorAddr}/rewards/withdraw | No | No |
| GET /distribution/{delegatorAddr}/distrInfo/{validatorAddr} | N/A | / | / |
| GET /distribution/{delegatorAddr}/distrInfos | N/A | / | / |
| GET /distribution/{validatorAddr}/valDistrInfo | N/A | / | / |
| N/A | GET /distribution/{address}/rewards | / | / |
| N/A | GET /distribution/community-tax | / | / | 

#### Service 模块

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- |--------------- | --------------- |
| POST /service/definition | POST /service/definitions | No | No |
| GET /service/definition/{defChainId}/{serviceName} | GET /service/definitions/{defChainId}/{serviceName} | No | No |
| POST /service/binding | POST /service/bindings | No | No |
| GET /service/binding/{defChainId}/{serviceName}/{bindChainId}/{provider} | GET /service/bindings/{defChainId}/{serviceName}/{bindChainId}/{provider} | No | No |
| PUT /service/binding/{defChainId}/{serviceName}/{provider} | PUT /service/bindings/{defChainId}/{serviceName}/{provider} | No | No |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/disable | PUT /service/bindings/{defChainId}/{serviceName}/{provider}/disable | No | No |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/enable | PUT /service/bindings/{defChainId}/{serviceName}/{provider}/enable | No | No |
| PUT /service/binding/{defChainId}/{serviceName}/{provider}/deposit/refund| PUT /service/bindings/{defChainId}/{serviceName}/{provider}/deposit/refund | No | No |
| POST /service/request | POST /service/requests | No | No |
| POST /service/response | POST /service/responses | No | No |
| GET /service/response/{reqChainId}/{reqId} | GET /service/responses/{reqChainId}/{reqId} | No | No |

#### 查询App版本

| [v0.13.1]      | [v0.14.0]        | 输入改变 | 输出改变 |
| --------------- | --------------- |--------------- | --------------- |
| GET /node_version | GET /node-version | No | No |
