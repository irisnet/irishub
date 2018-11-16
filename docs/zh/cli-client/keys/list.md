# iriscli keys list

## 描述

返回此密钥管理器存储的所有公钥的列表以及他们的相关名称和地址。

## 使用方式

```
iriscli keys list [flags]
```

## 标志

| 名称, 速记       | 默认值     | 描述                                                         | 是否必须  |
| --------------- | --------- | ------------------------------------------------------------ | -------- |
| --help, -h      |           | help for add                                                 |          |

## 例子

### 列出所有的密钥

```shell
iriscli keys list
```

执行命令后，你会得到所有存于本地客户端存储的密钥，包括它们的地址和公钥信息。

```txt
NAME:	TYPE:	ADDRESS:						            PUBKEY:
abc  	local	faa1va2eu9qhwn5fx58kvl87x05ee4qrgh44yd8teh	fap1addwnpepq02r0hts0yjhp4rsal627s2lqk4agy2g6tek5g9yq2tfrmkkehee2td75cs
test	local	faa1kkm4w5pvmcw0e3vjcxqtfxwqpm3k0zakl7lxn5	fap1addwnpepq0gsl90v9dgac3r9hzgz53ul5ml5ynq89ax9x8qs5jgv5z5vyssskww57lw
```

需要注意的是，如果本地有多个.iriscli存储，需要通过--home 参数来定位查询源。