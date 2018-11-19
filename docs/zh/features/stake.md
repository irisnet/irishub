# Stake用户手册

## 介绍

本文简要介绍了stake模块的功能以及常见用户接口。

## 核心概念

1. 投票权重

	投票权重是一个共识层面的概念。IRISHUB是一个拜占庭容错的POS区块链网络。在共识过程中，一个验证人集将对提案区块进行投票。如果验证者认为提案块有效，它将投赞成票，否则，它将投反对票。来自不同验证者的投票所占权重不同。投票的权重称为相应验证者的投票权重。
	
2. 验证人节点

	验证人节点是一个IRISHUB全节点。如果其投票权重为零，则它只是正常的全节点或验证者候选者。一旦其全节点为正数，那么它就是一个真正的验证人。

3. 委托人

	不能或不想运行验证人节点的人仍然可以作为委托人参与到POS网络中。委托人可以将token委托给验证人，委托人将从相应的验证人那里获得一定的股份。委托token也称为绑定token给验证人。稍后我们将对其进行详细说明。此外，验证人的维护者也是委托人。通常，验证人的维护者仅在其自己的验证人节点上具有股份。但它也可以在其他验证人节点上用于股份。
	
4. 候选验证人

	验证人的数量不能无限制地增加。太多的验证人可能会导致低效的共识，从而降低区块链吞吐率。因此，拜占庭容错的POS区块链网络都有验证人数量上限。通常，这个上限是100。如果有超过100个全节点升级加入验证人集。那么，只有具有前100多token数量的节点才能成功加入验证人集。其他人将是候选验证人，并将根据他们的保证金额进行降序排序。一旦一个或多个验证人被从验证人集中踢出，则顶部候选验证人将被自动添加到验证人集中。

5. 绑定，解绑定和解绑定期

	验证人维护者必须将他们自己流通的token绑定到自己的验证人节点。验证人节点投票权重与绑定的token数量成正比，包括维护者自己绑定的token和来自其他委托人的token。验证人维护者可以通过发送解绑定交易来降低他们自己的绑定token。委托人同样可以通过发送解绑定交易来降低绑定的token。但是，这些被解绑定的token不会立即成为流通的token。执行解绑定交易之后，在解绑定期结束之前，相应的验证人维护者或委托人不能再次在相同的验证人上发起解绑定交易。通常，解绑定期为三周。一旦解绑定期结束，被解绑的token将自动成为流通的token。解绑定期机制对POS区块链网络的安全性很重要。此外，如果验证人维护者在自己的验证人节点上绑定的token等于零，则相应的验证人节点会被验证器集踢出。

6. 转委托

	委托人可以将其股份从一个验证人转移到另一个验证人。这个可以分为两个步骤：从第一个验证人上解绑定和把解绑定的token绑定到另一个验证人上。正如我们上面所说，在解绑定期结束之前，解绑定操作不能立即完成，这意味着委托人不能立即发送另一个从同一验证人转委托的交易。

7. 作恶证据和惩罚

	拜占庭容错POS区块链网络假设拜占庭节点拥有不到总投票权重的1/3，而且要惩罚这些作恶节点。因此有必要收集作恶行为的证据。根据收集到的证据，stake模块将从相应的验证器和委托人中拿走一定数量的token。被拿走的token会被销毁。此外，作恶验证人将会被踢出验证人集，并被标记为关押状态，而且他们的投票权将立刻变为零。在关押期间，这些节点也不是候选验证人。当关押期结束，他们可以发送交易来解除关押状态并再次成为候选验证人。

8. 收益

	作为委托人，验证者拥有股份越多，获得的收益就越多。对于验证人的维护者，它将有额外的收益：验证人佣金。奖励来自代币通胀和交易费。至于如何计算奖励以及如何获得奖励，请参阅[mint](mint.md)和[distribution](distribution.md)。

## 用户操作

1. 运行全节点

	请参考[run_full_node](../get-started/Full-Node.md)来启动一个全节点。

2. 申请成为验证人

	这里加上您一个有一个IRISHUB的钱包，钱包里必须有一定数量的token，另外钱包的私钥已经被导入到iriscli中。

	发送申请成为验证人的交易，示例：
	```
	iriscli stake create-validator --amount=100iris --pubkey=$(iris tendermint show-validator) --moniker=<validator name> --fee=0.004iris --chain-id=<chain-id> --from=<key name> --commission-max-change-rate=0.01 --commission-max-rate=0.2 --commission-rate=0.1
	```
	`--amount`可以指定自己绑定的token数量，这个数越大你越有可能立刻成为验证人，否则只能成为候选验证人。

3. 查询自己的验证人节点

	把自己的钱包地址转换成验证人地址的编码格式
	```
	iriscli keys show [key name] --bech=val
	```
	示例返回结果：
	```
	NAME:   TYPE:   ADDRESS:                                      PUBKEY:
	faucet  local   fva1ljemm0yznz58qxxs8xyak7fashcfxf5l9pe40u    fvp1addwnpepqtdme789cpm8zww058ndlhzpwst3s0mxnhdhu5uyps0wjucaufha605ek3w
	```
	示例查询命令：
	```
	iriscli stake validator fva1ljemm0yznz58qxxs8xyak7fashcfxf5l9pe40u
	```
	返回示例：
	```text
    Validator 
    Operator Address: fva1ljemm0yznz58qxxs8xyak7fashcfxf5l9pe40u
    Validator Consensus Pubkey: fvp1zcjduepq8fw9p4zfrl5fknrdd9tc2l24jnqel6waxlugn66y66dxasmeuzhsxl6m5e
    Jailed: false
    Status: Bonded
    Tokens: 100.0000000000
    Delegator Shares: 100.0000000000
    Description: {node2   }
    Bond Height: 0
    Unbonding Height: 0
    Minimum Unbonding Time: 1970-01-01 00:00:00 +0000 UTC
    Commission: {{0.1000000000 0.2000000000 0.0100000000 0001-01-01 00:00:00 +0000 UTC}}
    ```

4. 修改验证人信息

	```
	iriscli stake edit-validator --from=<key name> --chain-id=<chain-id> --fee=0.004iris --commission-rate=0.15 --moniker=<new name>
	```

5. 增加在自己的验证人节点上委托的token

	```
	iriscli stake delegate --address-validator=<self-address-validator> --chain-id=<chain-id> --from=<key name> --fee=0.004iris  --amount=100iris 
	```

6. 委托

	向一个验证人委托一些token
	```
	iriscli stake delegate --address-validator=<other-address-validator> --chain-id=<chain-id> --from=<key name> --fee=0.004iris  --amount=100iris 
	```

7. 解绑定

	解绑定一半的token
	```
	iriscli stake unbond --address-validator=<address-validator> --chain-id=<chain-id> --from=<key name> --fee=0.004iris  --amount=100iris --share-percent=0.5
	```

8. 转委托

	转委托一半的token到另外一个验证人节点
	```
	iriscli stake redelegate --chain-id=<chain-id> --from=<key name> --fee=0.004iris --address-validator-source=<source validator address> --address-validator-dest=<destination validator address> --shares-percent=0.5
	```


对于其他查询stake状态的命令，请参考[stake_cli](../cli-client/stake/README.md)
