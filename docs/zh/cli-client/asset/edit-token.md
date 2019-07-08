# iriscli asset edit-token

## 描述

编辑指定ID的资产信息

## 使用方式

```bash
iriscli asset edit-token <token-id> --name=<name> --symbol-at-source=<symbol-at-source> --symbol-min-alias=<min-alias> --max-supply=<max-supply> --mintable=<mintable> --from=<your account name> --chain-id=<chain-id> --fee=0.6iris
```

## 特有的标志

| Name               | Type   | Required | Default | Description                                                  |
| ------------------ | ------ | -------- | ------- | ------------------------------------------------------------ |
| --name             | string | 否       | ""      | 资产名称，例如：IRIS Network                                 |
| --symbol-at-source | string | 否       | ""      | Source为 external 或 gateway 的时候，可以用来指定在源链上的Symbol |
| --symbol-min-alias | string | 否       | ""      | 资产最小单位别名                                             |
| --max-supply       | uint   | 否       | 0       | 以基准单位计的最大发行量                                     |
| --mintable         | bool   | 否       | false   | 初始发行后是否允许增发                                       |

## 示例

`max-supply` 只能减少，不能增加，且不能低于当前Token总量

```bash
iriscli asset edit-token cat --name="Cat Token" --symbol-at-source="cat" --symbol-min-alias=kitty --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=irishub --fee=0.4iris --commit
```
