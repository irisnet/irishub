# iris tendermint show-validator

## 介绍

在{irishome}/config获取验证人的私钥文件private_key.json，再推导出公钥，然后以`icp`为hrp进行bech32编码后的结果。

## 用法

```
iris tendermint show-validator [flags]
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
iris tendermint show-validator --home={iris-home}
```

示例输出
```$xslt
icp1zcjduepqzuz420weqehs3mq0qny54umfk5r78yup6twtdt7mxafrprms5zqsjeuxvx
```

返回的结果是bech32编码的地址, 关于bech32的详细文档请参阅 [bech32-prefix](../../features/basic-concepts/bech32-prefix.md) 