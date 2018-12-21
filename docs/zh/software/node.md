# IRIS Daemon

## 介绍

iris可执行文件是运行IRISnet网络节点的入口，包括验证人节点和其他全节点都需要通过安装iris命令并启动守护进程来加入到IRISnet网络。也可以使用该命令在本地启动自己的测试网络，如需加入IRISnet测试网请参阅[get-started](../get-started/README.md)。

## 如何在本地启动一个IRISnet网络

### 初始化节点

首先需要为自己创建对应的验证人账户
```bash
iriscli keys add {account_name}
```
得到账户信息，包括账户地址、公钥地址、助记词
```
NAME:	TYPE:	ADDRESS:						PUBKEY:
account_name	local	faa13t6jugwm5uu3h835s5d4zggkklz6rpns59keju	fap1addwnpepqdne60eyssj2plrsusd8049cs5hhhl5alcxv2xu0xmzlhphy9lyd5kpsyzu
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

witness exotic fantasy gaze brass zebra adapt guess drip quote space payment farm argue pear actress garage smile hawk bid bag screen wonder person
```

初始化genesis.json和config.toml等配置文件
```bash
iris init --home={path_to_your_home} --chain-id={your_chain_id} --moniker={your node name}
```
该命令会在home目录下创建相应文件

创建申请成为验证人的交易，并使用刚才创建的验证人账户对交易进行签名
```bash
iris gentx --name={account_name} --home={path_to_your_home}
```
生成好的交易数据存放在目录：{path_to_your_home}/config/gentx

### 配置genesis

使用下面命令修改genesis.json文件，为上述验证人账户分配初始账户余额，如：150个iris
```bash
iris add-genesis-account faa13t6jugwm5uu3h835s5d4zggkklz6rpns59keju 150iris
```

```json
    {
      "accounts": [
        {
          "address": "faa13t6jugwm5uu3h835s5d4zggkklz6rpns59keju",
          "coins": ["150iris"],
          "sequence_number": "0",
          "account_number": "0"
        }
      ]
    }
```

配置验证人信息
```bash
iris collect-gentxs --home={path_to_your_home}
```
该命令读取 {path_to_your_home}/config/gentx 下的 `CreateValidator` 交易，并将其写入到genesis.json中，完成在创世块中对初始验证人的指定。

### 启动节点

完成上述配置后，通过以下命令启动全节点。 
```bash
iris start --home {path_to_your_home}
```
该命令执行后，iris使用home下面所配置的genesis.json来生成创世区块并更新应用状态，后续根据chain-id所指链的当前区块高度，要么开始和其他peer同步区块，要么进入等待首次出块的流程。

## 启动多节点网络

如需在本地启动一个多节点的IRISnet网络，需要按下列步骤操作：

* 准备home目录：为每个节点创建各自独有的home目录
* 初始化：按照上述 `初始化节点` 步骤，在各自的home目录下为各节点完成初始化（注意需要使用同一个chain-id）
* 配置genesis：选择其中一个节点的home目录来配置genesis，参考`配置genesis`步骤，把各个节点的账户地址和初始余额配置进来，再把各个节点{path_to_your_home}/config/gentx下的文件拷贝到当前节点的{path_to_your_home}/config/gentx目录下，然后执行`iris collect-gentxs`来生成最终genesis.json，最后把该genesis.json拷贝到各个节点的home目录下，覆盖原有genesis.json
* 配置config.toml：修改各个节点home目录下的{path_to_your_home}/config/config.toml，为各个节点分配不同的端口，通过`iris tendermint show-node-id`查询各个节点的node-id，然后在`persistent_peers`中加入其它节点的`node-id@ip:port`，使得节点之间能够相互连接

为了简化上面的配置流程，可以使用下面命令来自动完成本地多节点的初始化和genesis、config的配置：
```bash
iris testnet --v 4 --output-dir ./output --chain-id irishub-test --starting-ip-address 127.0.0.1
```

启动节点加入公共测试网络Testnet的方法可参阅[Full-Node](../get-started/Full-Node.md)

## home目录介绍

home目录为iris节点的工作目录，home目录下包含了所有的配置信息和节点运行的所有数据。

在iris命令中可以通过flag `--home` 来指定该节点的home目录，如果在同一台机器上运行多个节点，则需要为他们指定不同的home目录。如果在iris命令中没有指定`--home` flag，则使用默认值 `$HOME/.iris` 作为本次iris命令所使用的home目录。

`iris init` 命令负责对所指定的`--home`目录进行初始化，创建默认的配置文件。除了`iris init` 命令外，其他任何`iris`相关命令所使用的home目录必须是被初始化过的，否则将会报错。

home的`data`目录下存放iris节点运行的数据，包括区块链数据、应用层数据、索引数据等。home的`config`目录下存放所有配置文件：

### genesis.json

genesis.json为创世块数据，指定了chain_id、共识参数、初始账户代币分配、创建验证人、stake/slashing/gov/upgrade等系统参数。详见[genesis-file](../features/basic-concepts/genesis-file.md)

### node_key.json

node_key.json用来存放节点的密钥，通过`iris tendermint show-node-id`查询到的node-id就是通过该密钥计算得来的，用来表示节点的唯一身份，在p2p连接中使用。

### pri_validator.json

pri_validator.json为验证人在每一轮出块共识投票中对Pre-vote/Pre-commit等进行签名的密钥，随着共识的进行，tendermint共识引擎会不断更新`last_height`/`last_round`/`last_step`等数据。

### config.toml

config.toml为该节点的非共识的配置信息，不同节点可根据自己的情况自行配置。常见修改项有`persistent_peers`/`moniker`/`laddr`等。
