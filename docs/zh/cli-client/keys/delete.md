# iriscli keys delete

## 描述

删除指定的密钥

## 使用方式

```
iriscli keys delete <name> [flags]
```

## 标志

| 名称, 速记       | 默认值     | 描述                                                         | 是否必须  |
| --------------- | --------- | ------------------------------------------------------------ | -------- |
| --help, -h      |           | 查询命令帮助                                                  |          |
| --force, -f     |   false   | 无条件地删除密钥，而无需密码 | 否 |
| --yes, -y       |   false   | 删除离线或ledger密钥引用时跳过确认提示 | 否 | 

## 例子

### 删除指定的密钥

```shell
iriscli keys delete MyKey
```

执行命令后，你会收到一个危险警告，并且要求你输入密码用于执行删除指令。

```txt
DANGER - enter password to permanently delete key:
```

输入正确的密码之后，你就完成了删除操作。

```txt
Password deleted forever (uh oh!)
```