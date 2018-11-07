# 参与Genesis文件生成


1. 每个希望成为验证人的参与者确保请根据一下[教程](Install-Iris.md) 在服务器上完成**Iris**的安装。

2. 执行gentx命令，获得一个node-id.json的文件。这个操作将默认生成一个余额为200IRIS的账户，该账户默认绑定100IRIS成为一个验证人候选人。

```
iris init gen-tx --name=your_name --home=<path_to_home> --ip=Your_public_IP
```
* 代码示例：
```
iris init gen-tx --name=alice 
```

```
       {
         "app_message": {
           "secret": "village venue about lend pause popular vague swarm blue unusual level drastic field broken moral north repair blue accident miss essay loan rail harbor"
         },
         "gen_tx_file": {
           "node_id": "1b45f5bb7ba1e00be01e8795dcaa0e8008f28cb5",
           "ip": "192.168.150.206",
           "validator": {
             "pub_key": {
               "type": "tendermint/PubKeyEd25519",
               "value": "NlMcGgz05K45ukGY10R8DApp8A0N0Jv4F2/OKtq9fCU="
             },
             "power": "100",
             "name": ""
           },
           "app_gen_tx": {
             "name": "tom",
             "address": "faa1mmnaknf87p7uu80m6uthyssd2ge0s73hcfr05h",
             "pub_key": "fap1zcjduepqxef3cxsv7nj2uwd6gxvdw3rups9xnuqdphgfh7qhdl8z4k4a0sjsxh3kgg"
           }
         }
       }
  ```
然后你可以发现在$IRISHOME/config目录下生成了一个gentx文件夹。里面存在一个gentx-node-ID.json文件。这个文件包含了如下信息：

```
{
    "node_id": "612db83e7facdd9abab879f7e465ed829f3f3487",
    "ip": "192.168.150.223",
    "validator": {
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "bzLIySQ4YDwBIkTgeyrnBx7VEoQ23zDnWhIV4FEEOZ4="
      },
      "power": "100",
      "name": ""
    },
    "app_gen_tx": {
      "name": "haoyang-virtualbox2",
      "address": "faa1k96h5cyppg6q2meftv6epuw39u5dd0sa8t84fv",
      "pub_key": "fap1zcjduepqduev3jfy8psrcqfzgns8k2h8qu0d2y5yxm0npe66zg27q5gy8x0qh7wt9l"
    }
  }
```
validator字段对应了home/config下的节点信息

`app_gen_tx`中说明了拥有这个节点的账户信息。这个账户的助记词就是刚刚的secret

3. 将上述提到的json文件以提交Pull Request的形式上传到`https://github.com/irisnet/testnets/tree/master/testnets/fuxi-4000/config/gentx`目录下：

> 注意:json文中的IP改成公网IP

4. 在收集完参与者的gentx文件后，团队将在一下目录公布fuxi-4000测试网的配置文件：`https://github.com/irisnet/testnets/tree/master/testnets/fuxi-4000/config`。然后你就可以下载genesis.json和config.toml文件了。


