# iriscli keys update

## 描述

更改用于保护私钥的密码

## 使用方式

```
iriscli keys update <name> [flags]
```

## 标志

| 名称, 速记       | 默认值     | 描述                                                         | 是否必须  |
| --------------- | --------- | ------------------------------------------------------------ | -------- |
| --help, -h      |           | help for add                                                 |          |

## 例子

### 修改密码

```shell
iriscli keys update MyKey
```

执行命令后，系统会要求你输入指定密钥的当前密码。

```txt
Enter the current passphrase:
```

然后系统会要求你输入新密码并重复输入确认密码。

```txt
Enter the new passphrase:
Repeat the new passphrase:
```

如果你输入新密码符合要求，则会执行更新操作。

```txt
Password successfully updated!
```