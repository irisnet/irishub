# IRISHub重的链上治理过程

## 什么是链上治理?

链上治理是让验证人对区块链网络的运行达成共识的一种投票机制。

### 线上治理提案的类型

* Text文字类型
* System Parameter Change系统参数修改
* Protocol Upgrade 软件征集

## 线上治理流程

链上治理投票分为以下三步骤：


* Mininimum Depost: `50IRIS`
* Deposit Period: 100 blocks
* Penalty for non-voting validtors: 1%
* Pass Threshold: 50%
* Voting Period: 100 blocks

## 如何提交一个提案?

任何人都可以提交链上治理提案，但你需要为此提案存入超过最低要求的押金。

如下的命令将执行提交一个 `Text`类型的提案:

```
iriscli gov submit-proposal --title="Text" --description="name of the proposal" --type="Text" --deposit="1000000000000000000000iris" --proposer=<account>  --from=<name>  --chain-id=fuxi-3001 --fee=400000000000000iris --gas=20000 --node=http://localhost:36657
```

The `<account>` for `proposer` field should start with `faa` which corresponds to `<name>`.


## 如何增加投票的抵押金额?

To add deposit to some proposal, you could execute this command to add `10IRIS` to the proposal's deposit:

```
iriscli gov deposit --proposalID=1 --depositer=<account> --deposit=1000000000000000000iris   --from=<name>  --chain-id=fuxi-3001  --fee=400000000000000iris --gas=20000  --node=http://localhost:36657 
```

##如何投票?

In the current version of governance module, you have the following choices for each proposal:
* Yes
* No
* NoWithVeto
* Abstien

You could put one of the choices in the `--option` field. 

To vote for a proposal, you need to get the correct `<proposal_id>`.You could execute the following command to vote on proposal with ID = 1:
```
iriscli  vote --from=jerry --voter=<account> --proposalID=1 --option=Yes --chain-id=fuxi-3001   --fee=2000000000000000iris --gas=20000  --node=http://localhost:36657
```

## 如何查询投票信息?

例如，查询第一个提案的信息：

```
iriscli gov query-proposal --proposalID=1 --chain-id=fuxi-3001 --node=http://localhost:26657

``````
也可以在浏览器上查询。
