# Staking

## 简介

本文简要介绍了staking模块的功能以及常见用户接口。

## 概念

### 投票权重

投票权重是一个共识层面的概念。IRIShub是一个拜占庭容错的POS区块链网络。在共识过程中，一个验证人集将对提案区块进行投票。如果验证人认为提案块有效，它将投赞成票，否则，它将投反对票。来自不同验证人的投票所占权重不同。投票的权重称为相应验证人的投票权重。

### 验证人节点

验证人节点是一个IRIShub全节点。如果其投票权重为零，则它只是一般的全节点或候选验证人。一旦其投票权重为正数，那么它就是一个真正的验证人。

### 委托人

不能或不想运行验证人节点的人仍然可以作为委托人参与到POS网络中。委托人可以将token委托给验证人，委托人将从相应的验证人那里获得一定的token份额。委托token也称为绑定token给验证人。稍后我们将对其进行详细说明。此外，验证节点的所有者也是委托人。验证节点的所有者不仅可以在其自己的验证节点上抵押token，而且也可以在其他验证节点上抵押token。

:::danger
**验证节点的所有者在解绑自己抵押的代币时，切勿完全解绑。 一旦完全解绑，该验证人节点将被处于jailed状态，该节点将收不到任何奖励或者佣金， 在该节点上委托代币的投资人的利益也会收到相应的损失。 所以，无论如何请保留至少1iris在抵押状态。**
**如果一旦验证人全部解委托，可以通过重新`delegate`和`unjail`的命令来恢复**
:::

### 候选验证人

验证人的数量不能无限增加。太多验证人可能会导致低效的共识，从而降低区块链吞吐率。因此，拜占庭容错的POS区块链网络都有验证人数量上限。通常，这个上限是100。如果有超过100个全节点申请加入验证人集，那么只有具有抵押token数量排名前100的节点才能成为真正的验证人，其他人将是候选验证人，并将根据他们抵押token的数量进行降序排序。一旦一个或多个验证人被从验证人集中踢出，则顶部候选验证人将被自动添加到验证人集中。

### 绑定，解绑和解绑期

验证人节点的所有者必须将他们自己流通的token绑定到自己的验证人节点。验证人节点投票权重与绑定的token数量成正比，包括所有者自己绑定的token和来自其他委托人的token。
验证人节点的所有者可以通过发送解绑交易来降低自己的绑定token。委托方还可以通过发送解绑交易来降低自己的绑定token。然而，这些解除了绑定的token不会立即变成流通token。
执行解除绑定交易后，相应的验证人节点或委托方在解除绑定周期结束之前不能在相同的验证人节点上再次发送解除绑定交易。
通常解除绑定的周期为三周。解除绑定周期结束后，解除绑定的token将自动成为流通token。解绑定周期机制对POS区块链网络的安全性做出了很大的贡献。此外，如果一个验证人节点绑定的token数量等于零，则此验证人节点将从验证人节点集合中删除。

### 转委托

委托人可以将其令牌从一个验证节点转移到另一个验证节点。转委托可以分为两个步骤：从第一个验证节点处取消绑定并绑定到另一个验证节点。如上所述，在解除绑定期结束之前无法立即完成未绑定操作，这意味着委托人无法立即发送其他转委托交易。

### 证据&&惩罚

拜占庭容错的POS区块链网络假定拜占庭节点的投票权不到总投票权的1/3。这些拜占庭节点必须受到惩罚。因此，有必要收集拜占庭行为的证据。根据证据，放样模块将自动从相应的验证者和委托者处削减一定数量的令牌。此外，拜占庭验证人节点将从验证人集合中删除并投入监狱，这意味着其投票权为零。在监禁期间，这些节点不是验证人的候选对象。监禁期结束后，他们可以发送交易取消监禁并再次成为验证者候选人。

### 奖励

作为委托人，他在验证人节点上拥有的令牌越多，它将获得的奖励就越多。对于验证人，它将获得额外的奖励：验证者佣金。奖励来自令牌通货膨胀和交易费。至于如何计算奖励以及如何获得奖励，请参考[mint](mint.md)和[distribution](distribution.md)。

## 用户操作

- 查询自己的验证人节点

查询验证人地址的编码格式的钱包地址：

```bash
iris keys show <key-name>
```

 示例输出：

  ```bash
  - name: node0
    type: local
    address: iaa1w9lvhwlvkwqvg08q84n2k4nn896u9pqx93velx
    pubkey: iap1addwnpepq03g7u43y3gwfz3pd4gkwz7d4mt600kzsc5cj2ysx58a5hp84qyduxtw28r
    mnemonic: ""
    threshold: 0
    pubkeys: []
  ```

查询验证人信息：

```bash
iris q staking validator iva14n9md3sq9xwscs96za8n85m0j9y2yu3cagxgke
```

 示例输出：

```json
{
    "operator_address": "iva14n9md3sq9xwscs96za8n85m0j9y2yu3cagxgke",
    "consensus_pubkey": "icp1zcjduepq9meszzqu54gpxvs4vzvuv85qvv5ef0egz3sde0ps4dvktcv77uds0kkhgf",
    "status": 3,
    "tokens": "100000000",
    "delegator_shares": "100000000.000000000000000000",
    "description": {
        "moniker": "node0"
    },
    "unbonding_time": "1970-01-01T00:00:00Z",
    "commission": {
        "commission_rates": {
            "rate": "1.000000000000000000",
            "max_rate": "1.000000000000000000",
            "max_change_rate": "1.000000000000000000"
        },
        "update_time": "2020-08-26T06:43:07.065305Z"
    },
    "min_self_delegation": "1"
}
```

- 修改验证人信息

```bash
iris tx staking edit-validator --from=<key-name> --chain-id=irishub --fees=0.3iris --commission-rate=0.15 --moniker=<new-name>
```

- 委托

```bash
iris tx staking delegate iva14n9md3sq9xwscs96za8n85m0j9y2yu3cagxgke 1000iris --chain-id=irishub --from=<key-name> --fees=0.3iris
```

- 解绑

```bash
iris tx staking unbond iva14n9md3sq9xwscs96za8n85m0j9y2yu3cagxgke 1000iris --chain-id=irishub --from=<key-name> --fees=0.3iris
```

- 转委托

```bash
iris tx staking redelegate iva14n9md3sq9xwscs96za8n85m0j9y2yu3cagxgke iva1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm 100iris --from mykey --chain-id=irishub --from=<key-name> --fees=0.3iris
```

对于其它Staking相关的命令，请参考[stake-cli](../cli-client/staking.md)
