# Upgrade

此模块提供软件版本升级的基本功能。

## 可用命令

| 名称                                                                            | 描述                           |
| ------------------------------------------------------------------------------- | ------------------------------ |
| [software-upgrade](#iris-tx-gov-submit-proposal-software-upgrade)               | 发起软件升级提案               |
| [cancel-software-upgrade](#iris-tx-gov-submit-proposal-cancel-software-upgrade) | 取消当前升级提案               |
| [plan](#iris-query-upgrade-plan)                                                | 查询当前正在进行的软件升级提案 |
| [applied](#iris-query-upgrade-applied)                                          | 查询已经执行的软件升级提案     |

## iris tx gov submit-proposal software-upgrade

通过治理模块发起一个软件升级提案。

```bash
iris tx gov submit-proposal software-upgrade <plan-name> [flags]
```

**标识：**

| 名称，速记       | 类型   | 必须 | 默认 | 描述                                             |
| ---------------- | ------ | ---- | ---- | ------------------------------------------------ |
| --deposit        | Coins  |      |      | 为提案抵押的代币数量                             |
| --title          | string |      |      | 提案标题                                         |
| --description    | string |      |      | 提案描述信息                                     |
| --upgrade-height | uint64 |      |      | 执行软件升级的高度（和`--upgrade-time`二选一）   |
| --upgrade-time   | Time   |      |      | 执行软件升级的时间（例如：2006-01-02T15:04:05Z） |
| --upgrade-info   | string |      |      | 软件升级信息                                     |

:::tip
如果需要支持[cosmovisor](#https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor)自动执行软件升级，`--upgrade-info`需要使用固定的格式，例如：

```json
{
    "binaries": {
        "linux/amd64":"https://example.com/irishub.zip?checksum=sha256:aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f"
    }
}
```

:::

## iris tx gov submit-proposal cancel-software-upgrade

通过治理模块发起取消当前正在进行的软件升级提案。

```bash
iris tx gov submit-proposal cancel-software-upgrade [flags]
```

**标识：**

| 名称，速记    | 类型   | 必须 | 默认 | 描述                 |
| ------------- | ------ | ---- | ---- | -------------------- |
| --deposit     | Coins  |      |      | 为提案抵押的代币数量 |
| --title       | string |      |      | 提案标题             |
| --description | string |      |      | 提案描述信息         |

## iris query upgrade plan

查询当前正在进行的软件升级计划。

```bash
iris query upgrade plan [flags]
```

## iris query upgrade applied

查询已近执行的软件升级计划

```bash
iris query upgrade applied <upgrade-name>
```
