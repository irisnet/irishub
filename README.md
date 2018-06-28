# IRIS-Hub

IRIS Hub是在Cosmos生态中的区域性枢纽，提供iService服务

## 下载发行版安装

- 进入下载页: https://github.com/irisnet/irishub/releases/
- 下载对应版本的可执行文件
- 解压缩`tar -C /usr/local/bin -xzf iris$VERSION.$OS-$ARCH.zip`
- 拷贝到`/usr/local/bin/`目录下
  执行

```
iris version
```

若出现对应的版本号则说明安装成功



## 通过源码编译安装

### 系统要求

操作系统：**Ubuntu 16.04** 

架构：**amd64**

要求安装1.9版本以上的Go https://golang.org/doc/install

### 安装步骤

```
go get github.com/irisnet/irishub
cd $GOPATH/src/github.com/irisnet/irishub
make all
iris version
```

若出现对应的版本号则说明安装成功

### 本地测试网络

1. 初始化irishub:

```
MYNAME1=<your name>
iris init --chain-id=test --name=$MYNAME1 --home=$HOME/.iris1
```

以上命令将在 `$HOME/.iris1／config`目录下创建默认的配置文件 `priv_validator.json`和`genesis.json`。`priv_validator.json` 包含验证人的私钥。`genesis.json`文件包含其他配置信息， 包括创世账户，创世账户的默认密码为`1234567890`，创世账户的助记词会在执行以上命令后输出在控制台，例如：

```
{
  "chain_id": "test",
  "node_id": "e4bcb3c7df9783d32c8a876028a45952b4146c78",
  "app_message": {
    "secret": "spare detail grass load memory robust expect plunge thank hen limb electric town slot annual abandon"
  }
}
```

其中secret便是创世账户的助记词，可以导入该助记词到iris客户端以重置密码：

```
iriscli keys add $MYNAME1 --recover
override the existing name XXX [y/n]: y
Enter a passphrase for your key: <your password>
Repeat the passphrase: <your password>
Enter your recovery seed phrase: <seed phrase>
```

1. 启动irishub:

```
iris start --home=$HOME/.iris1
```

1. 多节点配置：

同时，你也可以在本地再运行另一个Node，注意新的Node所对应的工作目录不同，但二者公用一个genesis文件。

```
MYNAME2=<your another name>
iris init --chain-id=test --name=$MYNAME2 --home=$HOME/.iris2
cp $HOME/.iris1/config/genesis.json $HOME/.iris2/config/genesis.json
```

修改 `$HOME/.iris2/config/config.toml` ，防止端口冲突，并让两个Node可以发现彼此：

```
proxy_app = "tcp://127.0.0.1:46668"
moniker = "anonymous"
fast_sync = true
db_backend = "leveldb"
log_level = "main:info,state:info,*:error"

[rpc]
laddr = "tcp://0.0.0.0:46667"

[p2p]
laddr = "tcp://0.0.0.0:46666"
seeds = "<iris1_node_id>@0.0.0.0:46656"
```

其中node_id可以通过以下命令获取：

```
iris tendermint show_node_id --home=$HOME/.iris1
iris tendermint show_node_id --home=$HOME/.iris2
```

然后就可以同时在后台运行：

```
iris start --home=$HOME/.iris1  &> iris1.log & NODE1_PID=$!
iris start --home=$HOME/.iris2  &> iris2.log & NODE2_PID=$!
```

`PID` 是为了方便之后杀死相关进程。

通过 `tail iris1.log`,或 `tail -f iris1.log`可以对Node 进行监控。

1. 现在我们可以使用客户端进行操作

查询本地账户：

```
iriscli keys list
```

查询账户余额：

```
iriscli account <address>
```

我们可以测试在账户间转账：

```
iriscli send --name=$MYNAME1 --to=$MYADDR2 --amount=10steak --chain-id=test
iriscli account $MYADDR1
iriscli account $MYADDR2
```

也可以查询Validator列表:

```
iriscli stake validators
```

接下来实现第一个validator的添加：首先获得节点的公钥：

```
iris tendermint show_validator --home=$HOME/.iris2
```

通过执行以下命令完成绑定：

```
iriscli stake create-validator --amount=10steak --pubkey=<your_node_pubkey> --address-validator=<your_address> --moniker=<your_node_moniker> --chain-id=test --name=<key_name>
```

再次查询 Validator列表:

```
iriscli stake validators
```

通过以下命令为Validator添加更多信息：

```
iriscli stake edit-validator --details="To the IRISnet !" --website="https://irisnet.org"
```

检查Validator是否生效：

```
iriscli advanced tendermint validator-set
```

如果你杀死了第一个Node，那么区块链将无法继续生成新的块。这是因为只有保证Voting Power总数超过 2/3 的节点在线的情况下，才能达成共识。

将第一个节点重新启动，那么区块链会回归正常。