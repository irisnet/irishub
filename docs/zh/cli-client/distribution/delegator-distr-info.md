# iriscli distribution delegator-distr-info

## 介绍

查询委托人全部委托的收益分配信息

## 用法

```
iriscli distribution delegator-distr-info [delegator address] [flags]
```

打印帮助信息

```
iriscli distribution delegator-distr-info --help
```

## 示例

```
iriscli distribution delegation-distr-info [delegator address]
```
执行结果示例
```json
[
  {
    "delegator_addr": "iaa1ezzh0humhy3329xg4avhcjtay985nll06lgq50",
    "val_operator_addr": "fva14a70gzu0v2w8dlfx462c9sldvja24qaz6vv4sg",
    "del_pool_withdrawal_height": "10859"
  },
  {
    "delegator_addr": "iaa1ezzh0humhy3329xg4avhcjtay985nll06lgq50",
    "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
    "del_pool_withdrawal_height": "4044"
  }
]
```