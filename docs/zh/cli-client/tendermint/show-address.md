# iris tendermint show-address

## 介绍

查询validator的公钥对于的地址

## 用法

```
iris tendermint show-address [flags]
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

### 验证人地址

```
iris tendermint show-address --home={iris-home}
```

示例返回：
``·
fca1jhax359kz6hm70swxpxumpgkaglr7yq80qhwvv
```

返回的结果是bech32编码的地址, 关于bech32的详细文档请参阅 [bech32-prefix](../../features/basic-concepts/bech32-prefix.md)
