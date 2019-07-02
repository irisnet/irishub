# iriscli asset query-tokens

## 描述

根据条件查询 IRIS Hub 链上发行的资产集合。

## 使用方式

```bash
iriscli asset query-tokens [flags]
```

## 标志

| 命令，缩写         | 类型    | 是否必须 | 默认值        | 描述                                                         |
| ------------------ | ------- | -------- | ------------- | ------------------------------------------------------------ |
| --source           | string  | false    | all           | 资产源: native / gateway / external                           |
| --gateway          | string  | false    |               | 网关的唯一标识，当 source 为 gateway 时必填                  |
| --owner            | string  | false    |               | 资产所有人                   |

## 查询规则

- 当 source 为 native 时
    - gateway 会被忽略
    - owner 可选
- 当 source 为 gateway 时
    - gateway 必填
    - owner 会被忽略（因为 gateway tokens 全部属于 gateway 的 owner ）
- 当 source 为 external 时
    - gateway 和 owner 都会被忽略
- 当 gateway 不为空时
    - source 可选    
    
## 示例

### 查询全部资产

```bash
iriscli asset query-tokens
```

### 查询全部 native 资产

```bash
iriscli asset query-tokens --source=native
```

### 查询名为 "cats" 的网关的全部资产

```bash
iriscli asset query-tokens --gateway=cats
```

### 查询指定 owner 的全部资产

```bash
iriscli asset query-tokens --owner=<address>
```