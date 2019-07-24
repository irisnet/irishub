# iriscli asset issue

## 描述

实现在 IRIS Hub 链上发行资产。

## 使用方式

```bash
iriscli asset issue-token [flags]
```

## 标志

| 命令，缩写         | 类型    | 是否必须 | 默认值        | 描述                                                         |
| ------------------ | ------- | -------- | ------------- | ------------------------------------------------------------ |
| --family           | string  | true     | fungible      | 资产类型: fungible, non-fungible (暂不支持)                  |
| --source           | string  | false    | native        | 资产源: native, gateway                                      |
| --name             | string  | true     |               | 资产名称, 长度限制在32个unicode字符, e.g. "IRIS Network"     |
| --gateway          | string  | false    |               | 网关的唯一标识, 当 source 为 gateway 时必填                  |
| --symbol           | string  | true     |               | Source内全局唯一的资产标识符，字母和数字的组合，首字符必须为字母，长度[3,8]，大小写无关 |
| --canonical-symbol | string  | false    |               | 当 source 为 gateway 时，用来指定源链上的 Symbol |
| --min-unit-alias   | string  | false    |               | 资产最小单位别名，字母和数字的组合，首字符必须为字母，长度[3,10]，大小写无关 |
| --initial-supply   | uint64  | true     |               | 以基准单位(Symbol)计的初始发行量；合法取值范围：[0,100 billion] |
| --max-supply       | uint64  | false    | 1000000000000 | 以基准单位(Symbol)计的最大发行量；合法取值范围：[初始发行量,1000 billion] |
| --decimal          | uint8   | true     |               | 合法取值范围：[0,18]                                         |
| --mintable         | boolean | false    | false         | 初始发行后是否允许增发(mint)                                 |

## 示例

### 发行原生资产

```bash
iriscli asset issue-token --family=fungible --source=native --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true --fee=1iris --from=<key-name> --commit
```

### 发行网关资产

#### 创建网关

在本例中，必须先创建名为 `cats` 的网关, [更多详情](./create-gateway.md)

```bash
iriscli asset create-gateway --moniker=cats --identity=<identity> --details=<details> --website=<website> --from=<key-name> --commit
```

#### 发行网关资产

```bash
iriscli asset issue-token --family=fungible --source=gateway --gateway=cats --canonical-symbol=cat --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true  --fee=1iris --from=<key-name> --commit
```

### 转账

你可以像[转账iris](../bank/send)一样，转账你所拥有的任何代币

**转账原生资产**

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10kitty --fee=0.3iris --chain-id=irishub
```

**转账网关资产**

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10cats.kitty --fee=0.3iris --chain-id=irishub
```
