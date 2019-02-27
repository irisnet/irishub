# iriscli stake redelegations-from

## 描述

基于某一验证者的所有重新委托查询

## 用法

```
iriscli stake redelegations-from [validator-address] [flags]
```
打印帮助信息
```
iriscli stake redelegations-from --help
```

## 示例

基于某一验证者的所有重新委托查询
```
iriscli stake redelegations-from [validator-address]
```

运行成功以后，返回的结果如下：

```json
[
  {
    "delegator_addr": "iaa10s0arq9khpl0cfzng3qgxcxq0ny6hmc9sytjfk",
    "validator_src_addr": "fva1dayujdfnxjggd5ydlvvgkerp2supknthajpch2",
    "validator_dst_addr": "fva1h27xdw6t9l5jgvun76qdu45kgrx9lqede8hpcd",
    "creation_height": "1130",
    "min_time": "2018-11-16T07:22:48.740311064Z",
    "initial_balance": {
      "denom": "iris-atto",
      "amount": "100000000000000000"
    },
    "balance": {
      "denom": "iris-atto",
      "amount": "100000000000000000"
    },
    "shares_src": "100000000000000000.0000000000",
    "shares_dst": "100000000000000000.0000000000"
  }
]
```
