# iriscli distribution

distribution模块用于管理自己的 [Staking 收益](../concepts/general-concepts.md#staking-收益)。

## 可用命令

| 名称                                                            | 描述                                                                                           |
| --------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| [withdraw-address](#iriscli-distribution-withdraw-address)      | 查询提现地址                                                                                   |
| [rewards](#iriscli-distribution-rewards)                        | 查询验证人或委托人的所有奖励                                                                   |
| [set-withdraw-address](#iriscli-distribution-set-withdraw-addr) | 设置提现地址                                                                                   |
| [withdraw-rewards](#iriscli-distribution-withdraw-rewards)      | 取回收益，有以下几种模式: 取回所有奖励、从指定的验证者取回委派奖励、验证人取回所有奖励以及佣金 |

## iriscli distribution withdraw-address

查询委托人的提现地址。

```bash
iriscli distribution withdraw-address <delegator-address> <flags>
```

### 查询提现地址

```bash
iriscli distribution withdraw-address <delegator-address>
```

如果委托人未指定提现地址，则查询结果为空。

## iriscli distribution rewards

查询验证人或委托人的所有奖励。

```bash
iriscli distribution rewards <address> <flags>
```

### 查询奖励

```bash
iriscli distribution rewards <iaa...>
```

输出:

```bash
Total:        270.33761964714393479iris
Delegations:  
  validator: iva..., reward: 2.899411557255275253iris
  validator: iva..., reward: 2.899411557255275253iris
  validator: iva..., reward: 2.899411557255275253iris
Commission:   267.438208089888659537iris
```

## iriscli distribution set-withdraw-addr

设置另一个地址以接收奖励，而不是使用委托人地址。

```bash
iriscli distribution set-withdraw-addr <withdraw-address> <flags>
```

### 设置提现地址

```bash
iriscli distribution set-withdraw-addr <iaa...> --from=<key-name> --fee=0.3iris --chain-id=irishub
```

## iriscli distribution withdraw-rewards

取回奖励到提现地址（默认为委托人地址，您可以通过 [set-withdraw-addr](#iriscli-distribution-set-withdraw-addr)重新设置提现地址)。

```bash
iriscli distribution withdraw-rewards <flags>
```

**标识：**

| 名称, 速记            | 类型   | 必须 | 默认 | 描述                                 |
| --------------------- | ------ | ---- | ---- | ------------------------------------ |
| --only-from-validator | string |      |      | 仅从此验证者地址中提取（以bech格式） |
| --is-validator        | bool   |      | 否   | 同时取回验证人的佣金                 |

:::tip
不要同时指定以上两个标志。
:::

### 从指定的验证者取回委派奖励

```bash
iriscli distribution withdraw-rewards --only-from-validator=<validator-address> --from=<key-name> --fee=0.3iris --chain-id=irishub
```

### 取回所有奖励

```bash
iriscli distribution withdraw-rewards --from=<key-name> --fee=0.3iris --chain-id=irishub
```

### 验证人取回所有奖励以及佣金

```bash
iriscli distribution withdraw-rewards --is-validator=true --from=<key-name> --fee=0.3iris --chain-id=irishub
```
