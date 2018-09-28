# Upgrade User Guide

## 基本功能描述

该模块支持区块链软件平滑升级的基础设施，通过UpgradeProposal和switch两阶段的投票来在约定高度切换到新版的代码，并对历史版本的链上数据完全兼容。
## 交互流程

### 软件升级提议治理流程
1. 用户提交升级软件提议
2. 治理流程详细见[]()
3. 
### 升级软件流程  
1. 用户安装新软件，并发送switch消息，广播全网自己已经安装新软件。
2. 

## 使用场景

iriscli upgrade submit-switch --name=x --from=$VADDR --proposalID=1 --chain-id=upgrade-test --fee=20000000000000000iris

iris start --replay

## 命令详情