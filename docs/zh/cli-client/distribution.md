# Distribution

distribution模块用于管理自己的 [Staking 收益](../concepts/general-concepts.md#staking-收益)。

## 可用命令

| 名称                                                                                      | 描述                                                                                           |
| ----------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| [commission](#iris-query-distribution-commission)                                         | 查询分配的验证人佣金                                                                                |
| [community-pool](#iris-query-distribution-community-pool)                                 | 查询社区池总币数                                                                   |
| [params](#iris-query-distribution-params)                                                 | 查询分配参数                                                                                   |
| [rewards](#iris-query-distribution-rewards)                                               | 查询所有分销委托人收益或来自指定验证人的收益 |
| [slashes](#iris-query-distribution-slashes)                                               | 查询验证人指定块范围内的分割                                                                                   |
| [validator-outstanding-rewards](#iris-query-distribution-validator-outstanding-rewards)      | 查询验证人的未付奖励分配及其所有授权                                                                                   |
| [fund-community-pool](#iris-tx-distribution-fund-community-pool)                          | 为社区基金池提供指定数额的资金                                                                                  |
| [set-withdraw-addr](#iris-tx-distribution-set-withdraw-addr)                              | 设置提现地址                                                                                   |
| [withdraw-all-rewards](#iris-tx-distribution-withdraw-all-rewards)                        | 取回委托人所有收益                                                                                   |
| [withdraw-rewards](#iris-tx-distribution-withdraw-rewards)                                | 取回收益，有以下几种模式: 取回所有奖励、从指定的验证者取回委派奖励、验证人取回所有奖励以及佣金  |

## iris query distribution commission

查询分配的验证人佣金。

```bash
iris query distribution commission [validator] [flags]
```

## iris query distribution community-pool

查询社区池总币数。

```bash
iris query distribution community-pool [flags]
```

## iris query distribution params

查询分配参数。

```bash
 iris query distribution params [flags]
```

## iris query distribution rewards

查询所有分销委托人收益或来自指定验证人的收益。

```bash
iris query distribution rewards [delegator-addr] [validator-addr] [flags]
```

## iris query distribution slashes

查询验证人指定块范围内的分割。

```bash
iris query distribution slashes [validator] [start-height] [end-height] [flags]
```

## iris query distribution validator-outstanding-rewards

查询验证人的未付奖励分配及其所有授权。

```bash
iris query distribution validator-outstanding-rewards [validator] [flags]
```
## iris tx distribution fund-community-pool

为社区基金池提供指定数额的资金。

```bash
iris tx distribution fund-community-pool [amount] [flags]
```
## iris tx distribution set-withdraw-addr

设置提现地址。

```bash
iris tx distribution set-withdraw-addr [withdraw-addr] [flags]
```

## iris tx distribution withdraw-all-rewards

取回委托人所有收益。

```bash
iris tx distribution withdraw-all-rewards [flags]
```

## iris tx distribution withdraw-rewards

取回收益，有以下几种模式: 取回所有奖励、从指定的验证者取回委派奖励、验证人取回所有奖励以及佣金。

```bash
iris tx distribution withdraw-rewards [validator-addr] [flags]
```
