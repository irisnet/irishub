# 运行一个验证人节点

在配置验证人节点之前，请保证已经按照此[文档](../software/How-to-install-irishub.md) 正确安装了`iris`, 且已经部署了全节点并完成同步。

在IRIShub枢纽中，验证人负责将交易打包并提交区块。成为一个验证人需要满足很多条件，不仅仅是技术和硬件上的投资。同时，因为只有在有限验证人的条件下，Tendermint才能发挥最大的作用。目前，我们将IRIShub枢纽的验证人上限定为100。也就是说只有前100个（根据委托量的多少排序）验证人能够获得奖励，而大部分IRIS持有者不会成为验证人而是通过委托的方式来获取收益并决定谁会成为验证人。

## 如何升级成一个验证人节点

### 创建一个账户
你首先需要安装正确版本的`iris` 和 `iriscli`。如果你没有可用的账户，需要执行以下操作创建一个新的账户。你也可以使用已有的账户。

```
iriscli keys add <key_name>
```

然后你需要输入至少8位的密码。

示例输出如下：
```
NAME:	TYPE:	ADDRESS:						PUBKEY:
tom	local	iaa1arlugktm7p64uylcmh6w0g5m09ptvklxrmsz9m	iap1addwnpepqvlmtpv7tke2k93vlyfpy2sxup93jfulll6r3jty695dkh09tekrz37h9q9
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

blast change tumble toddler rival ordinary chicken dirt physical club few language noise oak moment consider enemy claim elephant cruel people adult peanut garden
```

你可以查看到该账户的地址和公钥。在IRIShub中，地址经过bech32编码后将以`iaa1`为首字节 ，另外公钥将以 `iap1`为首字节.

账户的助记词(seed phrase)也将被显示出来。你可以使用该长度为24个单词的助记词在任意的机器上恢复你的账户。恢复账户的命令是:

```
iriscli keys add <key_name> --recover
```
### 获得IRIS代币(仅限Testnet 测试网)


一旦你完成了账户的创建，你可以通过[水龙头](https://testnet.irisplorer.io/#/faucet)获得用于测试网的IRIS token,然后你就可以将这部分IRIS用于绑定成为验证人。
水龙头每次将发送10IRIS，请按需使用！

以下命令将查询你的账户的余额：

```
iriscli bank account <account_address> --node=http://localhost:26657
```

## 执行成为验证人操作

### 确认你的全节点与网络保持同步

通过以下命令确认节点的状况：
```
iriscli status --node=tcp://localhost:26657 
```
若 `catching_up` 字段为 `false`那么你的节点就是同步的。

你需要获取当前节点的公钥信息来执行以下操作，公钥信息以 `icp`为首字节，想要了解更多的编码信息，请参考以下 [文档](../features/basic-concepts/bech32-prefix.md)

通过执行以下命令获得节点的公钥信息，公钥信息将以`icp`开头：

```
iris tendermint show-validator --home=<iris_home>
```
示例输出:
```
icp1zcjduepq9l2svsakh9946n42ljt0lxv0kpwrc4v9c2pnqhn9chnjmlvagans88ltuj
```
然后，使用以上输出作为`iriscli stake create-validator`命令的 `<pubkey>` 字段.

如果你想要抵押10IRIS并且在IRISnet链上创建验证人,你可以执行以下操作：

::: warning
**创建验证人需要更多的gas和fee， 和普通交易略有不同， 需要自己指定 --gas=100000 --fee=0.6iris**
**验证人上限是100人，如果超出这个数量则会按照委托抵押数量的多少排序（取前100人），所以尽量多抵押**
:::

```
iriscli stake create-validator --chain-id=<chain-id> --from=<key name> --gas=100000 --fee=0.6iris --pubkey=<validator public key> --amount=10iris --moniker=<your_custom_name> --commission-rate=0.1 --identity=<identity_string>
```
> 注意：
> 
> 1.`Fee`和`amount` 字段可以使用小数，例如`1.01iris` ,其中`amount`（验证人抵押额）需大于`0.5iris`。
> 
> 2.`identity`为可选字段，可以在[keybase](https://keybase.io/)上传头像并获取相对应的字段填入, 例如`--identity="357F80896B3311B4"`。


### 查询验证人信息

你可以通过以下命令查询验证人的信息：

```
iriscli stake validator <address-validator-operator> --chain-id=<chain-id> --node=tcp://localhost:26657 
```

请注意 `<address-validator-operator>` 字段是以`iva1`为首字母。


### 确认验证人是否在线

你可以通过以下命令查询验证人节点的运行状况，

```
iriscli status --node=tcp://localhost:26657 
```

你应该可以看到节点的`voting_power`字段返回值大于0。

### 编辑验证人信息

你可以通过以下命令修改验证人的描述信息，验证人的名称默认为`--moniker`字段。
你应该在`details`字段注明自定义的信息。

```
iriscli stake edit-validator --from=<key name> --moniker=<your_custom_name> --website=<your_website> --details=<your_details> --chain-id=<chain-id> --node=tcp://localhost:26657 --fee=0.3iris --identity=<identity_string>
```

> `identity`为可选字段。

### 使用浏览器

你可以通过[浏览器](https://www.irisplorer.io)确认验证人节点的运行状况。
