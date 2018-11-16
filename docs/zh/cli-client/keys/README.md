# iriscli keys

## 描述

Keys模块用于管理本地密钥库。

## 使用方式

```shell
iriscli keys [command]
```

## 相关命令

| 命令                    | 描述                                                                                          |
| ----------------------- | -------------------------------------------------------------------------------------------- |
| [mnemonic](mnemonic.md) | 通过读取系统熵来创建bip39助记词，也可以称为种子短语。                                               |
| [new](new.md)           | 使用交互式命令派生新的私钥，该命令将提示你每个输入。                                                 |
| [add](add.md)           | 创建新密钥，或从助记词导入已有密钥                                                                |
| [list](list.md)         | 列出所有密钥                                                                                   |
| [show](show.md)         | 显示指定名称的关键信息                                                                           |
| [delete](delete.md)     | 删除指定的密钥                                                                                  |
| [update](update.md)     | 更改用于保护私钥的密码                                                                           |

## 标志

| 名称, 速记       | 默认值   | 描述          | 是否必须  |
| --------------- | ------- | ------------- | -------- |
| --help, -h      |         | help for keys |          |

## 全局标志

| 名称, 速记       | 默认值          | 描述                                   | 是否必须  |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | [string] Binary encoding (hex|b64|btc) |          |
| --home          | $HOME/.iriscli | [string] directory for config and data |          |
| --output, -o    | text           | [string] Output format (text|json)     |          |
| --trace         |                | print out full stack trace on errors   |          |

## 补充说明

这些密钥可以是go-crypto支持的任何格式，并且可以由轻客户端，完整节点或需要使用私钥签名的任何其他应用程序使用。
