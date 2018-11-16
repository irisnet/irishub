# iriscli gov query-proposal

## 描述

查询单个提议的详情

## 使用方式

```
iriscli gov query-proposal [flags]
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                            | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --height        |                            | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                            | Help for query-proposal                                                                                                                              |          |
| --indent        |                            | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                            | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657      | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                            | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                       | Don't verify proofs for responses                                                                                                                    |          |

## 例子

### 查询指定的提议

```shell
iriscli gov query-proposal --chain-id=test --proposal-id=1
```

查询指定提议的详情，可以得到结果如下：

```txt
{
  "proposal_id": "1",
  "title": "test proposal",
  "description": "a new text proposal",
  "proposal_type": "Text",
  "proposal_status": "DepositPeriod",
  "tally_result": {
    "yes": "0.0000000000",
    "abstain": "0.0000000000",
    "no": "0.0000000000",
    "no_with_veto": "0.0000000000"
  },
  "submit_time": "2018-11-14T09:10:19.365363Z",
  "deposit_end_time": "2018-11-16T09:10:19.365363Z",
  "total_deposit": [
    {
      "denom": "iris-atto",
      "amount": "49000000000000000050"
    }
  ],
  "voting_start_time": "0001-01-01T00:00:00Z",
  "voting_end_time": "0001-01-01T00:00:00Z",
  "param": {
    "key": "",
    "value": "",
    "op": ""
  }
}
```