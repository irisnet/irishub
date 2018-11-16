# iriscli gov query-deposit

## 描述

查询保证金的充值明细

## 使用方式

```
iriscli gov query-deposit [flags]
```

## 标志

| 名称, 速记       | 默认值                 | 描述                                                                                                                                                 | 是否必须  |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --chain-id      |                       | [string] Chain ID of tendermint node                                                                                                                 | Yes      |
| --depositer     |                       | [string] bech32 depositer address                                                                                                                    | Yes      |
| --height        |                       | [int] block height to query, omit to get most recent provable block                                                                                  |          |
| --help, -h      |                       | Help for query-deposit                                                                                                                               |          |
| --indent        |                       | Add indent to JSON response                                                                                                                          |          |
| --ledger        |                       | Use a connected Ledger device                                                                                                                        |          |
| --node          | tcp://localhost:26657 | [string] \<host>:\<port> to tendermint rpc interface for this chain                                                                                  |          |
| --proposal-id   |                       | [string] proposalID of proposal depositing on                                                                                                        | Yes      |
| --trust-node    | true                  | Don't verify proofs for responses                                                                                                                    |          |
 
## 例子

### 查询充值保证金

```shell
iriscli gov query-deposit --chain-id=test --proposal-id=1 --depositer=faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07
```

通过指定提议、指定存款人查询保证金充值详情，得到结果如下：

```txt
{
  "depositer": "faa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
  "proposal_id": "1",
  "amount": [
    {
      "denom": "iris-atto",
      "amount": "30000000000000000000"
    }
  ]
}
```