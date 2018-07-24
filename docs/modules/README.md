# Upgrade Module

本module作为支持区块链软件平滑升级的基础设施，通过upgradeProposal和switch两阶段的投票来在约定高度切换到新版的代码，并对历史版本的链上数据完全兼容。

key point：

* upgradeProposal投票通过后各个节点下载安装新版本软件，并启动软件发送switch投票，表明已经可以切换到新版本
* switch的投票将在约定高度进行检查，需要95%的Voting Power才视为投票通过
* switch通过后开启切换流程：
   1. check_tx全部返回fail，以拒绝新tx的处理
   2. 处理mempool中留存的tx，直到生成一个empty block
   3. 配置新版本的路由开关，并打开check_tx接受新tx
* 发生老版本AppHash冲突，检查自己是否在switch voter list中，否则reset rootMultiStore到上一个commit，并下载新版本开始同步block
* 新增Module的升级方式（现有Module逻辑修改也通过新Module完成），新老版本的Module共享同一个store，对于查询需要路由到相应版本的Module
* Hardcord的升级方式（bug fix），Upgrade Module提供便利函数来决定指定的代码段在当前区块高度是否执行

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
