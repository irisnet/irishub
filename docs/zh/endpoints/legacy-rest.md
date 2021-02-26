# Legacy Amino JSON REST

irishub v1.0.0（依赖Cosmos-SDK v0.41）和更早版本提供了 REST 端点来查询状态和广播交易。 这些端点在 irishub v1.0 中仍然保留，但已标记为已弃用，并计划在几个版本后删除。因此，我们将这些端点称为 Legacy REST 端点。

Legacy REST 端点相关的重要信息：

- 这些端点中的大多数都是向后兼容的。下一部分将介绍所有的不兼容更新。
- 值得注意的是，这些端点仍会输出 Amino JSON。 Cosmos-SDK v0.41 引入了 Protobuf 作为整个代码库的默认编码库，但是传统的 REST 端点是少数几个使用 Amino 编码的部分之一。

## API 端口、激活方式和配置

所有路由都在 `~/.iris/config/app.toml` 中的以下字段中配置：

- `api.enable = true|false` 字段定义是否启用 REST 服务，默认为 `true`。
- `api.address = {string}` 字段定义服务器应绑定到的地址（实际为端口，因为主机应保持在 `0.0.0.0`）。 默认为 `tcp://0.0.0.0:1317`。
- 在 `~/.iris/config/app.toml` 中定义了一些其他 API 配置选项并附有注释，请直接参考该文件。

### Legacy REST API 路由

