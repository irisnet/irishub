# iriscli tendermint txs

## 介绍

搜索查询所有符合指定匹配条件的交易

## 用法

```shell
iriscli tendermint txs <flags>
```

## 标志

| 名称，速记   | 默认值                | 功能介绍         | 是否必填 |
| ------------ | --------------------- | ---------------- | -------- |
| --chain-id   | ""                    | 区块链Chain ID   | yes      |
| --node       | tcp://localhost:26657 | 节点查询rpc接口  |          |
| --help, -h   |                       | 帮助信息         |          |
| --trust-node | true                  | 是否信任查询节点 |          |
| --tags       | ""                    | 匹配条件         |          |
| --page       | 0                     | 分页的页码       |          |
| --size       | 100                   | 分页的大小       |          |

## 示例

### 查询交易

```shell
iriscli tendermint txs --tags "action:send&sender:iaa1c6al0vufl8efggzsvw34hszua9pr4qqymthxjw" --chain-id=<chain-id> --trust-node
```

示例结果：

```json
{
  "hash": "CD117378EC1CE0BA4ED0E0EBCED01AF09DA8F6B7",
  "height": "100722",
  "tx": {
    "type": "auth/StdTx",
    "value": {
      "msg": [
        {
          "type": "cosmos-sdk/Send",
          "value": {
            "inputs": [
              {
                "address": "iaa1c6al0vufl8efggzsvw34hszua9pr4qqymthxjw",
                "coins": [
                  {
                    "denom": "iris-atto",
                    "amount": "3650000000000000000"
                  }
                ]
              }
            ],
            "outputs": [
              {
                "address": "iaa1v2ezk7yvkgjq87ey54etfuxc87353ulr50vlzc",
                "coins": [
                  {
                    "denom": "iris-atto",
                    "amount": "3650000000000000000"
                  }
                ]
              }
            ]
          }
        }
      ],
      "fee": {
        "amount": [
          {
            "denom": "iris-atto",
            "amount": "4787310000000000"
          }
        ],
        "gas": "6631"
      },
      "signatures": [
        {
          "pub_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A/ZQqJkDnqiN7maj4N9we8u8hE1dUpFD72+bD2PZgH+V"
          },
          "signature": "MEQCIEiNg0y3Xp9YgpY00cuYV6yoRIIXS1/Z7rOJeRwK8WipAiABfHZAS/yDMqPnBEPud1eJX8cZ6hhex1C7CGq286oclw==",
          "account_number": "162",
          "sequence": "3"
        }
      ],
      "memo": ""
    }
  },
  "result": {
    "log": "Msg 0: ",
    "gas_wanted": "6631",
    "gas_used": "4361",
    "tags": [
      {
        "key": "c2VuZGVy",
        "value": "ZmFhMWM2YWwwdnVmbDhlZmdnenN2dzM0aHN6dWE5cHI0cXF5cnkzN2pu"
      },
      {
        "key": "cmVjaXBpZW50",
        "value": "ZmFhMXYyZXprN3l2a2dqcTg3ZXk1NGV0ZnV4Yzg3MzUzdWxydnEyOHo5"
      },
      {
        "key": "Y29tcGxldGVDb25zdW1lZFR4RmVlLWlyaXMtYXR0bw==",
        "value": "MzE0ODQ2MzExNDE2MDc2MA=="
      }
    ]
  }
}

```

## Actions 列表

| module       | Msg                                            | action                          |
| ------------ | ---------------------------------------------- | ------------------------------- |
| bank         | irishub/bank/MsgSend                           | send                            |
|              | irishub/bank/MsgBurn                           | burn                            |
|              | irishub/bank/MsgSetMemoRegexp                  | set-memo-regexp                 |
| distribution | irishub/distr/MsgSetWithdrawAddress            | set_withdraw_address            |
|              | irishub/distr/MsgWithdrawDelegationRewardsAll  | withdraw_delegation_rewards_all |
|              | irishub/distr/MsgWithdrawDelegationReward      | withdraw_delegation_reward      |
|              | irishub/distr/MsgWithdrawValidatorRewardsAll   | withdraw_validator_rewards_all  |
| gov          | irishub/gov/MsgSubmitProposal                  | submit_proposal                 |
|              | irishub/gov/MsgSubmitSoftwareUpgradeProposal   | submit_proposal                 |
|              | irishub/gov/MsgSubmitCommunityTaxUsageProposal | submit_proposal                 |
|              | irishub/gov/MsgSubmitTokenAdditionProposal     | submit_proposal                 |
|              | irishub/gov/MsgDeposit                         | deposit                         |
|              | irishub/gov/MsgVote                            | vote                            |
| stake        | irishub/stake/MsgCreateValidator               | create_validator                |
|              | irishub/stake/MsgEditValidator                 | edit_validator                  |
|              | irishub/stake/MsgDelegate                      | delegate                        |
|              | irishub/stake/BeginUnbonding                   | begin_unbonding                 |
|              | irishub/stake/BeginRedelegate                  | begin_redelegate                |
| slashing     | irishub/slashing/MsgUnjail                     | unjail                          |
| asset        | irishub/asset/MsgCreateGateway                 | create_gateway                  |
|              | irishub/asset/MsgEditGateway                   | edit_gateway                    |
|              | irishub/asset/MsgTransferGatewayOwner          | transfer_gateway_owner          |
|              | irishub/asset/MsgIssueToken                    | issue_token                     |
|              | irishub/asset/MsgEditToken                     | edit_token                      |
|              | irishub/asset/MsgMintToken                     | mint_token                      |
|              | irishub/asset/MsgTransferTokenOwner            | transfer_token_owner            |
| coinswap     | irishub/coinswap/MsgSwapOrder                  | swap_order                      |
|              | irishub/coinswap/MsgAddLiquidity               | add_liquidity                   |
|              | irishub/coinswap/MsgRemoveLiquidity            | remove_liquidity                |
| htlc         | irishub/htlc/MsgCreateHTLC                     | create_htlc                     |
|              | irishub/htlc/MsgClaimHTLC                      | claim_htlc                      |
|              | irishub/htlc/MsgRefundHTLC                     | refund_htlc                     |
| service      | irishub/service/MsgSvcDef                      | define_service                  |
|              | irishub/service/MsgSvcBinding                  | bind_service                    |
|              | irishub/service/MsgSvcBindingUpdate            | update_service_binding          |
|              | irishub/service/MsgSvcDisable                  | disable_service                 |
|              | irishub/service/MsgSvcEnable                   | enable_service                  |
|              | irishub/service/MsgSvcRefundDeposit            | refund_service_deposit          |
|              | irishub/service/MsgSvcRequest                  | call_service                    |
|              | irishub/service/MsgSvcResponse                 | respond_service                 |
|              | irishub/service/MsgSvcRefundFees               | refund_service_fees             |
|              | irishub/service/MsgSvcWithdrawFees             | withdraw_service_fees           |
|              | irishub/service/MsgSvcWithdrawTax              | withdraw_service_tax            |
| rand         | irishub/rand/MsgRequestRand                    | request_rand                    |
