# iriscli asset query-fee

## 描述

查询资产相关的费用，包括网关创建费用以及Token发行和增发费用

## 使用方式

```
iriscli asset query-fee [flags]
```

打印帮助信息:
```
iriscli asset query-fee --help
```

## 特有的标志

| 命令, 速记     | 类型   | 是否必须 | 默认值  | 描述                                          |
| --------------------| -----  | -------- | -------- | ------------------------------------------------------------------- |
| --gateway           | string | 否    | ""       | 网关名字; 如果查询网关费用,则必须指定 |
| --token             | string | 否    | ""       | Token ID; 如果查询Token费用,则必须指定  |


## 示例

```
iriscli asset query-fee --gateway=tgw
```

输出信息:
```txt
Fee: 600000iris
```

```json
{
  "exist": false,
  "fee": {
    "denom": "iris",
    "amount": "600000"
  }
}
```

```
iriscli asset query-fee --token=i.sym
```

输出信息:
```txt
Fees:
  IssueFee: 300000iris
  MintFee:  30000iris
```

```json
{
  "Exist": false,
  "issue_fee": {
    "denom": "iris",
    "amount": "300000"
  },
  "mint_fee": {
    "denom": "iris",
    "amount": "30000"
  }
}
```