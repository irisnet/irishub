# iriscli keys show

## 描述

查询本地密钥的详细信息

## 使用方式

```
iriscli keys show <name> [flags]
```

## 标志

| 名称, 速记            | 默认值             | 描述                                                           | 是否必须  |
| -------------------- | ----------------- | -------------------------------------------------------------- | -------- |
| --address            |                   | 仅输出地址                                                      |          |
| --bech               | acc               | [string] 密钥的Bech32前缀编码 (acc|val|cons)                     |          |
| --help, -h           |                   | 查询命令帮助                                                    |          |
| --multisig-threshold | 1                 | [uint] K out of N required signatures                          |          |
| --pubkey             |                   | 仅输出公钥                                                      |          |

## 例子

### 查询指定密钥

```shell
iriscli keys show MyKey
```

执行命令后，你会得到本地客户端存储的指定密钥详情。

```txt
NAME:	TYPE:	ADDRESS:						            PUBKEY:
MyKey	local	iaa1kkm4w5pvmcw0e3vjcxqtfxwqpm3k0zakl7lxn5	iap1addwnpepq0gsl90v9dgac3r9hzgz53ul5ml5ynq89ax9x8qs5jgv5z5vyssskww57lw
```