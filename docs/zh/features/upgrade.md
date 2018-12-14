# 软件升级用户文档

## 基本功能描述

该模块支持区块链软件平滑升级的基础设施，通过UpgradeProposal在约定高度切换到新版的代码，并对历史版本的链上数据完全兼容。

## 交互流程

### 软件升级提议治理流程
1. 用户提交升级软件的提议
2. 治理流程详细见GOV的[用户手册](governance.md)


### 升级软件流程  
1. 用户安装新软件，并发送switch消息，广播全网已经安装新软件。
2. 到达限定的时间，链上会统计升级到新软件的voting power比例是否超过95%。
3. 如果超过95%，软件进行升级，否则升级失败。
4. 对于没有及时参与升级的节点,需要重新下载新软件，同步区块。

## 使用场景

### 创建使用环境

```
rm -rf iris                                                                         
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=upgrade-test -o --home=iris
iris start --home=iris
```
### 提交软件升级的提议

```
# 发送升级提议
iriscli gov submit-proposal --title=Upgrade --description="SoftwareUpgrade" --type="SoftwareUpgrade" --deposit=10iris --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000 --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80

# 对提议进行抵押
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000

# 对提议投票
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000

# 查询提议情况
iriscli gov query-proposal --proposal-id=1 --trust-node
```

### 升级软件

* 场景一

用户在指定的高度（例如80），完成以下动作：

```
# 1. 下载新版本iris1

# 2. 关闭旧软件
kill -f iris

# 3. 安装新版本 iris1 并启动（copy to bin）
iris1 start --home=iris

# 4. 到达规定的时间，自动升级

# 5. 查询当前版本是否升级成功
iriscli upgrade info --trust-node
```

* 场景二

用户在指定的高度（例如80），没有安装新软件，软件无法继续运行：

```
# 1. 下载新版本iris1

# 2. 关闭旧软件
kill -f iris

# 3. 安装新版本 iris1 并启动
iris1 start --home=iris

# 4. 查询当前版本是否升级成功
iriscli upgrade info --trust-node
```

## 命令详情

```
iriscli gov submit-proposal --title=Upgrade --description="SoftwareUpgrade" --type="SoftwareUpgrade" --deposit=10iris --from=x --chain-id=upgrade-test --fee=0.05iris --gas=20000 --software=https://github.com/irisnet/irishub/tree/v0.9.0 --version=2 --switch-height=80
```

* `--type`  "SoftwareUpgrade" 软件升级提议的类型
* `--version`  "Version" 新软件协议版本号
* `--software`  新软件的下载地址
* `--switch-height` 新软件升级的高度
* 其他参数可参考GOV的[用户手册](governance.md)

```
iriscli upgrade submit-switch --name=x --from=$VADDR --proposalID=1 --chain-id=upgrade-test --fee=0.05iris --gas=20000
```

* `--proposalID` 当前通过的软件升级提议的ID

```
iris start --replay
```

* 重新同步区块，清理因升级而写脏的AppHash

```
iriscli upgrade info --trust-node
```

* 查询当前软件具体版本信息
