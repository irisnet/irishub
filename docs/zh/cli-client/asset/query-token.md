# iriscli asset query-token

## 描述

查询 IRIS Hub 链上发行的资产。

## 使用方式

```bash
iriscli asset query-token <token-id>
```

### 全局唯一 Token ID 生成规则
    
- Source 为 native 时：ID = [Symbol]，例：iris
    
- Source 为 external 时：ID = x.[Symbol]，例：x.btc
    
- Source 为 gateway 时：ID = [Gateway].[Symbol]，例：cats.kitty

## 示例

### 查询名为 "kitty" 的 native 资产

```bash
iriscli asset query-token kitty
```

### 查询名为 "kitty" 的 cats 网关资产

```bash
iriscli asset query-token cats.kitty
```

### 查询名为 "btc" 的 external 资产

```bash
iriscli asset query-token x.btc
```
