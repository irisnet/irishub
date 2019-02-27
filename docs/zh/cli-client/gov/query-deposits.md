# iriscli gov query-deposits

## 描述

查询指定提议的保证金详细情况

## 使用方式

```
iriscli gov query-deposits [flags]
```
打印帮助信息:

```
iriscli gov query-deposits --help
```
## 标志

| 名称, 速记       | 默认值                      | 描述                                                                                                                                                 | 是否必须  |
| --------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --proposal-id   |                            | [string] 提议ID                                                                                                        | Yes      |

## 例子

### 查询所有保证金

```shell
iriscli gov query-deposits --chain-id=<chain-id> --proposal-id=1
```

你可以查询到指定提议的所有保证金代币，包括每个存款人的充值详情。

```txt
[
  {
    "depositor": "iaa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
    "proposal_id": "1",
    "amount": [
      {
        "denom": "iris-atto",
        "amount": "35000000000000000000"
      }
    ]
  }
]
```
