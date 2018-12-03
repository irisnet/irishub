# 如何运行一个全节点

## 配置

### 设置软件运行的目录

iris在运行过程中所依赖的配置文件和数据会存放在\$IRISHOME下，所以在运行iris前，需要指定一个目录作为\$IRISHOME。\$IRISHOME默认为：/Users/$user/.iris。

在\$IRISHOME需要设置两个文件夹：/config 和 /data

### 下载配置文件文件
iris运行中需要用到两个重要的文件：genesis.json 和config.toml

genesis文件中定义了区块链网络的初始状态，而config.toml指定了iris软件模块的重要组成部分。

下载这两个文件到/$IRISHOME/config目录下：

```
cd $IRISHOME/config/
wget https://raw.githubusercontent.com/irisnet/testnets/master/fuxi/fuxi-5000/config/config.toml
wget https://raw.githubusercontent.com/irisnet/testnets/master/fuxi/fuxi-5000/config/genesis.json
```
### 修改配置文件
在config.toml文件中可以配置以下信息：
* 将`moniker`字段配置称为自定义的名称，这样便于区分不同的节点
* `seed`字段用语设置种子节点，在fuxi-5000中的官方种子节点为：
```
c16700520a810b270206d59f0f02ea9abd85a4fe@ts-1.bianjie.ai:26656
a12cfb2f535210ea12731f94a76b691832056156@ts-2.bianjie.ai:26656
```

你也可以配置 `moniker` 和 `external_address` 字段. 

```
moniker = "<your_custom_name>"
external_address = "your-public-IP:26656"
```


另外，如果你需要与其他节点通过内网链接，请设置 `addr_book_strict` 为 `false` 。

```
addr_book_strict = false
```
###  配置端口

如果你的节点需要与其他节点建立链接，则需要开放 `26656` 端口；若需要通过rpc端口查询Tendermint提供的信息，则需要开放 `26657` 端口。

通过以下命令启动全节点，并将日志输出到文件中：
```
iris start --home {path_to_your_home} > log文件地址 &
```
通过执行以下操作确认节点的运行状态：
```
iriscli status
```
示例输出：
```json
{"node_info":{"id":"3fb472c641078eaaee4a4acbe32841f18967672c","listen_addr":"172.31.0.190:26656","network":"fuxi-5000","version":"0.22.6","channels":"4020212223303800","moniker":"name","other":["amino_version=0.10.1","p2p_version=0.5.0","consensus_version=v1/0.2.2","rpc_version=0.7.0/3","tx_index=on","rpc_addr=tcp://0.0.0.0:26657"]},"sync_info":{"latest_block_hash":"7B1168B2055B19F811773EEE56BB3C9ECB6F3B37","latest_app_hash":"B8F7F8BF18E3F1829CCDE26897DB905A51AF4372","latest_block_height":12567,"latest_block_time":"2018-08-25T11:33:13.164432273Z","catching_up":false},"validator_info":{"address":"CAF80DAEC0F4A7036DD2116B56F89B07F43A133E","pub_key":{"type":"AC26791624DE60","value":"Cl6Yq+gqZZY14QxrguOaZqAswPhluv7bDfcyQx2uSRc="},"voting_power":0}}
```
通过以上命令可以查看状态：

* `"catching_up":false`: 表示节点与网络保持同步

* `"catching_up":true`: 表示节点正在同步区块

* `"latest_block_height"`: 表示最新的区块高度


之后你就应该可以在浏览器中看到

## 重置一个全节点

若需要将一个节点重启，则可以通过以下命令让节点再次通过与网络保持同步。

### 重置IRIShub节点流程如下：

1. 关闭iris进程
```
kill -9 <PID>
```

若Genesis文件有变动，则需要下载新的文件到$IRISHOME/config目录下。

2. 重置iris
```
iris unsafe-reset-all --home=<home>
```

3. 重新启动

通过以下命令启动全节点，并将日志输出到文件中：
```
iris start --home <path_to_your_home> > log文件地址 &
```
通过执行以下操作确认节点的运行状态：
```
iriscli status
```
示例输出：
```json
{"node_info":{"id":"3fb472c641078eaaee4a4acbe32841f18967672c","listen_addr":"172.31.0.190:26656","network":"fuxi-5000","version":"0.22.6","channels":"4020212223303800","moniker":"name","other":["amino_version=0.10.1","p2p_version=0.5.0","consensus_version=v1/0.2.2","rpc_version=0.7.0/3","tx_index=on","rpc_addr=tcp://0.0.0.0:26657"]},"sync_info":{"latest_block_hash":"7B1168B2055B19F811773EEE56BB3C9ECB6F3B37","latest_app_hash":"B8F7F8BF18E3F1829CCDE26897DB905A51AF4372","latest_block_height":12567,"latest_block_time":"2018-08-25T11:33:13.164432273Z","catching_up":false},"validator_info":{"address":"CAF80DAEC0F4A7036DD2116B56F89B07F43A133E","pub_key":{"type":"AC26791624DE60","value":"Cl6Yq+gqZZY14QxrguOaZqAswPhluv7bDfcyQx2uSRc="},"voting_power":100}}
```
通过以上命令可以查看状态：

* `"catching_up":false`: 表示节点与网络保持同步

* `"latest_block_height"`: 表示最新的区块高度
