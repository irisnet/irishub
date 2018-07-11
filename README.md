# Iris-Hub
Iris Hub是在Cosmos生态中的区域性枢纽，提供iService服务

## 下载发行版安装
 * 进入下载页: https://github.com/irisnet/iris-hub/releases/tag/0.1.0
 * 下载对应版本的可执行文件
 * 解压缩`tar -C /usr/local -xzf iris$VERSION.$OS-$ARCH.zip`
 * 拷贝到`/usr/local/`目录下
执行
 ```
iris version
 ```

若出现`v 0.1.0`则说明安装成功



## 通过源码编译安装

### 系统要求
操作系统：**Ubuntu 16.04** 

架构：**amd64**

要求安装1.9版本以上的Go https://golang.org/doc/install

GO 安装步骤

* 下载
curl -O https://dl.google.com/go/go1.9.4.linux-386.tar.gz

* 解压缩 
```
tar xvf go1.9.4.linux-386.tar.gz
```

* 获得权限
```
sudo chown -R root:root ./go

sudo mv go /usr/local
```

* 设置GOPATH
```
export GOPATH=$HOME/work
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

* 测试是否安装成功
建立一个简单的工程。

```
mkdir -p work/src/github.com/user/hello
```

创建 "Hello World" Go文件.


```
nano ~/work/src/github.com/user/hello/hello.go
```
修改文件内容如下：

```
package main

import "fmt"

func main() {
    fmt.Printf("hello, world\n")
}
```
编译&安装
```
go install github.com/user/hello
```
执行命令
```
hello
```
如果安装成功，将会出现 "hello, world"输出：



## 安装步骤
```
go get github.com/irisnet/iris-hub
cd $GOPATH/src/github.com/irisnet/iris-hub
make all
iris version
```

若出现`v 0.1.0`则说明安装成功

###本地测试网络

Here is a quick example to get you off your feet: 

1. 生成和保存测试网络中的测试账户

```
MYNAME=<your name>
iris client keys new $MYNAME
iris client keys list
MYADDR=<your newly generated address>
```
2. 初始化iris-hub:

```
iris node init $MYADDR --home=$HOME/.iris1 --chain-id=test 
```

This will create all the files necessary to run a single node chain in以上命令将在 `$HOME/.iris1／config`目录下创建默认的配置文件 `priv_validator.json`和`genesis.json`。`priv_validator.json` 包含验证人的私钥。`genesis.json`文件包含其他配置信息， `$MYADDR`是默认生成的创世账户。

获得节点的Node_ID

需要下载编辑develop分支的tendermint，或者下载可执行[[文件](https://github.com/kidinamoto01/gaia-testnet/blob/master/gaia-testnet/tendermint-develop-liunx-amd64.zip)]

执行以下命令获得iris1的Node_ID
```
tendermint show_node_id --home=$HOME/.iris1
```

同时，你也可以在本地再运行另一个Node，注意新的Node所对应的工作目录不同，但二者公用一个genesis文件。


```
iris node init $MYADDR --home=$HOME/.iris2 --chain-id=test
cp $HOME/.iris1/config/genesis.json $HOME/.iris2/config/genesis.json
```

修改 `$HOME/.iris2/config/config.toml` 让两个Node可以发现彼此：

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

然后就可以同时在后台运行：

```
iris node start --home=$HOME/.iris1  &> iris1.log &
NODE1_PID=$!
iris node start --home=$HOME/.iris2  &> iris2.log &
NODE2_PID=$!
```

`PID` 是为了方便之后杀死相关进程。

通过 `tail iris1.log`,或 `tail -f iris1.log`可以对Node 进行监控。

现在我们初始化客户端，然后查询账户：

```
iris client init --chain-id=test --node=tcp://localhost:46657
iris client query account $MYADDR
```

我们可以测试在账户间转账：

```
MYNAME1=<your name>
iris client keys new $MYNAME1
iris client keys list
MYADDR1=<your newly generated address>
iris client tx send --amount=1000fermion --to=$MYADDR1 --name=$MYNAME
iris client query $MYADDR1
iris client query $MYADDR
```

也可以查询 candidate列表:

```
iris client query candidates
```

很奇怪的是，iris并没有发现initial时生成的validator。这是因为这类validator比较特殊。通过查询tendermint接口，你可以获得初始的validator信息。

```
curl localhost:46657/validators
```

接下来实现第一个validator的添加：首先获得节点的公钥：

```
cat $HOME/.iris1/config/priv_validator.json 
```

使用`jq`命令可以快速获得公钥:

```
cat $HOME/.iris1/config/priv_validator.json | jq .pub_key.data

```

通过执行以下命令完成绑定
```
iris client tx declare-candidacy --amount=1fermion --pubkey=<validator pubkey> --moniker=<moniker>  --name=$MYNAME

```
再次查询 candidate列表:

```
iris client query candidates
```
也可以查询具体的validator信息

```
iris client query candidate --pubkey=<validator pubkey>
```

查询初始账户的变化，candidates页有变化：

```
iris client query account $MYADDR
iris client query candidates
```

通过以下命令可以修改绑定的金额
```
iris client tx edit-candidacy --pubkey=<validator pubkey> --moniker=<new moniker> --name=$MYNAME
```
通过unbond可将抵押的代币收回，你会发现账户余额增加了。

```
iris client tx unbond --shares=1 --pubkey=<validator pubkey> --name=$MYNAME
iris client query validators
iris client query account $MYADDR
```
同样的，你也可以绑定用$HOME/.iris2/config/priv_validator.json代表的第二个Node
通过查询tendermint接口，你可以获得新的validator信息。


```
curl localhost:46657/validators
```

如果你杀死了第二个Node，那么区块链将无法继续生成新的块。这是因为只有保证2/3节点在线的情况下，才能达成共识。

将第二个节点重新启动，那么区块链会回归正常。