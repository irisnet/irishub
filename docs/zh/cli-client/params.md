# iriscli params

Params模块允许查询系统里预设的参数，查询结果中除了Gov模块的参数，其他都可以通过[Gov模块](./gov.md)发起提议来修改。

## 可用命令

| 名称                                       | 描述                                                                            |
| ------------------------------------------ | -------------------------------------------------------------------------------------- |
| [subspace](#iris-query-params-subspace)    | 根据subspace和key查询原始参数                                   |

**标志:**

| 名称，速记 | 默认 | 描述     | 必须 |
| ---------- | ---- | -------- | ---- |
| --module   |      | 模块名称 |      |

## iris query params subspace

该命令用于根据subspace和key查询原始参数

```bash
iris query params subspace [subspace] [key] [flags]
```

**标识：**

| 名称, 速记    | 类型    | 必须    | 默认   | 描述                 |
| -------------| ------ | ------ | ----- | ------------------------- |
| --help ,-h   |        |        |       | params帮助信息           |

### 根据subspace和key查询原始参数

```bash
iris query params subspace [subspace] [key] [flags]
```

