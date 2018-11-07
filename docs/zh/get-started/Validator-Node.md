# 运行一个验证人节点

在配置验证人节点之前，请保证已经按照此[文档](Install-Iris.md)正确安装了**Iris**

在IRISHub枢纽中，验证人负责将交易打包并提交区块。成为一个验证人需要满足很多条件，不仅仅是技术和硬件上的投资。同时，因为只有在有限验证人的条件下，Tendermint才能发挥最大的作用。目前，我们将IRISHub枢纽的验证人上限定为100。也就是说只有前100个验证人能够获得奖励，而大部分IRIS持有者不会成为验证人而是通过委托的方式决定谁会成为验证人。

## 如何升级成一个验证人节点

### 获取IRIS Token

#### 创建一个账户
你首先需要安装`iris` 和 `iriscli`。然后执行以下操作创建一个新的账户：

```
iriscli keys add <NAME_OF_KEY>
```

然后你需要输入至少8位的密码。

示例输出如下：
```
NAME:	TYPE:	ADDRESS:						PUBKEY:
tom	local	faa1arlugktm7p64uylcmh6w0g5m09ptvklxm5k69x	fap1addwnpepqvlmtpv7tke2k93vlyfpy2sxup93jfulll6r3jty695dkh09tekrzagazek
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

blast change tumble toddler rival ordinary chicken dirt physical club few language noise oak moment consider enemy claim elephant cruel people adult peanut garden
```

你可以查看到该账户的地址和公钥。在IRISHub中，地址经过bech32编码后将以`faa1`为首字节 ，另外公钥将以 `fap1`为首字节.

账户的助记词(seed phrase)也将被显示出来。你可以使用该长度为24个单词的助记词在任意的机器上恢复你的账户。恢复账户的命令是:

```
iriscli keys add <NAME_OF_KEY> --recover
```
### Claim tokens


一旦你完成了账户的创建，你可以通过[水龙头](https://testnet.irisplorer.io/#/faucet)获得用于测试网的IRIS token,然后你就可以将这部分IRIS用于绑定成为验证人。
水龙头每次将发送10IRIS，请按需使用！

以下命令将查询你的账户的余额：

```
iriscli bank account <ACCOUNT> --node=http://localhost:26657
```

## 执行成为验证人操作

### 确认你的全节点与网络保持同步

通过以下命令确认节点的状况：
```
iriscli status --node=tcp://localhost:26657 
```
若 `catching_up` 字段为 `false`那么你的节点就是同步的。

你需要获取当前节点的公钥信息来执行以下操作，公钥信息以 `fvp`为首字节，想要了解更多的编码信息，请参考以下 [文档](Bech32-on-IRISnet.md)

通过执行以下命令获得节点的公钥信息，公钥信息将以`fvp1`开头：

```
iris tendermint show_validator --home= < IRIS-HOME >
```
示例输出:
```
fvp1zcjduepqv7z2kgussh7ufe8e0prupwcm7l9jcn2fp90yeupaszmqjk73rjxq8yzw85
```
然后，使用以上输出作为`iriscli stake create-validator`命令的 `<pubkey>` 字段：

```
iriscli stake create-validator  --from= < name > --amount= < amount >iris --pubkey= < pubkey >  --moniker= < moniker > --fee=0.05iris  --gas=2000000 --chain-id=fuxi-4000   --node=http://localhost:26657
```
> 注意：**amount** 应为整数， **Fee** 字段可以使用小数，例如`0。01iris` 。

也就是说，如果你想要抵押1IRIS,你可以执行以下操作：

```
iriscli stake create-validator --pubkey=pubkey  --fee=0.04iris  --gas=2000000 --from= < name > --chain-id=fuxi-4000  --node=tcp://localhost:26657  --amount=1iris
```

### 查询验证人信息

你可以通过以下命令查询验证人的信息：

```
iriscli stake validator  < address-validator-operator >  --chain-id=fuxi-4000 --node=tcp://localhost:26657 
```

请注意 `<address-validator>` 字段是以`faa1`为首字母。


### 确认验证人是否在线

你可以通过以下命令查询验证人节点的运行状况，

```
iriscli status --node=tcp://localhost:26657 
```

你应该可以看到节点的`power`字段返回值大于0。

### 编辑验证人信息

你可以通过以下命令修改验证人的描述信息，验证人的名称默认为`--moniker`字段。
你应该在`details`字段注明自定义的信息。

```
iriscli stake edit-validator --from= < name >  --moniker="choose a moniker"  --website="https://irisnet.org"  --details="team" --chain-id=fuxi-4000 
  --details="details"--node=tcp://localhost:26657 --fee=0.04iris  --gas=2000000
```
### 查询验证人信息

你可以通过以下命令查询验证人的信息：

```
iriscli stake validator < address-validator-operator > --chain-id=fuxi-4000
```

### 使用浏览器：IRISPlorer

你可以通过[浏览器](https://testnet.irisplorer.io)确认验证人节点的运行状况。

### 部署IRISHub Monitor监控

请根据以下[链接](../tools/Deploy-IRIS-Monitor.md) 部署一个Monitor监控验证人。

