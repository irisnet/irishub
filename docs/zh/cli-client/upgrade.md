# iriscli upgrade

此命令用于查询软件升级状态

## 可用命令

| 名称                                            | 描述               |
| ----------------------------------------------- | ------------------ |
| [info](#iriscli-upgrade-info)                   | 查询升级模块的信息 |
| [query-signals](#iriscli-upgrade-query-signals) | 查询signals的信息  |

## iriscli upgrade info

### 查询链上正在使用的版本

```bash
iriscli upgrade info
```

这将显示当前的协议信息以及准备升级的协议信息，例如

```json
{
  "version": {
    "ProposalID": "1",
    "Success": true,
    "Protocol": {
      "version": "0",
      "software": "https://github.com/irisnet/irishub/tree/v0.10.0",
      "height": "1"
    }
  },
  "upgrade_config": {
    "ProposalID": "3",
    "Definition": {
      "version": "1",
      "software": "https://github.com/irisnet/irishub/tree/v0.10.1",
      "height": "8000"
    }
  }
}
```

## iriscli upgrade query-signals

查询当前的signals信息

**标志：**

| 名称，速记 | 默认 | 描述        | 必须 |
| --------------- | ------- | ------------------ | -------- |
| --detail        | false   | signals详情 |          |

### 查询已升级的voting power统计信息

```bash
iriscli upgrade query-signals
```

示例输出：

```bash
signalsVotingPower/totalVotingPower = 0.5000000000
```

### 查询升级signals详情

```bash
iriscli upgrade query-signals --detail
```

示例输出：

```bash
iva15cv33a67cfey5eze7238hck6yngw36949evplx   100.0000000000
iva15cv33a67cfey5eze7238hck6yngw36949evplx   100.0000000000
iva15cv33a67cfey5eze7238hck6yngw36949evplx   100.0000000000
siganalsVotingPower/totalVotingPower = 0.5000000000
```
