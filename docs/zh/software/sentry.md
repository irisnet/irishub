# 哨兵节点及其搭建

为了保证验证人节点的安全性和可用性，我们建议为验证人节点配置2个以上的哨兵节点。使用哨兵节点的好处在于可以有效地防止DoS攻击等其他针对验证人节点的攻击。

## 初始化一个全节点

为了搭建哨兵节点，首先我们需要初始化一些全节点。执行以下命令创建一个全节点(建议在多台不同的服务器上创建多个哨兵节点以提高可用性和安全性)
```
iris  init --moniker=<your name> --home=<sentry home>
```
`<sentry home>`是你指定的哨兵节点的地址。示例：
```
iris init --moniker="sentry" --home=sentry --home-client=sentry
{
  "chain_id": "test-chain-hfuDmL",
  "node_id": "937efdf8526e3d9e8b5e887fa953ff1645cc096d",
  "app_message": {
    "secret": "issue envelope dose rail busy glass treat crop royal resemble city deer hungry govern cable angle cousin during mountain december spare stick unveil great"
  }
}
```


## 修改哨兵节点的配置

然后将验证人节点中的genesis.json文件复制到 `<sentry home>/config/`目录下。接下来对`<sentry home>/config/`目录下的config.toml进行编辑。需要进行如下修改：
```
private_peers_ids="validator_node_id"
```

这里的`<validator node id>`可以在验证人节点上使用iriscli status命令获得。经过这样设置之后然后使用

```
iris init --home=<sentry home>
```

启动哨兵节点。对每个哨兵节点都需要进行这些操作。

## 修改验证人节点的配置

接下来需要对验证人节点的`<validator home>/config/`目录下的config.toml进行修改：

```
persistent_peers="sentry node id@sentry listen address" 
```

这里只写sentry节点的node id和地址，多个哨兵节点的信息使用逗号分开。

设置`pex=false` 不与其他节点进行peers交换，这样验证人节点就不会连接除persistent_peers之外的节点。
这里的`<sentry node id>`可以在哨兵节点上使用iriscli status命令获得。修改完成后需要重启验证人节点使修改生效。

```
iris  init --home=<validator node home>
```