# 如何运行一个全节点

## 配置

## 初始化全节点


首先执行以下命令：
```
iris init --moniker=<your_custom_name> --home=$IRISHOME --chain-id=irishub
```

iris在运行过程中所依赖的配置文件和数据会存放在\$IRISHOME下，所以在运行iris前，需要指定一个目录作为\$IRISHOME。\$IRISHOME默认为：/Users/$user/.iris。

在\$IRISHOME需要设置两个文件夹：/config 和 /data



### 下载配置文件

iris运行中需要用到两个重要的文件：genesis.json 和config.toml

genesis文件中定义了区块链网络的初始状态，而config.toml指定了iris软件模块的重要组成部分。

下载这两个文件到/$IRISHOME/config目录下：

```
cd $IRISHOME/config/
rm genesis.json
rm config.toml
wget https://raw.githubusercontent.com/irisnet/betanet/master/config/genesis.json
wget https://raw.githubusercontent.com/irisnet/betanet/master/config/config.toml
```
### 修改配置文件
在config.toml文件中可以配置以下信息：
* 将`moniker`字段配置称为自定义的名称，这样便于区分不同的节点
* `seed`字段用语设置种子节点，在irishub mainnet中的官方种子节点为：
```
6a6de770deaa4b8c061dffd82e9c7f4d40c2165d@seed-1.mainnet.irisnet.org:26656
a17d7923293203c64ba75723db4d5f28e642f469@seed-2.mainnet.irisnet.org:26656
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
{"node_info":{"protocol_version":{"p2p":"5","block":"8","app":"0"},"id":"8fa36b85e98f986b70889da52b733fa925908947","listen_addr":"tcp://0.0.0.0:26656","network":"irishub","version":"0.27.3","channels":"4020212223303800","moniker":"test","other":{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"sync_info":{"latest_block_hash":"DF2F64D56863C5516586112B9A954DFB2257C65FF178267E75D85D160E5E0E2B","latest_app_hash":"","latest_block_height":"1","latest_block_time":"2019-01-23T03:42:17.268038878Z","catching_up":false},"validator_info":{"address":"7B9052D259643E5B9AF0BD481B843C89B27AACAA","pub_key":{"type":"tendermint/PubKeyEd25519","value":"Mt9WvNPDd8F4Dcq7lP+GFIhW0/K4jAt8nTq/ljut94E="},"voting_power":"100"}}
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
{"node_info":{"id":"3fb472c641078eaaee4a4acbe32841f18967672c","listen_addr":"172.31.0.190:26656","network":"irishub","version":"0.22.6","channels":"4020212223303800","moniker":"name","other":["amino_version=0.10.1","p2p_version=0.5.0","consensus_version=v1/0.2.2","rpc_version=0.7.0/3","tx_index=on","rpc_addr=tcp://0.0.0.0:26657"]},"sync_info":{"latest_block_hash":"7B1168B2055B19F811773EEE56BB3C9ECB6F3B37","latest_app_hash":"B8F7F8BF18E3F1829CCDE26897DB905A51AF4372","latest_block_height":12567,"latest_block_time":"2018-08-25T11:33:13.164432273Z","catching_up":false},"validator_info":{"address":"CAF80DAEC0F4A7036DD2116B56F89B07F43A133E","pub_key":{"type":"AC26791624DE60","value":"Cl6Yq+gqZZY14QxrguOaZqAswPhluv7bDfcyQx2uSRc="},"voting_power":100}}
```
通过以上命令可以查看状态：

* `"catching_up":false`: 表示节点与网络保持同步

* `"latest_block_height"`: 表示最新的区块高度
