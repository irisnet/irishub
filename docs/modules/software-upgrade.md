# Upgrade Module

本module作为支持区块链软件平滑升级的基础设施，通过upgradeProposal和switch两阶段的投票来在约定高度切换到新版的代码，并对历史版本的链上数据完全兼容。

key point：

* upgradeProposal投票通过后各个节点下载安装新版本软件，并启动软件发送switch投票，表明已经可以切换到新版本
* switch的投票将在约定高度进行检查，需要95%的Voting Power才视为投票通过（switch只接受validator签名的）
* switch通过后开启切换流程：
   1. check_tx全部返回fail，以拒绝新tx的处理
   2. 处理mempool中留存的tx，直到生成一个empty block（社区会公告在第多少高度进行版本升级，提醒届时会终止服务，不要在升级时间段内发交易。所以这里算是一个防御检查，保守一点可以先等两个空块，后续压测发现不足的话再调整。）
   3. 配置新版本的路由开关，并打开check_tx接受新tx
* 发生老版本AppHash冲突:
   1. 检查自己是否在switch voter list中，否则reset rootMultiStore到上一个commit
   2. 下载新版本，运行iris start --replay命令启动节点
   3. --replay在完成app的setup后，start子命令会先回滚tendermint的statedb到上一高度。如果revert前的statedb对于的高度为H，那么现在state db高度为H-1。然后对于ABCI App store，只load到H-1的高度。然后就是正常启动节点，节点再检测到state db的高度和app store的高度都比block store低一个高度，就会自动replay最后一个block
   4. 开启tendermint的block sync，进行正常的区块追赶
* 新增Module的升级方式（现有Module逻辑修改也通过新Module完成），新老版本的Module共享同一个store，对于查询需要iriscli提供不同版本的数据解析能力
* Hardcord的升级方式（bug fix），Upgrade Module提供便利函数来决定指定的代码段在当前区块高度是否执行

## Data Struct

```
type ModuleLifeTime struct {
	Start	int64
	End	int64
        Handler sdk.Handler
	store	sdk.KVStoreKey
}

type Version struct {
	Id		int	 // should be equal with corresponding upgradeProposalID
	Start		int64
	ModuleList	[]ModuleLifeTime
}

```

## TxMsg

```
type MsgSwitch struct {
	Title          string
	ProposalId     int
	Voter          sdk.AccAddress
}

```

## Storage

| Key | Type   | Value | Description | Note|
| --------- | ------ | ------- | -------- | -----------|
| CurrentVersionIDKey | int | CurrentVersionID    | c/     |    |
| VersionKey | Version | Version    | v/%010d/     |  v/proposalId  |
| VersionListKey | ListOfVersionKey | [][]byte{}    | l/     |  list of the version_key ordered by proposalId  |
| SwitchKey | MsgSwitch | MsgSwitch    | s/%010d/%d/     | s/proposalId/switchVoterAddress | 

## Gov Propsal 设计
### add the module
Title:ADD

Description:"NEW MODULE NAME"
### change the module
Title:CHANGE

Description:"NEW MODULE NAME:OLD MOUDLE NAME"

### bugfix
Title:BUGFIX

Description:"Bug MODULE NAME"

## 不同阶段对store数据结构的读写
### store存储数据
* currentVersion 当前版本
* versionList 按时间顺序排列的versionID的list
* version 各个版本信息
* upgradeproposalID 当前升级的proposalID

### 创世阶段
1. 创建第一个version ，设置为currentVersion（id=0）
2. 创建一个versionList=[0]
3. 保存currentVersion，versionList，version，upgradeproposalID（-1）

### gov proposal 通过阶段
1. 把当前升级的proposalID 写入 upgradeproposalID

### MsgSwitch阶段
1. 判断voter是否为validator
2. 判断当前proposalID 是否为 upgradeproposalID
3. 1，2都成立，则把投票信息写入DB

## DoSwitch 路由开关

### 创世版本

从IrisApp中加载信息读取并构造version数据

### 升级版本

新版本代码中预定义了改版本的version信息,将其设置为当前版本
