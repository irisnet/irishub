# distribution模块的用户手册

## 简介

该模块负责将收集的交易费和通胀的代币分发给所有验证人和委托人。为了减少计算压力，引入了一种被动分配策略。
`被动`意味着不会直接向贡献者自动支付利益。委托人或者验证人必须手动发送交易以取回其收益，否则，他们的收益将保留在全局收益池中。

## 收益

### 收益的来源

1.交易所产生的交易费 `fee` (由交易的第一个签名者支付`fee`)
2.通胀的代币 `inflation`   (目前IRISnet系统设置的通胀为4%每年，通证总量为`20亿`)

### 收益的去向

1.验证人 （验证人本身也是委托人）

验证人自己委托部分的委托收益
验证人节点所有委托所产生的佣金

2.委托人

委托人的委托收益

例子：
```
假设： 当前委托"万份日收益"大约为 4iris（委托10000iris一天产生的总收益）， 验证人节点设置佣金比例 10% ，
       当前验证人自己委托 =10000iris ， 其他人委托给验证人 = 10000iris

则：节点一天的总收益：2 * 4iris = 8iris

验证人节点收益 ： 佣金 8iris*10% + 委托收益 4iris*90% =  4.4iris
其他委托人收益 ： 委托收益 4iris*90% =3.6iris
```

3.基金会

4.出块人收益
每个验证人在轮到出块的这一轮会获得`出块奖励`, 计算方法：
```
 proposerReward = (TxFee + InflatedToken) * (0.01 + 0.04 * sumPowerPrecommitValidators / totalBondedTokens)
```

## 使用场景

### 设置收益取回地址

默认情况下，收益将支付给发送委托交易的钱包。

委托人可能更改自己收益取回钱包的地址。将委托地址对应的钱包(标记为`A`)，希望收益取回钱包地址(标记为`B`)。

设置钱包B为取回收益的钱包：

```bash
iriscli distribution set-withdraw-addr <address_of_wallet_B> --fee=0.3iris --from=<key_name_of_ wallet_A> --chain-id=<chain-id>
```  

查询委托账户对应的取回收益的钱包地址：

```bash
iriscli distribution withdraw-address <address_of_wallet_A> 
```

### 取回收益

根据取回场景的不同，有3种方式可以取回收益：

1.`WithdrawDelegationRewardsAll` : 提取所有在外的委托收益（从一个或者多个验证人处）。

```bash
iriscli distribution withdraw-rewards --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
```

2.`WithdrawDelegatorReward` : 从指定验证人处提取委托收益。

```bash
iriscli distribution withdraw-rewards --only-from-validator=<validator_address>  --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
```

3.`WithdrawValidatorRewardsAll` : 仅验证人可用，同时提取自己节点的抵押收益和佣金收益。

```bash
iriscli distribution withdraw-rewards --is-validator=true --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
```

### 查询收益

根据不同场景，有2种方式查询收益：

1.使用`rewards`命令。注：此方法不需账户余额大于fee。 

执行命令：
```bash
iriscli distribution rewards <delegator_address>
```

返回如下：总收益，委托收益，佣金收益
```bash
Total:        270.33761964714393479iris
Delegations:  
  validator: iva1q7602ujxxx0urfw7twm0uk5m7n6l9gqsgw4pqy, reward: 2.899411557255275253iris
Commission:   267.438208089888659537iris
```

2.使用`dry-run`模式(模拟执行并不会广播交易)。注：此方法需要账户余额大于fee，实际执行不会扣除fee。 

执行命令(仅验证人可用)：
```bash
iriscli distribution withdraw-rewards --is-validator=true --from=node0 --dry-run --chain-id=irishub-stage --fee=0.3iris --commit
```

返回如下，`withdraw-reward-total`就是预计的抵押获益：
```bash
estimated gas = 16768
simulation code = 0
simulation log = Msg 0: 
simulation gas wanted = 50000
simulation gas used = 11179
simulation fee amount = 0
simulation fee denom = 
simulation tag action = withdraw_validator_rewards_all
simulation tag source-validator = iva1rulhmls7g9cjh239vnkjnw870t5urrut9cyrxl
simulation tag withdraw-reward-total = 2035775375047308887487iris-atto
simulation tag withdraw-address = iaa18cgtskr6cgqyyady8mumk05xk2g9c95qgw5556
simulation tag withdraw-reward-from-validator-iva1rulhmls7g9cjh239vnkjnw870t5urrut9cyrxl = 1052484144134629789682iris-atto
simulation tag withdraw-reward-commission = 983291230912679097804iris-atto
```

执行命令：
```bash
iriscli distribution withdraw-rewards --from=<key_name> --dry-run --chain-id=<chain-id> --fee=0.3iris --commit
```

返回如下：
```bash
estimated gas = 14329
simulation code = 0
simulation log = Msg 0: 
simulation gas wanted = 50000
simulation gas used = 9553
simulation fee amount = 0
simulation fee denom = 
simulation tag action = withdraw_delegation_rewards_all
simulation tag delegator = iaa1rulhmls7g9cjh239vnkjnw870t5urrutsfwvmc
simulation tag withdraw-reward-total = 1052472042330962430914iris-atto
simulation tag withdraw-address = iaa18cgtskr6cgqyyady8mumk05xk2g9c95qgw5556
simulation tag withdraw-reward-from-validator-iva1rulhmls7g9cjh239vnkjnw870t5urrut9cyrxl = 1052472042330962430914iris-atto
```

