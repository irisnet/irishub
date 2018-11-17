# iriscli distribution 

## 介绍

这里主要介绍distribution模块提供的命令行接口

## 用法

```
iriscli distribution [subcommand]
```

打印所以子命令和参数

```
iriscli distribution --help
```

## 子命令

| 名称                            | 功能                                                   |
| --------------------------------| --------------------------------------------------------------|
| [delegation-distr-info](delegation-distr-info.md) | 查询委托(delegation)的收益分配记录 |
| [delegator-distr-info](delegator-distr-info.md) | 查询委托人所有的委托(delegation)的收益分配记录 |
| [validator-distr-info](validator-distr-info.md) | 查询验证人收益分配记录 |
| [withdraw-address](withdraw-address.md) | 查询收益取回地址 |
| [set-withdraw-address](set-withdraw-address.md)  | 设置收益取回地址 |
| [withdraw-rewards](withdraw-rewards.md) | 发起取回收益的交易 |
