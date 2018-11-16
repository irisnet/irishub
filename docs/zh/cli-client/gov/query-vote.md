# iriscli gov query-vote

## 描述

查询指定提议、指定投票者的投票情况

## 使用方式

```
iriscli gov query-vote [flags]
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | help for submit-proposal                                                                                                                             |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |
| --voter         |                            | [string] bech32 voter address                                                                                                                        | Yes      |

## 例子

### 查询投票

```shell
iriscli gov query-vote --chain-id=test --proposal-id=1 --voter=faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd
```

通过指定提议、指定投票者查询投票情况。

```txt
{
  "voter": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
  "proposal_id": "1",
  "option": "Yes"
}
```