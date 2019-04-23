# iriscli keys show

## 描述

查询本地密钥的详细信息

## 使用方式

```
iriscli keys show <name> <flags>
```

## 标志

| 名称, 速记            | 默认值             | 描述                                                           | 是否必须  |
| -------------------- | ----------------- | -------------------------------------------------------------- | -------- |
| --address            |                   | 仅输出地址                                                      |          |
| --bech               | acc               | 密钥的Bech32前缀编码 (acc/val/cons)                     |          |
| --help, -h           |                   | 查询命令帮助                                                    |          |
| --multisig-threshold | 1                 | K out of N required signatures                          |          |
| --pubkey             |                   | 仅输出公钥                                                      |          |

## 例子

### 查询指定密钥

```shell
iriscli keys show MyKey
```

执行命令后，你会得到本地客户端存储的指定密钥详情。

```txt
NAME:	TYPE:	ADDRESS:						            PUBKEY:
MyKey	local	iaa1kkm4w5pvmcw0e3vjcxqtfxwqpm3k0zak83e7nf	iap1addwnpepq0gsl90v9dgac3r9hzgz53ul5ml5ynq89ax9x8qs5jgv5z5vyssskzc7exa
```

### 查询验证人（iva/ivp）地址

如果某个地址已绑定成为验证人，则可以使用`--bech val`获取其`iva/ivp`的地址： 

```$xslt
iriscli keys show alice --bech val
```

Then you could see the following:
```$xslt
NAME: TYPE: ADDRESS: PUBKEY:
alice local iva12nda6xwpmp000jghyneazh4kkgl2tnzyx7trze ivp1addwnpepqfw52vyzt9xgshxmw7vgpfqrey30668g36f9z837kj9dy68kn2wxqm8gtmk
```

`iva/ivp`地址可被用于委托命令， 请参阅 [create a delegation](../stake/delegate.md)