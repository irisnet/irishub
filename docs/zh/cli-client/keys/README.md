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
| [mnemonic](mnemonic.md) | Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. |
| [new](new.md)           | Derive a new private key using an interactive command that will prompt you for each input.   |
| [add](add.md)           | Create a new key, or import from seed                                                        |
| [list](list.md)         | List all keys                                                                                |
| [show](show.md)         | Show key info for the given name                                                             |
| [delete](delete.md)     | Delete the given key                                                                         |
| [update](update.md)     | Change the password used to protect private key                                              |

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
