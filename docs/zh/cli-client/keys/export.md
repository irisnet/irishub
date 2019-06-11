# iriscli keys export

## 描述

通过json文件的形式备份单个账户的密钥信息。

## 使用方式

```shell
iriscli keys export Mykey --output-file=<path_to_backup_keystore>
```


## 标志

| 名称, 速记       | 默认值     | 描述                                                              | 是否必须  |
| --------------- | --------- | ----------------------------------------------------------------- | -------- |
| --output-file    |           | 导出keystore的地址                                              |          |
| --help, -h      |           | 查询命令帮助                                                       |          |


## 例子

### 备份密钥

该命令会导出一个加密的json文件，并要求指定密码。
```shell
iriscli keys export Mykey --output-file=<path_to_backup_keystore>
```

### 导入密钥

使用备份时指定的密码,导入key。
```shell
iriscli keys add Mykey --recover --keystore=<path_to_backup_keystore>
```
