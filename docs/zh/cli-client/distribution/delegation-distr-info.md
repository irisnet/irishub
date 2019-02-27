# iriscli distribution delegation-distr-info

## 介绍

查询某个委托的收益分配信息

## 用法

```
iriscli distribution delegation-distr-info [flags]
```

打印帮助信息

```
iriscli distribution delegation-distr-info --help
```

## 特有的flags

| 名称                | 类型   | 是否必填 | 默认值  | 功能描述        |
| --------------------| -----  | -------- | -------- | -------------- |
| --address-validator | string | true     | ""       | 验证人bech地址 |
| --address-delegator | string | true     | ""       | 委托人bech地址 |

## 示例

```
iriscli distribution delegation-distr-info --address-delegator=<delegator address> --address-validator=<validator address>
```
执行结果示例
```json
{
  "delegator_addr": "iaa1ezzh0humhy3329xg4avhcjtay985nll06lgq50",
  "val_operator_addr": "fva1ezzh0humhy3329xg4avhcjtay985nll0hpyhf4",
  "del_pool_withdrawal_height": "4044"
}
```