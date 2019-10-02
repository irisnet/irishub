# iriscli tendermint

Tendermint状态查询子命令

## 可用命令

| 名称，速记                                           | 描述                                                                       |
| ---------------------------------------------------- | -------------------------------------------------------------------------- |
| [tx](#iriscli-tendermint-tx)                         | 在所有提交的区块上寻找匹配此txhash的交易                                   |
| [txs](#iriscli-tendermint-txs)                       | 搜索查询所有符合指定匹配条件的交易                                         |
| [block](#iriscli-tendermint-block)                   | 在给定高度获取区块的验证数据。如果未指定高度，则将使用最新高度作为默认高度 |
| [validator-set](#iriscli-tendermint-validator-set)   | 查询指定高度的验证人信息                                                   |
| [show-address](#iriscli-tendermint-show-address)     | 查询验证人的私钥对应的地址                                                 |
| [show-validator](#iriscli-tendermint-show-validator) | 获取验证人的私钥对应的公钥                                                 |

## iriscli tendermint tx

### 按交易Hash查询交易

```bash
iriscli tendermint tx CD117378EC1CE0BA4ED0E0EBCED01AF09DA8F6B7 --chain-id=irishub
```

## iriscli tendermint txs

搜索所有匹配指定条件的交易

```bash
iriscli tendermint txs <flags>

```

**标志：**

| 名称，速记 | 默认 | 描述                    | 必须 |
| ---------- | ------ | ----------------------- | -------- |
| --tags     |        | Tag：必须匹配的标签列表 |          |
| --page     | 0      | 分页页码                |          |
| --size     | 100    | 分页大小                |          |

### 搜索交易

```bash
iriscli tendermint txs --tags="action:send&sender:iaa1c6al0vufl8efggzsvw34hszua9pr4qqymthxjw" --chain-id=irishub
```

## 操作列表

| module       | Msg                                           | action                          |
| ------------ | --------------------------------------------- | ------------------------------- |
| asset        | irishub/asset/MsgCreateGateway                | create_gateway                  |
|              | irishub/asset/MsgEditGateway                  | edit_gateway                    |
|              | irishub/asset/MsgTransferGatewayOwner         | transfer_gateway_owner          |
|              | irishub/asset/MsgIssueToken                   | issue_token                     |
|              | irishub/asset/MsgEditToken                    | edit_token                      |
|              | irishub/asset/MsgMintToken                    | mint_token                      |
|              | irishub/asset/MsgTransferTokenOwner           | transfer_token_owner            |
| bank         | irishub/bank/Send                             | send                            |
|              | irishub/bank/Burn                             | burn                            |
| distribution | irishub/distr/MsgModifyWithdrawAddress        | set_withdraw_address            |
|              | irishub/distr/MsgWithdrawDelegationRewardsAll | withdraw_delegation_rewards_all |
|              | irishub/distr/MsgWithdrawDelegationReward     | withdraw_delegation_reward      |
|              | irishub/distr/MsgWithdrawValidatorRewardsAll  | withdraw_validator_rewards_all  |
| gov          | irishub/gov/MsgSubmitProposal                 | submit_proposal                 |
|              | irishub/gov/MsgSubmitTxTaxUsageProposal       | submit_proposal                 |
|              | irishub/gov/MsgSubmitAddTokenProposal         | submit_proposal                 |
|              | irishub/gov/MsgDeposit                        | deposit                         |
|              | irishub/gov/MsgVote                           | vote                            |
| stake        | irishub/stake/MsgCreateValidator              | create_validator                |
|              | irishub/stake/MsgEditValidator                | edit_validator                  |
|              | irishub/stake/MsgDelegate                     | delegate                        |
|              | irishub/stake/BeginUnbonding                  | begin_unbonding                 |
|              | irishub/stake/BeginRedelegate                 | begin_redelegate                |
| slashing     | irishub/slashing/MsgUnjail                    | unjail                          |
| service      | irishub/service/MsgSvcDef                     | define_service                  |
|              | irishub/service/MsgSvcBinding                 | bind_service                    |
|              | irishub/service/MsgSvcBindingUpdate           | update_service_binding          |
|              | irishub/service/MsgSvcDisable                 | disable_service                 |
|              | irishub/service/MsgSvcEnable                  | enable_service                  |
|              | irishub/service/MsgSvcRefundDeposit           | refund_service_deposit          |
|              | irishub/service/MsgSvcRequest                 | call_service                    |
|              | irishub/service/MsgSvcResponse                | respond_service                 |
|              | irishub/service/MsgSvcRefundFees              | refund_service_fees             |
|              | irishub/service/MsgSvcWithdrawFees            | withdraw_service_fees           |
|              | irishub/service/MsgSvcWithdrawTax             | withdraw_service_tax            |

## iriscli tendermint block

在给定高度获取区块的验证数据。如果未指定高度，则将使用最新高度作为默认高度。

### 按高度获取区块信息

```bash
iriscli tendermint block <block-height> --chain-id=irishub
```

### 获取最新的区块信息

```bash
iriscli tendermint block --chain-id=irishub
```

## iriscli tendermint validator-set

在指定高度获取全部验证人。如果未指定高度，则将使用最新高度作为默认高度。

### 按高度获取验证人

```bash
 iriscli tendermint validator-set <block-height> --chain-id=irishub
```

### 获取最新的验证人

```bash
 iriscli tendermint validator-set --chain-id=irishub
```
