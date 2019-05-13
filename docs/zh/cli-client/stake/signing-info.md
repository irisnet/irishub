# iriscli stake signing-info

## 描述

查询验证者签名信息

## 用法

```
iriscli stake signing-info <validator-pubkey> <flags>
```

打印帮助信息
```
iriscli stake signing-info --help
```

## 示例

查询验证者签名信息
```
iriscli stake signing-info <validator-pubkey>
```

运行成功以后，返回的结果如下：
```txt
  Signing Info
  Start Height:          0
  Index Offset:          3506
  Jailed Until:          1970-01-01 00:00:00 +0000 UTC
  Missed Blocks Counter: 0
```
