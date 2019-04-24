# iriscli gov query-proposal

## 描述

查询单个提议的详情

## 使用方式

```
iriscli gov query-proposal <flags>
```
打印帮助信息:

```
iriscli gov query-proposal --help
```

## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | 提议ID                                                                                                        | Yes      |
## 例子

### 查询指定的提议

```shell
iriscli gov query-proposal --chain-id=<chain-id> --proposal-id=<proposal-id>
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
