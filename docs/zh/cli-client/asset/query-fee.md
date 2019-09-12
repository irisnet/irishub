# iriscli asset query-fee

## 描述

查询资产相关的费用，包括网关创建费用以及Token发行和增发费用

## 使用方式

```bash
iriscli asset query-fee [flags]
```

## 特有的标志

| 命令, 速记 | 类型    | 是否必须 | 默认值 | 描述                               |
| --------- | ------ | ------- | ---- | ---------------------------------- |
| --gateway | string |         |      | 网关名字; 如果查询网关费用,则必须指定    |
| --token   | string |         |      | Token ID; 如果查询Token费用,则必须指定 |

## 示例

### 查询网关创建费用

```bash
iriscli asset query-fee --gateway=cats
```

### 查询发行/增发原生资产费用

```bash
iriscli asset query-fee --token=kitty
```

### 查询发行/增发网关资产费用

```bash
iriscli asset query-fee --token=cats.kitty
```
