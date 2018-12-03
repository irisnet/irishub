# iriscli upgrade submit-switch

## 描述

安装完新软件后，向这次升级相关的提议发送switch消息，表明自己已经安装新软件并把消息广播到全网。

## 用例

```
iriscli upgrade submit-switch [flags]
```
打印帮助信息:

```
iriscli upgrade submit-switch --help
```
## 标志

| 名称, 速记       | 默认值    | 描述                                                         | 必需     |
| ---------------  | --------- | ------------------------------------------------------------ | -------- |
| --proposalID    |           | 软件升级提议的ID                                             | 是       |
| --title          |           | switch消息对标题                                             |          |

## 用例

发送对软件升级提议（ID为5）switch消息

```
iriscli upgrade submit-switch --chain-id=IRISnet --from=x --fee=0.004iris --proposalID 5 --title="Run new verison"
```
