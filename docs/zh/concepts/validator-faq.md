---
order: 2
---

# 验证人问答

:::tip
IRISnet 基础概念请参考 [基础概念](general-concepts.md)
:::

## 常见问题

### 如何成为 IRIShub 验证人

参考 [加入主网](../get-started/mainnet.md)

同时也推荐您运行主网验证人之前，在[测试网](../get-started/testnet.md)上面测试验证人相关操作。

:::tip
想要获取更多的委托，您可以：

- 进行安全审计，并公开审计结果
- 开源一些 IRIShub 相关的开发工具和工作流程
- 建立自己的网站，以建立良好的品牌形象
:::

### 最低硬件要求是什么

最低硬件要求请参考：[硬件要求](../daemon/intro.md#硬件要求)

### 验证人节点的状态有哪些

通过`create-validator`交易创建验证人后，它们可能会处于以下三种状态：

- `bonded`：验证人在活跃验证人集合中，可参与共识并获得奖励，若行为不当会导致抵押的部分通证被罚没。
- `unbonding`：
  - 验证人节点行为异常且已被监禁（Jail），即被移出活跃验证人集合。监禁时长是确定的，超过监禁时长后，验证人可以发送`unjail`交易以解除监禁。
  - 验证人抵押的IRIS数量（包括受委托）脱离了前100名，因此而成为候选人。可以通过抵押或委托更多的IRIS使自己进入前100名，他将即时重新获得验证人身份。
- `unbonded`：验证人不在活跃验证人集合中，因此不能参与共识。此类验证人不会受到惩罚，也不会获得任何奖励。但仍然可以接受委托。

### IRIShub 有哪些不同类型的密钥

简而言之，有两种类型的密钥：

#### Tendermint 密钥

这是用于共识投票签名的唯一密钥。

- 由`iris init`创建节点时生成
- 查询bech32前缀为 `icp` 的共识公钥

  ```bash
  iris tendermint show-validator
  ```

- 查询bech32前缀为 `ica` 的节点地址

  ```bash
  iris tendermint show-address
  ```

- 查询对应的节点Hex地址

  ```bash
  iriscli status | jq .validator_info.address
  ```
  
  它也存储在`config/priv_validator.json`中。

#### 应用程序密钥

该密钥可以通过`iriscli`创建，用于签名交易。应用程序密钥与以`iap`为前缀的公钥和以`iaa`为前缀的地址相关联。两者都是由`iriscli keys add`生成的帐户密钥派生出来的。

注意：验证人操作员的密钥直接与应用程序密钥绑定，地址和公钥分别使用保留的前缀：iva和ivp。

### 如何备份验证人节点

安全备份验证人节点私钥非常**重要**，这是恢复验证人节点的唯一方法。请注意，这里指的是[Tendermint密钥](#tendermint-key)。

如果您使用的是软件签名（tendermint的默认签名方法），则您的[Tendermint密钥](#tendermint-key)位于`<iris-home>/config/priv_validator.json`中。最简单的方法是备份整个config文件夹。

或者，您可以使用硬件更安全地管理[Tendermint密钥](#tendermint-key)，例如[YubiHSM2](https://developers.yubico.com/YubiHSM2/)。

### 如何迁移/恢复验证人节点

迁移验证人的方法有很多，最推荐的方法是：

1.在新服务器上运行[运行全节点](../get-started/mainnet.md#运行全节点)

2.追赶上最新区块之后，停止验证人节点和全节点

3.将全节点的 `config` 文件夹替换为验证人的

4.启动新的验证人节点

### 什么是“自抵押”？我如何增加“自抵押”

自抵押就是给自己验证人节点的委托。可以通过从用于创建验证人的操作员帐户发送“委托”交易来增加此金额。

### 想要成为活跃的验证人最少要抵押多少IRIS

最低抵押 `1iris` 即可创建验证人，但能否成为活跃的验证人取决于您的抵押（包括受委托）数量是否超过第100名验证人。

### 验证人可以卷走委托人的资金吗

通过委托IRIS给一个验证人节点，用户同时委托了对应的投票权。验证人的投票权越高，他们在共识和治理流程中的权重就越大。但这并不意味着验证人可以托管其委托人的资金。 **验证人绝不可能动用其委托人的资金**。

即使验证人无法窃取委托资金，但如果验证人行为不当，委托人仍会被动承担责任和惩罚。

### 多久选择一次验证人来提议下一个区块？它是否与抵押的IRIS数量有关

被选择来提议下一个区块的验证人称为提议者。验证人被选为提议者的概率是确定的，与验证人的投票权（即绑定的IRIS数量）成正比。例如，如果所有验证人的投票权总额为 100 IRIS，而其中一个验证人的投票权为10，则此验证人将提议出约10％的区块。

### 激励计划是什么

请参考[Staking 收益](./general-concepts.md#staking-收益)

### 运行验证人节点有哪些收益

除了自抵押收益，验证人还可以从受委托的收益中赚取佣金。

验证人在治理中也起着重要作用。如果委托人未投票，则他们将从其验证人那里继承投票。这使验证人在生态系统中负有主要责任。

### 什么是验证人佣金

验证人收益池中的收益会在验证人及其委托人之间分配，验证人可以对委托人获取的收益收取一定比例的佣金。每个验证人可以自由设置其初始佣金，每日修改佣金幅度上限和佣金比例上限。IRIShub强制每个验证人设置该参数，创建验证人后，只有佣金比例可以更改。

### 收益计算公式是什么

请参考[Staking 收益计算公式](general-concepts.md#staking-收益计算公式)

### 如何查询我的验证人地址

验证人地址有2种：

- 验证人操作员地址，即用来创建验证人节点的[应用程序密钥](#应用程序密钥)

  查询验证人的操作员地址（iva ...）和pubkey（ivp ...）：

  ```bash
  iriscli keys show MyKey --bech=val
  ```

- 验证人节点地址，即[Tendermint密钥](#Tendermint-密钥)

  查询关联的地址（ica ...）和共识pubkey（icp ...），请参考[Tendermint密钥](#Tendermint-密钥)

### 如何将验证人的Logo上传到[区块浏览器](../get-started/explorers.md)

1. 使用验证人的名称注册[Keybase](https://keybase.io/)

2. 上传Logo作为Keybase帐户的头像

3. 点击`Add a PGP key`，创建一个PGP私钥，完成后您将获得一个16位的字符串

4. [编辑验证人](../cli-client/stake.md#iriscli-stake-edit-validator)信息并指定`--identity=<16位的PGP字符串>`

## 常见错误

### 验证人的投票权为0

可能是您的验证人被监禁或由于抵押（包括受委托）的IRIS数量排到了前100名以外。

如果您的验证人被监禁，可以按以下步骤操作来恢复：

- 如果`iris`没有运行，请重新启动：

  ```bash
  iris start
  ```

- 等待节点赶上最新的区块，检查验证人会被监禁到什么时间（了解[监禁时长](gov-params.md#slashing-模块可治理参数)）：

  ```bash
  # 查询验证人节点共识公钥
  iris tendermint show-validator --home=<iris-home>

  # 使用共识公钥查询节点状态
  iriscli stake signing-info <validator-consensus-public-key>
  ```

  您将可以看到 `Jailed Until` 的时间，只有在该时间之后，您才可以执行接下来的步骤。

- 如果当前时间已经超过了 `Jailed Until`，即可执行[解禁](../cli-client/stake.md#iriscli-stake-unjail)操作：
  
  ```bash
  iriscli stake unjail --from=<key-name> --fee=0.3iris --chain-id=irishub
  ```

- 再次检查您的验证人，看看您的投票权是否恢复。

  ```bash
  iris status
  ```

  您可能会注意到您的投票权比以前低，那是因为你被惩罚了（了解[罚款金额](gov-params.md#slashing-模块可治理参数)）。

### `iris` 异常退出：too many open files

Linux可以打开（每个进程）的默认文件数是 `1024`，而 `iris` 进程会打开超过1024个文件，进而导致进程崩溃。一个快速的解决方法是执行 `ulimit -n 4096`（增加允许的打开文件数量，仅对当前会话有效），然后使用 `iris start` 重新启动。如果您使用的是systemd或其他进程管理器来启动 `iris`，则最好在该级别进行一些配置。

- 示例`systemd`配置：
  
    ```toml
    # /etc/systemd/system/iris.service
    [Unit]
    Description=IRIS Hub Node
    After=network.target

    [Service]
    Type=simple
    User=ubuntu
    WorkingDirectory=/home/ubuntu
    ExecStart=/home/ubuntu/go/bin/iris start
    Restart=on-failure
    RestartSec=3
    LimitNOFILE=65535

    [Install]
    WantedBy=multi-user.target
    ```

- 在Ubuntu系统中修改全局ulimit示例:

    ```bash
    # Edit limits.conf
    vim /etc/security/limits.conf

    # Append the following lines at the bottom
    * hard nofile 65535
    * soft nofile 65535
    root hard nofile 65535
    root soft nofile 65535

    # Reboot the system
    reboot

    # Re-login & Check whether ulimit is updated to 65535
    ulimit -n
    ```

### Uptime始终为0％，即使节点已经完成同步

比较两个`Consensus Pubkey`：

- 从[区块浏览器](https://www.irisplorer.io/#/validators)中，您可以在“验证人详情”页中找到该验证人声明的`Consensus Pubkey`。

- 通过 `iris tendermint show-validator --home=<iris-home>` 检查正在使用的`Consensus Pubkey`。

如果它们不相同，则意味着您正在运行的只是一个普通全节点，而不是验证人节点。

#### 最好的情况是您已经备份了[Tendermint 密钥](#如何备份验证人节点)

那么您可以执行以下操作：

- 停止节点
- 用您的备份替换当前的 `<iris-home>/config/priv_validator.json`
- 通过`iris tendermint show-validator --home=<iris-home>` 确认 `Consensus Pubkey` 是正确的
- 启动节点
- 完成同步后，检查`voting_power`现在应该大于0：`iriscli status`

#### 如果我丢失了Tendermint密钥怎么办

这意味着您 **永远失去了您的验证人！**您只能创建一个新的验证人，并将所有通证[转委托](../cli-client/stake.md#iriscli-stake-redelegate)给新的验证人。

## 加入社区

- 英文 Riot：<https://riot.im/app/#/room/#irisvalidators:matrix.org>
- 中文 QQ： 834063323
