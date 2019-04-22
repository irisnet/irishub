# iris tendermint 和 iriscli tendermint 的子命令

## iriscli tendermint 介绍

tendermint的状态查询的子命令

### 用法
```
iriscli tendermint <subcommand>
```

### 相关子命令

|命令|描述|
|---|---|
|[tx](tx.md)|在所有提交的块上匹配此txhash|
|[txs](txs.md)|搜索与给定标签匹配的所有事务|
|[block](block.md)|获取给定高度的块的验证数据|
|[validator-set](validator-set.md)|在指定高度查询验证人集合|


## iris tendermint 介绍

### 介绍

获取全节点相应的信息

### 用法

```shell
iris tendermint <subcommand>
```

### 相关子命令

| Name                    | Description                                                                                  |
| ----------------------- | -------------------------------------------------------------------------------------------- |
| [show-node-id](show-node-id.md) | 获取全节点p2p网络的id |
| [show-validator](show-validator.md) | 获取验证人公钥 |
| [show-address](show-address.md)  |    获取验证人地址 |


### 全局标志

| 名称，速记 | 默认值        | 功能描述                            | 是否必须 |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | 编码方式 (hex\b64\btc) |          |
| --home          | $HOME/.iris    | 存放运行数据和配置文件的目录 |   |
| --output, -o    | text           | 输出格式 (text,json)     |   |
| --trace         |                | 是否打印callstack和所有错误信息   |    |