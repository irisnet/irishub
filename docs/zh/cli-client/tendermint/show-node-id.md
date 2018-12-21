# iris tendermint show-node-id

## 介绍

获取节点p2p的id

## 用法

```
iris tendermint show-node-id [flags]
```

## 标志

| 名称，速记      | 默认值           | 介绍                                                    | 是否必填 |
| -------------------- | ----------------- | -------------------------------------------------------------- | -------- |
| --help, -h           |                   | 打印帮助信息                                                 |  否        |

### 全局标志

| 名称，速记 | 默认值        | 功能描述                            | 是否必须 |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | [string] 编码方式 (hex|b64|btc) |          |
| --home          | $HOME/.iris    | [string] 存放运行数据和配置文件的目录 |   |
| --output, -o    | text           | [string] 输出格式 (text,json)     |   |
| --trace         |                | 是否打印callstack和所有错误信息   |    |

## 示例

```shell
iris tendermint show-node-id --home={iris-home}
```

示例返回
```$xslt
b18d3d1990c886555241f91331f9c00fe69421aa
```

