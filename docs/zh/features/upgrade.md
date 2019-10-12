# 软件升级

## 简介

该模块是支持区块链软件平滑升级的基础设施，通过UpgradeProposal在约定高度切换到新版的代码，并对历史版本的链上数据完全兼容。

## 交互流程

### 软件升级提议治理流程

1. 用户提交升级软件的提议并且经过投票使该提议通过
2. 治理流程详细见GOV的[用户手册](governance.md)

### 升级软件流程

1. 用户一旦安装新软件，节点就会自动广播全网，本节点已经安装新软件。
2. 到达限定的时间（由软件升级提议决定），链上会统计升级到新软件的voting power比例是否超过软件升级的阈值（由软件升级提议决定）。
3. 如果超过，软件进行升级，否则升级失败。
4. 对于没有及时参与升级的节点，需要安装并运行新版本软件。

## 使用场景

### 创建使用环境

```bash
rm -rf iris
rm -rf .iriscli
iris init gen-tx --name=<key-name> --home=<path-to-your-home>
iris init --gen-txs --chain-id=<chain-id> -o --home=<path-to-your-home>
iris start --home=<path-to-your-home>
```

### 提交软件升级的提议

```bash
# 发送升级提议
iriscli gov submit-proposal --title=<title> --description=<description> --type="SoftwareUpgrade" --deposit=100iris --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --software=https://github.com/irisnet/irishub/tree/v0.13.1 --version=2 --switch-height=80 --threshold=0.9 --commit

# 对提议进行抵押
iriscli gov deposit --proposal-id=<proposal-id> --deposit=1000iris --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit

# 对提议投票
iriscli gov vote --proposal-id=<proposal-id> --option=Yes --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --commit

# 查询提议情况
iriscli gov query-proposal --proposal-id=<proposal-id>
```

### 升级软件

* 场景一

用户在指定的高度（例如80），完成以下动作：

```bash
# 1. 下载新版本iris1

# 2. 关闭旧软件
kill -f iris

# 3. 安装新版本 iris1 并启动（copy to bin）
iris1 start --home=<path-to-your-home>

# 4. 区块到达指定高度，自动升级

# 5. 查询当前版本是否升级成功
iriscli upgrade info --trust-node
```

* 场景二

用户在指定的高度（例如80），没有安装新软件，软件无法继续运行：

```bash
# 1. 下载新版本iris1

# 2. 关闭旧软件
kill -f iris

# 3. 安装新版本 iris1 并启动
iris1 start --home=<path-to-your-home>

# 4. 查询当前版本是否升级成功
iriscli upgrade info --trust-node
```

## 命令详情

```bash
iriscli gov submit-proposal --title=<title> --description=<description> --type="SoftwareUpgrade" --deposit=100iris --from=<key-name> --chain-id=<chain-id> --fee=0.3iris --software=https://github.com/irisnet/irishub/tree/v0.13.1 --version=2 --switch-height=80 --threshold=0.9 --commit
```

* `--type`  "SoftwareUpgrade" 软件升级提议的类型
* `--version`  新软件协议版本号
* `--software`  新软件的下载地址
* `--switch-height` 新软件升级的高度
* `--threshold`  软件升级的阈值
* 其他参数可参考Governance的[用户手册](governance.md)

```bash
iriscli upgrade info --trust-node
```

* 查询当前软件具体版本信息
