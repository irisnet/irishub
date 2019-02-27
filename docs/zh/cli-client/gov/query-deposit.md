# iriscli gov query-deposit

## 描述

查询保证金的充值明细

## 使用方式

```
iriscli gov query-deposit [flags]
```
打印帮助信息:

```
iriscli gov query-deposit --help
```
## 标志

| 名称, 速记       | 默认值                 | 描述                                                                                                                                                 | 是否必须  |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --depositor     |                       | [string] bech32编码的存款人地址                                                                                                                    | Yes      |
| --proposal-id   |                       | [string] 提议ID                                                                                                        | Yes      |

## 例子

### 查询充值保证金

```shell
iriscli gov query-deposit --chain-id=test --proposal-id=1 --depositor=iaa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07
```

通过指定提议、指定存款人查询保证金充值详情，得到结果如下：

```txt
{
  "depositor": "iaa1c4kjt586r3t353ek9jtzwxum9x9fcgwetyca07",
  "proposal_id": "1",
  "amount": [
    {
      "denom": "iris-atto",
      "amount": "30000000000000000000"
    }
  ]
}
```
