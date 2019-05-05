# iris tendermint show-node-id

## 介绍

从<path_to_your_home>/config获取节点的私钥文件node_key.json，再推导公钥对应的地址，然后进行hex编码的结果。

## 用法

```
iris tendermint show-node-id <flags>
```

### 全局标志

| 名称，速记 | 默认值        | 功能描述                            | 是否必须 |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | 编码方式 (hex/b64/btc) |          |
| --home          | $HOME/.iris    | 存放运行数据和配置文件的目录 |   |
| --output, -o    | text           | 输出格式 (text,json)     |   |
| --trace         |                | 是否打印callstack和所有错误信息   |    |

## 示例

```shell
iris tendermint show-node-id --home=<path_to_your_home>
```

示例返回
```
b18d3d1990c886555241f91331f9c00fe69421aa
```

