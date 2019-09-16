# iriscli tendermint txs

## Description

Search all transactions which match the given tag list

## Usage

```
iriscli tendermint txs <flags>

```

## Flags

| Name, shorthand | Default               | Description                                                   | Required |
| --------------- | --------------------- | ------------------------------------------------------------- | -------- |
| --chain-id      |                       | Chain ID of Tendermint node                                   | true     |
| --node string   | tcp://localhost:26657 | Node to connect to                                            |          |
| --help, -h      |                       | Help for txs                                                  |          |
| --trust-node    | true                  | Trust connected full node (don't verify proofs for responses) |          |
| --tags          |                       | Tag: value list of tags that must match                       |          |
| --page          | 0                     | Pagination page                                               |          |
| --size          | 100                   | Pagination size                                               |          |

## Examples

### Search transactions

```shell
iriscli tendermint txs --tags=`action:send&sender:iaa1c6al0vufl8efggzsvw34hszua9pr4qqymthxjw` --chain-id=<chain-id> --trust-node
```

You will get the following result.

```
[
  {
    "hash": "50F8D75FC1F0C2643A0D09189B7FB44246AB00AF89779215FFBC0740E6C59F3A",
    "height": "3411",
    "tx": {
      "type": "irishub/bank/StdTx",
      "value": {
        "msg": [
          {
            "type": "irishub/bank/Send",
            "value": {
              "inputs": [
                {
                  "address": "faa10t6tn0ntgrzetmzwlr9x8fj4j29qrcax0p52dm",
                  "coins": [
                    {
                      "denom": "iris-atto",
                      "amount": "10000000000000000000"
                    }
                  ]
                }
              ],
              "outputs": [
                {
                  "address": "faa1m9m9t8paa48xgmaxg7gxzq3a5rcl4neecm4f94",
                  "coins": [
                    {
                      "denom": "iris-atto",
                      "amount": "10000000000000000000"
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
              "amount": "400000000000000000"
            }
          ],
          "gas": "50000"
        },
        "signatures": [
          {
            "pub_key": {
              "type": "tendermint/PubKeySecp256k1",
              "value": "AtRMRPAKXstV6/NN8cizi55lMeqrtyzkR6UmSMcujYpG"
            },
            "signature": "ctjaNgszonLxoVd2weWe1TleCxg8vmSoYNuJNI1OEE5Ll/+NY0PEnDHeUsTkq71t8HgYkFkM636EssP9TAmttQ==",
            "account_number": "2",
            "sequence": "1"
          }
        ],
        "memo": ""
      }
    },
    "result": {
      "Code": 0,
      "Data": null,
      "Log": "Msg 0: ",
      "Info": "",
      "GasWanted": "50000",
      "GasUsed": "6678",
      "Tags": [
        {
          "key": "action",
          "value": "send"
        },
        {
          "key": "sender",
          "value": "faa10t6tn0ntgrzetmzwlr9x8fj4j29qrcax0p52dm"
        },
        {
          "key": "recipient",
          "value": "faa1m9m9t8paa48xgmaxg7gxzq3a5rcl4neecm4f94"
        }
      ],
      "Codespace": "",
      "XXX_NoUnkeyedLiteral": {},
      "XXX_unrecognized": null,
      "XXX_sizecache": 0
    },
    "timestamp": "2019-07-01T07:40:05Z"
  }
]
```

## Actions list
| module          | Msg                  | action                                                    |
| --------------- | -------------------- | --------------------------------------------------------- |
| asset        | irishub/asset/MsgCreateGateway | create_gateway |
|              | irishub/asset/MsgEditGateway | edit_gateway |
|              | irishub/asset/MsgTransferGatewayOwner | transfer_gateway_owner |
|              | irishub/asset/MsgIssueToken | issue_token |
|              | irishub/asset/MsgEditToken | edit_token |
|              | irishub/asset/MsgMintToken | mint_token |
|              | irishub/asset/MsgTransferTokenOwner | transfer_token_owner |
| bank         | irishub/bank/Send | send |
|              | irishub/bank/Burn | burn |
| distribution | irishub/distr/MsgModifyWithdrawAddress | set_withdraw_address |
|              | irishub/distr/MsgWithdrawDelegationRewardsAll | withdraw_delegation_rewards_all |
|              | irishub/distr/MsgWithdrawDelegationReward | withdraw_delegation_reward |
|              | irishub/distr/MsgWithdrawValidatorRewardsAll | withdraw_validator_rewards_all |
| gov          | irishub/gov/MsgSubmitProposal | submit_proposal |
|              | irishub/gov/MsgSubmitTxTaxUsageProposal | submit_proposal |
|              | irishub/gov/MsgSubmitAddTokenProposal | submit_proposal |
|              | irishub/gov/MsgDeposit | deposit |
|              | irishub/gov/MsgVote | vote |
| stake        | irishub/stake/MsgCreateValidator | create_validator |
|              | irishub/stake/MsgEditValidator | edit_validator |
|              | irishub/stake/MsgDelegate | delegate |
|              | irishub/stake/BeginUnbonding | begin_unbonding |
|              | irishub/stake/BeginRedelegate | begin_redelegate |
| slashing     | irishub/slashing/MsgUnjail | unjail |
| service      | irishub/service/MsgSvcDef | define_service |
|              | irishub/service/MsgSvcBinding | bind_service |
|              | irishub/service/MsgSvcBindingUpdate | update_service_binding |
|              | irishub/service/MsgSvcDisable | disable_service |
|              | irishub/service/MsgSvcEnable | enable_service |
|              | irishub/service/MsgSvcRefundDeposit | refund_service_deposit |
|              | irishub/service/MsgSvcRequest | call_service |
|              | irishub/service/MsgSvcResponse | respond_service |
|              | irishub/service/MsgSvcRefundFees | refund_service_fees |
|              | irishub/service/MsgSvcWithdrawFees | withdraw_service_fees |
|              | irishub/service/MsgSvcWithdrawTax | withdraw_service_tax |
