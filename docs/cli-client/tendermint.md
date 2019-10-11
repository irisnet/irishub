# iriscli tendermint

Tendermint state querying subcommands

## Available Commands

| Name, shorthand                                    | Description                                           |
| -------------------------------------------------- | ----------------------------------------------------- |
| [tx](#iriscli-tendermint-tx)                       | Matches this txhash over all committed blocks         |
| [txs](#iriscli-tendermint-txs)                     | Search for all transactions that match the given tags |
| [block](#iriscli-tendermint-block)                 | Get verified data for a the block at given height     |
| [validator-set](#iriscli-tendermint-validator-set) | Get the full tendermint validator set at given height |

## iriscli tendermint tx

### Search transaction by hash

```bash
iriscli tendermint tx CD117378EC1CE0BA4ED0E0EBCED01AF09DA8F6B7 --chain-id=irishub
```

## iriscli tendermint txs

Search all transactions which match the given tag list

```bash
iriscli tendermint txs <flags>

```

**Flags:**

| Name, shorthand | Default | Description                             | Required |
| --------------- | ------- | --------------------------------------- | -------- |
| --tags          |         | Tag: value list of tags that must match |          |
| --page          | 0       | Pagination page                         |          |
| --size          | 100     | Pagination size                         |          |

### Search transactions

```bash
iriscli tendermint txs --tags="action:send&sender:iaa1c6al0vufl8efggzsvw34hszua9pr4qqymthxjw" --chain-id=irishub
```

## Actions list

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

Get verified data of a the block at the given height. If no height is specified, the latest height will be used as default.

### Get block info by height

```bash
iriscli tendermint block <block-height> --chain-id=irishub
```

### Get the latest block info

```bash
iriscli tendermint block --chain-id=irishub
```

## iriscli tendermint validator-set

Get the full tendermint validator set at given height. If no height is specified, the latest height will be used as default.

### Get validator set by height

```bash
 iriscli tendermint validator-set <block-height> --chain-id=irishub
```

### Get the latest validator-set

```bash
 iriscli tendermint validator-set --chain-id=irishub
```