irishub v0.16 和更早版本中存在的 REST 路由通过 [HTTP 弃用标头](https://tools.ietf.org/id/draft-dalal-deprecation-header-01.html)标记为已弃用，它们仍然被维护以保持向后兼容，但是将在几个版本后删除。

对于应用程序开发人员而言，传统的 REST API 路由需要连接到 REST 服务器，这是通过在 ModuleManager 上调用 `RegisterRESTRoutes` 方法来完成的。

## Legacy REST 端点

### Legacy REST 端点的中不兼容更新 （对比 Cosmos-SDK v0.39 及更早的版本）

| Legacy REST 端点                                                             | 描述                                         | 不兼容更新                                                                                                                                                                                                                                                                                                                                                   |
| ---------------------------------------------------------------------------- | -------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `POST` `/txs`                                                                | 广播交易                                     | 尝试广播不支持 Amino 序列化的交易时端点将返回错误（例如 IBC 交易）<sup>1</sup>。                                                                                                                                                                                                                                                                             |
| `POST` `/txs/encode`，`POST` `/txs/decode`                                   | Amino 格式交易在 JSON 与二进制文件间的编解码 | 尝试编/解码不支持 Amino 序列化的交易时端点将返回错误（例如 IBC 交易）<sup>1</sup>。                                                                                                                                                                                                                                                                          |
| `GET` `/txs/{hash}`                                                          | 通过哈希查询交易                             | 尝试输出不支持 Amino 序列化的交易时端点将返回错误（例如 IBC 交易）<sup>1</sup>。                                                                                                                                                                                                                                                                             |
| `GET` `/txs`                                                                 | 通过事件查询交易                             | 尝试输出不支持 Amino 序列化的交易时端点将返回错误（例如 IBC 交易）<sup>1</sup>。                                                                                                                                                                                                                                                                             |
| `GET` `/gov/proposals/{id}/votes`，`GET` `/gov/proposals/{id}/votes/{voter}` | Gov 模块查询投票信息的端点                   | 所有 gov 模块端点返回值中含有投票信息的，值为 int32 而不是 string：`1=VOTE_OPTION_YES，2=VOTE_OPTION_ABSTAIN，3=VOTE_OPTION_NO，4=VOTE_OPTION_NO_WITH_VETO`.                                                                                                                                                                                                 |
| `GET` `/staking/*`                                                           | Staking 模块查询端点                         | Staking 模块中所有返回验证人的端点都有两处不兼容更新。第一，验证人 `consensus_pubkey` 字段返回 Amino 编码的 `Any` 结构，而不是 bech32 编码的公钥字符串。`Any` 的 `value` 字段是公钥原始密钥的 base64 编码的字节数组。第二，验证人的 `status` 字段目前为 int32 类型而不是 string：`1=BOND_STATUS_UNBONDED`,`2=BOND_STATUS_UNBONDING`,`3=BOND_STATUS_BONDED`。 |
| `GET` `/staking/validators`                                                  | 获取所有的验证人                             | BondStatus 现在是一个 protobuf 枚举值而不是 int32，并且 JSON 序列化时使用它的 protobuf 字段名，所以期望查询参数像 `?status=BOND_STATUS_{BONDED,UNBONDED,UNBONDING}` 而不是 `?status={bonded,unbonded,unbonding}`。                                                                                                                                           |

<sup>1</sup>： 不支持 Amino 序列化的交易是那些包含一个或多个未在 Amino 编解码器中注册的 `Msg` 的交易。 当前在 IRIShub 中，只有 IBC `Msg`s 属于这种情况。

### 迁移到新的 REST 端点 （从 Cosmos-SDK v0.39）

**IRIShub API 端点**

| Legacy REST 端点                                                                  | 描述                                 | 新的 gRPC-gateway REST 端点                                                                                   |
| --------------------------------------------------------------------------------- | ------------------------------------ | ------------------------------------------------------------------------------------------------------------- |
| `GET` `/bank/balances/{address}`                                                  | 查询一个地址的余额                   | `GET` `/cosmos/bank/v1beta1/balances/{address}`                                                               |
| `POST` `/bank/accounts/{address}/transfers`                                       | 从一个账户向另一个账户转账           | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/bank/total`                                                               | 获取所有 token 的总量                | `GET` `/cosmos/bank/v1beta1/supply`                                                                           |
| `GET` `/bank/total/{denom}`                                                       | 获取一种 token 的总量                | `GET` `/cosmos/bank/v1beta1/supply/{denom}`                                                                   |
| `GET` `/auth/accounts/{address}`                                                  | 获取账户信息                         | `GET` `/cosmos/auth/v1beta1/accounts/{address}`                                                               |
| `GET` `/staking/delegators/{delegatorAddr}/delegations`                           | 获取一个委托人的所有委托信息         | `GET` `/cosmos/staking/v1beta1/delegations/{delegator_addr}`                                                  |
| `POST` `/staking/delegators/{delegatorAddr}/delegations`                          | 提交委托                             | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/staking/delegators/{delegatorAddr}/delegations/{validatorAddr}`           | 查询一个委托人和验证人之间的委托     | `GET` `/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}`                      |
| `GET` `/staking/delegators/{delegatorAddr}/unbonding_delegations`                 | 获取一个委托人的所有解委托信息       | `GET` `/cosmos/staking/v1beta1/delegators/{delegator_addr}/unbonding_delegations`                             |
| `POST` `/staking/delegators/{delegatorAddr}/unbonding_delegations`                | 提交一个解委托                       | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/staking/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}` | 查询委托人和验证人之间的所有解委托   | `GET` `/cosmos/staking/v1beta1/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation` |
| `GET` `/staking/redelegations`                                                    | 查询重委托                           | `GET` `/cosmos/staking/v1beta1/delegators/{delegator_addr}/redelegations`                                     |
| `POST` `/staking/delegators/{delegatorAddr}/redelegations`                        | 提交一个重委托                       | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/staking/delegators/{delegatorAddr}/validators`                            | 查询一个委托人委托的所有验证人       | `GET` `/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators`                                        |
| `GET` `/staking/delegators/{delegatorAddr}/validators/{validatorAddr}`            | 查询一个委托人委托的验证人           | `GET` `/cosmos/staking/v1beta1/delegators/{delegator_addr}/validators/{validator_addr}`                       |
| `GET` `/staking/validators`                                                       | 获取所有验证人                       | `GET` `/cosmos/staking/v1beta1/validators`                                                                    |
| `GET` `/staking/validators/{validatorAddr}`                                       | 获取一个验证人信息                   | `GET` `/cosmos/staking/v1beta1/validators/{validator_addr}`                                                   |
| `GET` `/staking/validators/{validatorAddr}/delegations`                           | 获取一个验证人的所有委托信息         | `GET` `/cosmos/staking/v1beta1/validators/{validator_addr}/delegations`                                       |
| `GET` `/staking/validators/{validatorAddr}/unbonding_delegations`                 | 查询一个验证人所有未解绑的委托       | `GET` `/cosmos/staking/v1beta1/validators/{validator_addr}/unbonding_delegations`                             |
| `GET` `/staking/pool`                                                             | 获取当前抵押池的状态                 | `GET` `/cosmos/staking/v1beta1/pool`                                                                          |
| `GET` `/staking/parameters`                                                       | 获取当前 staking 模块参数值          | `GET` `/cosmos/staking/v1beta1/params`                                                                        |
| `GET` `/slashing/signing_infos`                                                   | 获取所有签名信息                     | `GET` `/cosmos/slashing/v1beta1/signing_infos`                                                                |
| `POST` `/slashing/validators/{validatorAddr}/unjail`                              | 解禁一个被监禁的验证人               | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/slashing/parameters`                                                      | 查询 slashing 模块参数值             | `GET` `/cosmos/slashing/v1beta1/params`                                                                       |
| `POST` `/gov/proposals`                                                           | 提交一个提议                         | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/gov/proposals`                                                            | 获取所有的提议                       | `GET` `/cosmos/gov/v1beta1/proposals`                                                                         |
| `POST` `/gov/proposals/param_change`                                              | 构造一个发起修改参数的提议的交易     | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/gov/proposals/{proposal-id}`                                              | 根据 ID 查询提议                     | `GET` `/cosmos/gov/v1beta1/proposals/{proposal_id}`                                                           |
| `GET` `/gov/proposals/{proposal-id}/proposer`                                     | 查询一个提议的发起者                 | `GET` `/cosmos/gov/v1beta1/proposals/{proposal_id}` (Get proposer from `Proposal` struct)                     |
| `GET` `/gov/proposals/{proposal-id}/deposits`                                     | 获取一个提议的抵押金额               | `GET` `/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits`                                                  |
| `POST` `/gov/proposals/{proposal-id}/deposits`                                    | 向一个提议抵押代币                   | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/gov/proposals/{proposal-id}/deposits/{depositor}`                         | 获取一个抵押者在一个提议中的抵押信息 | `GET` `/cosmos/gov/v1beta1/proposals/{proposal_id}/deposits/{depositor}`                                      |
| `GET` `/gov/proposals/{proposal-id}/votes`                                        | 获取一个提议的投票信息               | `GET` `/cosmos/gov/v1beta1/proposals/{proposal_id}/votes`                                                     |
| `POST` `/gov/proposals/{proposal-id}/votes`                                       | 对一个提议投票                       | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/gov/proposals/{proposal-id}/votes/{voter}`                                | 获取一个投票者的投票信息             | `GET` `/cosmos/gov/v1beta1/proposals/{proposal_id}/votes/{voter}`                                             |
| `GET` `/gov/proposals/{proposal-id}/tally`                                        | 获取一个提议的统计信息               | `GET` `/cosmos/gov/v1beta1/proposals/{proposal_id}/tally`                                                     |
| `GET` `/gov/parameters/deposit`                                                   | 获取 gov 模块抵押参数                | `GET` `/cosmos/gov/v1beta1/params/{params_type}`                                                              |
| `GET` `/gov/parameters/tallying`                                                  | 获取 gov 模块统计参数                | `GET` `/cosmos/gov/v1beta1/params/{params_type}`                                                              |
| `GET` `/gov/parameters/voting`                                                    | 获取 gov 模块投票参数                | `GET` `/cosmos/gov/v1beta1/params/{params_type}`                                                              |
| `GET` `/distribution/delegators/{delegatorAddr}/rewards`                          | 获取所有委托的奖励金额               | `GET` `/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards`                                   |
| `POST` `/distribution/delegators/{delegatorAddr}/rewards`                         | 取出所有委托人奖励                   | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}`          | 查询委托奖励                         | `GET` `/cosmos/distribution/v1beta1/delegators/{delegator_address}/rewards/{validator_address}`               |
| `POST` `/distribution/delegators/{delegatorAddr}/rewards/{validatorAddr}`         | 取出委托奖励                         | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/distribution/delegators/{delegatorAddr}/withdraw_address`                 | 获取奖励提取地址                     | `GET` `/cosmos/distribution/v1beta1/delegators/{delegator_address}/withdraw_address`                          |
| `POST` `/distribution/delegators/{delegatorAddr}/withdraw_address`                | 重设奖励提取地址                     | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/distribution/validators/{validatorAddr}`                                  | 获取一个验证人奖励分配信息           | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/distribution/validators/{validatorAddr}/outstanding_rewards`              | 获取一个验证人所有未偿付的奖励       | `GET` `/cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards`                       |
| `GET` `/distribution/validators/{validatorAddr}/rewards`                          | 获取一个验证人的佣金和自委托奖励     | `GET` `/cosmos/distribution/v1beta1/validators/{validator_address}/commission` <br> `GET` `/cosmos/distribution/v1beta1/validators/{validator_address}/outstanding_rewards`   |
| `POST` `/distribution/validators/{validatorAddr}/rewards`                         | 取出验证人的奖励                     | N/A，直接使用 Protobuf                                                                                        |
| `GET` `/distribution/community_pool`                                              | 获取社区池中持有的金额               | `GET` `/cosmos/distribution/v1beta1/community_pool`                                                           |
| `GET` `/distribution/parameters`                                                  | 获取 distribution 模块参数值         | `GET` `/cosmos/distribution/v1beta1/params`                                                                   |

**Tendermint API 端点**

| Legacy REST 端点                | 描述                                            | 新的 gRPC-gateway REST 端点                                    |
| ------------------------------- | ----------------------------------------------- | -------------------------------------------------------------- |
| `GET` `/node_info`              | 获取连接的节点的属性值                          | `GET` `/cosmos/base/tendermint/v1beta1/node_info`              |
| `GET` `/syncing`                | 获取节点的同步状态                              | `GET` `/cosmos/base/tendermint/v1beta1/syncing`                |
| `GET` `/blocks/latest`          | 获取最新高度的区块                              | `GET` `/cosmos/base/tendermint/v1beta1/blocks/latest`          |
| `GET` `/blocks/{height}`        | 获取指定高度的区块                              | `GET` `/cosmos/base/tendermint/v1beta1/blocks/{height}`        |
| `GET` `/validatorsets/latest`   | 获取最新高度的验证人集合                        | `GET` `/cosmos/base/tendermint/v1beta1/validatorsets/latest`   |
| `GET` `/validatorsets/{height}` | 获取指定高度的验证人集合                        | `GET` `/cosmos/base/tendermint/v1beta1/validatorsets/{height}` |
| `GET` `/txs/{hash}`             | 通过哈希查询交易                                | `GET` `/cosmos/tx/v1beta1/txs/{hash}`                          |
| `GET` `/txs`                    | 通过事件查询交易                                | `GET` `/cosmos/tx/v1beta1/txs`                                 |
| `POST` `/txs`                   | 广播交易                                        | `POST` `/cosmos/tx/v1beta1/txs`                                |
| `POST` `/txs/encode`            | 将 Amino JSON 格式的交易编码为 Amino 二进制格式 | N/A，直接使用 Protobuf                                         |
| `POST` `/txs/decode`            | 将 Amino 二进制格式的交易解码为 Amino JSON 格式 | N/A，直接使用 Protobuf                                         |

## 高优先级查询端点的不兼容更新

**高优先级查询端点**

- Staking
  - Validators
  - Delegators
- Bank
  - Balances
- Gov
- Auth
- Distributions

### Bank

**端点名称：** QueryBalance

- **端点路径：** `/bank/balances/{address}`
- **更新内容：**
  - 无更改
  - 有关如何为非原生的 IBC 资产创建 Token，请参见[跨链资产源路径](https://github.com/cosmos/cosmos-sdk/pull/6662)。 这将包括每种 Token 源路径的哈希值。核心决定了哈希值替换 Denom 还是在 Denom 之前

### Validators

**端点名称：** QueryValidators

- **端点路径：** `/staking/validators`
- **更新内容：**
  - `unbonding_height` 和 `jailed` 字段不再支持。
  - 如果字段为空则将被忽略，而不是返回一个职位空字符串的字段。如果验证人不配置该字段则不返回该字段。例如，在启动时验证人没有填写安全联系人，只有在他们配置之后该字段才会出现

- **JSON 示例：**

    ```JSON
    {
        "commission": {
            "commission_rates": {
                "max_change_rate": "0.000000000000000000",
                "max_rate": "0.000000000000000000",
                "rate": "0.000000000000000000"
            },
            "update_time": "1970-01-01T00:00:00Z"
        },
        "consensus_pubkey": "cosmosvalconspub1zcjduepqwuxd2yevzmsrmrjx2su8kdlk44eqfdzeqx27zejuen6m0nkcpzps0qavpw",
        "delegator_shares": "0.000000000000000000",
        "description": {
            "details": "security",
            "identity": "identity",
            "moniker": "moniker",
            "security_contact": "details",
            "website": "website"
        },
        "min_self_delegation": "1",
        "operator_address": "cosmosvaloper1pcpl7xhxq0wm72e9ljls2sxr5h3vqwytnq44sr",
        "status": 1,
        "tokens": "0",
        "unbonding_time": "1970-01-01T00:00:00Z"
    }
    ```

### Delegators

**端点名称：** QueryDelegatorDelegations

- **端点路径：** `/staking/delegators/delegations`
- **更新内容：**
  - `balance` 不再是数字类型，它包含两个子字段：`amount` 和 `Denom`
  - `delegator_address` 不再是字符串类型。更改为 `delegation` 并包含以下三个子字段：`delegator_address"、"shares"和"validator_address`
  - 旧字段 `validator_address` 不再使用。在 `redelegation` 中使用 `validator_dst_address` and `validator_src_address` 这两个新字段代替
- **JSON 示例：**

    ```JSON
    {
        "balance": {
            "amount": "5",
            "denom": "stake"
        },
        "delegation": {
            "delegator_address": "cosmos1n2k9ygw2ws9sg86mrx84pdcre5geqd5ugt44h0",
            "shares": "5.000000000000000000",
            "validator_address": "cosmosvaloper155998a4hv5kqvuxr9jryjxrtnlydvqu8c0cy03"
        }
    }
    ```

**端点名称：** QueryRedelegations

- **端点路径：** `/staking/redelegations`
- **更新内容：** 以下旧字段更改为新字段 `redelegation_entry` 的子字段：
  - `completion_time`
  - `initial_balance`
  - `shares_dst`
- `creation_height` 字段不再支持
- 以下为新字段：
  - `redelegation` 包含以下子字段
    - `delegator_address` （新增）
    - `entries` （新增）
    - `valdiator_dst_address`
    - `validator_src_address`
- **JSON 示例：**

    ```JSON
    {
        "entries": [
            {
                "balance": "5",
                "redelegation_entry": {
                    "completion_time": "1969-12-31T16:00:00-08:00",
                    "initial_balance": "5",
                    "shares_dst": "5.000000000000000000"
                }
            },
            {
                "balance": "5",
                "redelegation_entry": {
                    "completion_time": "1969-12-31T16:00:00-08:00",
                    "initial_balance": "5",
                    "shares_dst": "5.000000000000000000"
                }
            }
        ],
        "redelegation": {
            "delegator_address": "cosmos104yggz5x4ype50c59vu84ze2w36pc3swm2u698",
            "entries": null,
            "validator_dst_address": "cosmosvaloper1td8yl7g5662m0mpptaxjmcn9jtzvl0wgulvv23",
            "validator_src_address": "cosmosvaloper1gqv70e79a8q0yz5s5qhsjhdl2c79496faer0vz"
        }
    }
    ```

**端点名称：** QueryUnbondingDelegation

- **端点路径：** `/staking/unbondingDelegation`
- **更新内容：**
  - 字段 `creation_height` 不再支持

### Distributions

**端点名称：** getQueriedValidatorOutstandingRewards

- **端点路径：** `/distribution/validators/{validatorAddr}/outstanding_rewards`
- **更新内容：**
  - 新字段 `rewards` 是输出中的根字段
- **JSON 示例：**

    ```JSON
    {
        "rewards": [
            {
                "denom": "mytoken",
                "amount": "3.000000000000000000"
            },
            {
                "denom": "myothertoken",
                "amount": "0.000000300000000000"
            }
        ]
    }
    ```

**端点名称：** getQueriedValidatorCommission

- **端点路径：** `/distribution/validators/{validatorAddr}`
- **更新内容：**
  - 新字段 `commission` 是输出中的根字段
- **JSON 示例：**

    ```JSON
    {
        "commission": [
            {
                "denom": "token1",
                "amount": "4.000000000000000000"
            },
            {
                "denom": "token2",
                "amount": "2.000000000000000000"
            }
        ]
    }
    ```

**端点名称：** getQueriedValidatorSlashes

- **端点路径：** `/distribution/validators/{validatorAddr}`
- **更新内容：** 无更改
  
- **端点名称：** getQueriedDelegationRewards
- **端点路径：** `/distribution/delegators/{delegatorAddr}/rewards`
- **更新内容：** 无更改

## 构造和签名交易（完全向后兼容）

与 IRIShub 主网集成的代码相同，交易结构如下：

```json
{
    "type": "cosmos-sdk/StdTx",
    "value": {
        "msg": [
            {
                "type": "cosmos-sdk/MsgSend",
                "value": {
                    "from_address": "iaa1rkgdpj6fyyyu7pnhmc3v7gw9uls4mnajvzdwkt",
                    "to_address": "iaa1q6t5439f0rkvkzl38m0f43e0kpv3mx7x2shlq8",
                    "amount": [
                        {
                            "denom": "uiris",
                            "amount": "1000000"
                        }
                    ]
                }
            }
        ],
        "fee": {
            "amount": [
                {
                    "denom": "uiris",
                    "amount": "30000"
                }
            ],
            "gas": "200000"
        },
        "signatures": null,
        "memo": "Sent via irishub client"
    }
}
```

IRIShub 地址前缀使用 `iaa` 代替，这会影响以下字段：

- value.msg.value.from_adress
- value.msg.value.to_address

Denom 替换为 `uiris` （1iris = 10<sup>6</sup>uiris），这会影响到以下字段：

- value.msg.value.amount.denom
- value.fee.amount.denom

## 广播交易（完全向后兼容）

与 IRIShub 主网集成的代码相同，调用`POST` `/txs` 发送交易，示例如下：

```bash
curl -X POST "http://localhost:1317/txs" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"tx\": {\"msg\":[{\"type\":\"cosmos-sdk/MsgSend\",\"value\":{\"from_address\":\"iaa1rkgdpj6fyyyu7pnhmc3v7gw9uls4mnajvzdwkt\",\"to_address\":\"iaa1q6t5439f0rkvkzl38m0f43e0kpv3mx7x2shlq8\",\"amount\":[{\"denom\":\"uiris\",\"amount\":\"1000000\"}]}}],\"fee\":{\"amount\":[{\"denom\":\"uiris\",\"amount\":\"30000\"}],\"gas\":\"200000\"},\"signatures\":[{\"pub_key\":{\"type\":\"tendermint/PubKeySecp256k1\",\"value\":\"AxGagdsRTKni/h1+vCFzTpNltwoiU7SwIR2dg6Jl5a//\"},\"signature\":\"Pu8yiRVO8oB2YDDHyB047dXNArbVImasmKBrm8Kr+6B08y8QQ7YG1eVgHi5OIYYclccCf3Ju/BQ78qsMWMniNQ==\"}],\"memo\":\"Sent via irishub client\"}, \"mode\": \"block\"}"
```

## 查询交易的不兼容更新

### Tx

- **端点名称：** QueryTx
- **端点路径：** `GET /txs`&&`GET /txs/{hash}`
- **更新内容：**
  - `tags` 字段不再使用，使用 `event` 字段代替
  - `result` 字段不再使用，result 中原有字段移到第一级
  - `coin_flow` 字段不再使用

- **JSON 示例：**

  ```json
  {
      "height": "5",
      "txhash": "E663768B616B1ACD2912E47C36FEBC7DB0E0974D6DB3823D4C656E0EAB8C679D",
      "data": "0A060A0473656E64",
      "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"send\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"module\",\"value\":\"bank\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc\"},{\"key\":\"sender\",\"value\":\"iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn\"},{\"key\":\"amount\",\"value\":\"1000000uiris\"}]}]}]",
      "logs": [
          {
              "events": [
                  {
                      "type": "message",
                      "attributes": [
                          {
                              "key": "action",
                              "value": "send"
                          },
                          {
                              "key": "sender",
                              "value": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"
                          },
                          {
                              "key": "module",
                              "value": "bank"
                          }
                      ]
                  },
                  {
                      "type": "transfer",
                      "attributes": [
                          {
                              "key": "recipient",
                              "value": "iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc"
                          },
                          {
                              "key": "sender",
                              "value": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn"
                          },
                          {
                              "key": "amount",
                              "value": "1000000uiris"
                          }
                      ]
                  }
              ]
          }
      ],
      "gas_wanted": "200000",
      "gas_used": "69256",
      "tx": {
          "type": "cosmos-sdk/StdTx",
          "value": {
              "msg": [
                  {
                      "type": "cosmos-sdk/MsgSend",
                      "value": {
                          "from_address": "iaa18awn3k70u05tlcul8w2qnl64g002uj4kjn93rn",
                          "to_address": "iaa1w976a5jrhsj06dqmrh2x9qxzel74qtcmapklxc",
                          "amount": [
                              {
                                  "denom": "uiris",
                                  "amount": "1000000"
                              }
                          ]
                      }
                  }
              ],
              "fee": {
                  "amount": [
                      {
                          "denom": "uiris",
                          "amount": "30000"
                      }
                  ],
                  "gas": "200000"
              },
              "signatures": [],
              "memo": "",
              "timeout_height": "0"
          }
      },
      "timestamp": "2021-01-18T07:29:21Z"
  }
  ```
