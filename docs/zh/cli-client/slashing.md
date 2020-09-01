# Slashing
Slashing 模块可以解禁被监禁的验证人

## 可用命令

| 名称                                                | 描述                         |
| --------------------------------------------------- | ---------------------------- |
| [unjail](#iris-tx-slashing-unjail)                  | 解禁被监禁的验证人           |
| [params](#iris-query-slashing-params)               | 查询当前`Slashing`的参数信息 |
| [signing-info](#iris-query-slashing-signing-info)   | 查询验证人的签名信息         |
| [signing-infos](#iris-query-slashing-signing-infos) | 查询所有验证人的签名信息     |

## iris tx slashing unjail

解禁被监禁的验证人。

```bash
iris tx slashing unjail [flags]
```

## iris query slashing params

查询当前`Slashing`的参数信息。

```bash
iris query slashing params [flags]
```

## iris query slashing signing-info

查询验证人的签名信息。

```bash
iris query slashing signing-info [validator-conspub] [flags]
```

## iris query slashing signing-infos

查询所有验证人的签名信息。

```bash
iris query slashing signing-infos [flags]
```