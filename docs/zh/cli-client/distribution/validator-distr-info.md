# iriscli distribution validator-distr-info

## 介绍

查询验证人的收益分配信息

## 用法

```
iriscli distribution validator-distr-info [flags]
```

打印帮助信息:

```
iriscli distribution validator-distr-info --help
```

## 特有的flags

这个命令没有特有的flag，它有一个输入参数：验证人地址


## 示例

```
iriscli distribution validator-distr-info <validator address>
```
执行结果示例
```json
[
  {
    "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
    "val_operator_addr": "fva14a70gzu0v2w8dlfx462c9sldvja24qaz6vv4sg",
    "del_pool_withdrawal_height": "10859"
  },
  {
    "delegator_addr": "faa1ezzh0humhy3329xg4avhcjtay985nll0zswc5j",
    "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
    "del_pool_withdrawal_height": "4044"
  }
]
```