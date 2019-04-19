# IRISnet主网

## IRIShub 简介

IRIS Hub是在Cosmos生态中的区域性枢纽，提供iService服务

## 如何加入IRIShub网络

### 第一步: 安装IRIShub

请根据以下[教程](../software/How-to-install-Irishub.md) 在服务器上完成`iris`的安装。

### 第二步: 运行一个全节点

请根据以下[步骤](Full-Node.md) 完成初始化并且在服务器上部署一个全节点。

### 第三步: 将全节点升级成为一个验证人节点

请根据以下[步骤](Validator-Node.md) 将一个全节点升级成为验证人节点。

### 部署IRIShub Monitor监控

请根据以下[链接](../software/monitor.md) 在服务器上部署一个Monitor监控。


### 如何成为一个验证人节点

如何你的节点已经完成同步了，那么接下来你应该：

如果你参与到了genesis文件的生成过程中，那么只要你的节点与网络同时启动，它就会保持验证人的状态。

如果你并没有参与到genesis文件的生成过程中，那么你依然可以通过执行相关操作升级成为一个验证人。目前IRIShub的验证人上限是100(根据委托总量排名，超过100名的将成为候选人)。升级的流程在[这里](Validator-Node.md).

### 部署哨兵节点

验证人有遭受攻击的风险。你可以根据以下[教程](../software/sentry.md)部署一个哨兵节点来保护验证人。

### 使用KMS
如果您打算使用KMS（密钥管理系统），则应首先执行以下步骤：[使用KMS](../software/kms/kms.md)。

##  更多链接


* Explorer: https://www.irisplorer.io/#/home

* Riot chat: #irisvalidators:matrix.org

* IRIShub验证人工作QQ群：834063323