
# IRISnet 验证人常见问题

## IRIS网络介绍

IRIS网络是跨链服务基础设施和协议，用于构建可信的分布式商业应用。IRIS网络将面向服务的基础设施融入到Cosmos网络中。它整合由异构系统提供的商业服务，包括公共链、联盟链以及现有系统。通过区块链互联网实现服务的互联互通。


## IRIS枢纽

在IRIS网络的第一个区块链称为 *IRIS枢纽* 的区块链，它是一个基于Tendermint共识引擎构建的Bond-Proof-of-Stake（BPoS）区块链。它将成为第一个连接Cosmos枢纽的区域性枢纽。IRIS网络通过标准的ABCI交易实现IRIS服务（也称为iServices）的注册，绑定，调用，查询，分析和管理。iService提供商充当公共链、联盟区块链以及现有企业系统中业务逻辑的适配器。IRIShub将于从2019年春节后启动，这也标志着[主网启动](https://github.com/irisnet/iris-foundation/blob/master/iris-betanet-plan_cn.md)的第一步。

## IRIS服务

![IRIS网络图](https://github.com/irisnet/irisnet/blob/master/images/chap2-1.png?raw=true)

引入 *IRIS服务*的目标是弥合区块链世界与传统商业应用世界之间的鸿沟，居中协调链下服务的整个生命周期 - 从定义，绑定（提供者注册），调用，直到它们的治理。

### 生命周期

- **定义：** 根据接口定义语言（IDL）文件定义链下iService可以做什么。
- **绑定：** 声明实现给定iService定义的提供者端点的位置（地址），定价和服务质量。
- **调用：** 处理针对给定iService提供者端点的消费者请求以及相应的提供者响应。

### IRIS服务提供者

*提供者* 是提供一个或多个iService定义实现的网络用户，通常充当位于其他公有链、联盟链以及企业现有系统中的链下服务和资源的 *适配器*。它们监听和处理传入的请求，并将响应发送回网络。提供者可以通过向其他提供者发送请求来同时充当消费者。按照计划，提供者将被要求为他们可能提供的服务收取费用，默认情况下，服务费将以IRIS通证定价。

### IRIS服务消费者

*消费者* 是那些使用iServices的用户，他们向指定的提供者端点发送请求并接收相关提供者的响应。

### IRIS服务分析员

*分析员* 是一种特殊用户，代表了发起建立IRIS网络的IRIS基金会有限公司（IRIS Foundation Limited），这是一家注册在香港的股份有限公司。分析员是在分析模式中调用iServices的唯一授权用户，旨在帮助创建和维护服务提供者的概要文件，通过这些客观的概要文件服务消费者可以选择合适的服务提供者。

### IRIS服务仲裁员

*仲裁员* 是自我声明的一类用户，他们通过协作为消费者对提供者绩效的投诉进行仲裁。有关仲裁机制的细节正在积极设计中，请关注我们的[白皮书](resources/whitepaper-zh.md)。



## IRISnet网络节点类型

IRISnet是一个基于 [Tendermint](https://cosmos.network/docs/introduction/tendermint.html) 的BPoS 区块链网络，验证节点是网络运行的核心支柱。验证节点是全节点，他们的维护者成为验证人，其主要任务是**运行并保障验证节点稳定、安全地运行**。想要成为IRIS网络中的验证节点，验证人必须运维一个全节点，并抵押一定数量的**IRIS**通证，抵押通证总量也包括委托人委托给验证人的数量。抵押通证最多的100个全节点将被选为验证节点。各类节点分类如下：  
- **全节点**
  全节点即是一个具备完全功能的IRISnet 节点，需要[安装完整的IRIShub软件](https://www.irisnet.org/docs/zh/get-started/Install-the-Software.html)，并完成相应的全节点[配置](https://www.irisnet.org/docs/zh/get-started/Full-Node.html)。全节点具有以下特点： 
  - 投票权重为零；
  - 保存完整的交易账本；
  - 可以作为候选验证人节点。 
- **验证人节点**
  - 完成交易的验证打包，对区块进行共识投票。简而言之，IRISnet网络的验证人所起的作用类似与PoW网络的矿工。
  - 参与网络治理，对社区提案进行投票；
  - 及时升级软件，保障IRISnet的持续成长
- **候选验证节点**
  如果有超过100个全节点申请加入验证节点集，那么只有具有抵押通证数量排名前100的节点才能成为真正的验证节点，其他人将是候选验证节点。随着抵押数量的动态变化，验证节点集合和候选验证节点集合会不断变化
  - 投票权重为零；
  - 没有抵押获利


## IRIS通证

IRIS枢纽有自己的原生通证，称为 *IRIS*，在网络中有三个作用。

- **抵押：** 与Cosmos Hub中的ATOM通证类似，IRIS通证将用作抵押通证以保护PoS区块链的安全运行。
- **交易费用：** IRIS通证也将用于支付IRIS网络中所有交易的费用。
- **服务费：** IRIS网络中的服务提供者需要以IRIS通证为单位收取服务费。

IRIS网络最终将支持来自Cosmos网络的所有列入白名单的费用通证，它们可用于支付交易费用和服务费用。

## IRIS网络中用户类型
###  IRISnet验证人
在IRISHub枢纽中，运维验证节点的用户成为验证人。验证人的核心资产是验证节点运行的系统，最重要的是验证节点用于对区块和投票签名的密钥。
###  IRISnet委托人
部署一个验证节点需要付出很多努力，作为IRIS通证持有者，你可以通过委托(Delegate)的方式，将自己的通证抵押出去，同样获得抵押获利。但是你可以随时转换验证人，同时也需要按照验证人的要求缴纳一定比例的佣金。

### IRISnet Profiler账户
只有Profiler账户可以提交链上软件升级SoftwareUpgrade/停止共识Halt提议。
### IRISnet Trustee账户
只有Trustee账户能够发起链上治理中的TxTaxUsage提议，把社区资金转移给Trustee账户使用



##  IRISnet验证人的权益和责任

###  验证人收益分析
验证人可以获得收益有以下3种：
* **出块奖励**
  在IRISnet网络中，所有验证节点将轮流出块。出块的概率和抵押IRIS的数量成正比。作为出块人，验证人将获得额外的出块奖励。
  
* **抵押获利**
  接受IRIS通证持有人的委托，抵押委托人的代币用做共识投票的权益证明，并与委托人分享所获得的收益；IRISnet是基于Tendermint的PoS网络，验证人在网络共识中的投票权取决于验证人（包括受委托）[抵押通证（IRIS）的数量](https://www.irisnet.org/docs/zh/features/stake.html#%E4%BB%8B%E7%BB%8D)。网络中抵押的通证数量越多，攻击网络所需的成本也越大，网络也越安全。为了维护验证人及其委托人抵押通证的价值，IRISnet设定通胀增发通证，用于激励验证人及所有IRIS通证持有人将通证抵押。抵押获利[定时发放](https://www.irisnet.org/docs/zh/features/mint.html#%E5%8C%BA%E5%9D%97%E6%97%B6%E9%97%B4)。 通胀率也是社区对[IRIS通证抵押率的调节器](https://www.irisnet.org/docs/zh/features/mint.html#%E9%80%9A%E8%83%80%E7%8E%87)，IRISnet在第一年将通胀设定在4%， 后期可以通过社区提案的形式通过在线治理投票调整。网络中IRIS数量将逐年增加。没有被抵押的IRIS的价值将会逐渐被稀释。通胀参数可以通过发起链上治理投票来修改。

* **手续费** 
在IRIS网络中的各种交易都需要支付一定的手续费。手续费的多少取决于每种交易的[手续费上限 (fee)](https://www.irisnet.org/docs/zh/features/basic-concepts/fee.html#fee)和[交易消耗资源 (Gas) ](https://www.irisnet.org/docs/zh/features/basic-concepts/fee.html#gas)。IRISnet网络的中设定全局的手续费/Gas最小比例。

抵押获利和手续费会在验证人和委托人之间按比例分配。


###  验证人风险分析
Cosmos有一个关键的特性来激励验证人节点保持安全：惩罚那些看起来做了错事的验证人节点。
验证人在获得收益的同时，也必然要履行相应的义务，如果验证人出现双签，不在线或者不投票的情况，都将受到相应的惩罚，包括抵押的通证被核销，取消验证人资格并被标记为关押（无验证人候选资格）。 

-   **双重签名**：如果验证节点对同一次共识过程多次投票，并且这些投票相互矛盾，则对应的验证人抵押5%的通证将被罚没。一旦验证人被发现出现了双重签名的作恶行为，则验证节点会被处罚不能参加共识2天，这样也无法获得抵押获利。
-   **不在线**： 验证节点长期不参与网络共识，即如果一个验证节点的签名没有被最近的2万个区块包含，则对应的验证人抵押1%的通证将被罚没。一旦验证人被发现出现了双重签名的作恶行为，则验证节点会被处罚不能参加共识1天，这样也无法获得抵押获利。
-   **提交无效交易**：验证人有义务保证区块数据的有效性。验证人通过提交不合法的交易进入区块来扰乱网络共识，则对应的验证人抵押2%的通证将被罚没。一旦验证人被发现出现了双重签名的作恶行为，则验证节点会被处罚不能参加共识7天，这样也无法获得抵押获利。
-   **不投票**： 验证人有义务参与网络治理，对社区提案进行投票。如果一个验证人没有对某个提议投票表决，验证人抵押0.7%的通证将被罚没。

以上参数都可以通过发起链上治理投票来修改。

### 通证管理
* IRIS通证类型
在IRISnet网络中存在两种类型的IRIS，一种为*可流通*的，一种为*绑定*的。可流通的通证可以在账户间互相转账，也可以在交易所交易。验证人和委托人通过[委托操作](https://www.irisnet.org/docs/cli-client/stake/delegate.html)将可自由流通的IRIS通证变为*绑定*的IRIS。验证节点的[投票权重](https://www.irisnet.org/docs/zh/features/stake.html#%E4%BB%8B%E7%BB%8D)与验证人绑定的通证数量(包括委托通证)成正比。
* IRIS通证解绑
若委托人需要取回已委托的通证，可以通过发送[解绑交易](https://www.irisnet.org/docs/cli-client/stake/unbond.html)。在IRISnet网络中，解绑期默认为**三周**。一旦解绑期结束，被绑定的通证将自动成为流通通证。解绑期机制对PoS区块链网络的安全性很重要。

* IRIS通证再委托
委托人可以随时更换受委托验证人，当更换受委托验证人时，委托人抵押的通证可直接转入新的受委托人抵押池，而不无需等待3周的解绑期。
* IRIS通证收益管理
委托人通过抵押获得的抵押获利和手续费收益都是可流通的IRIS，可随时交易。


## 如何成为IRISnet验证人
若想成为验证人，你需要先参与Fuxi测试网。
1. 安装软件：
https://github.com/irisnet/irishub/blob/master/docs/zh/get-started/Install-the-Software.md
2. 部署全节点：https://github.com/irisnet/irishub/blob/master/docs/zh/get-started/Full-Node.md
3. 升级成为验证人节点：
https://github.com/irisnet/irishub/blob/master/docs/zh/get-started/Validator-Node.md

测试网包含许多奖励将在主网上线后发放给社区成员。你需要用你的Keybase签名一个irishub的地址，然后将其发送给团队。
**建议**
为了保障验证人节点安全，稳定的持续工作；提高自身形象，以便获得更多委托人信任，增加自身的投票权重，你还可以： 
 - 公布安全审计结果
 - 建立并公布相应的节点运维流程
 - 设立自己的验证人网站向委托人提供更加透明的信息，树立一个可信的受托人形象。

## IRISnet网络生态
**钱包**
- Rainbow 钱包 https://www.rainbow.one
- IRIS通证：测试网可以通过[水龙头](https://www.irisnet.org/docs/zh/get-started/Validator-Node.html#%E8%8E%B7%E5%BE%97iris%E4%BB%A3%E5%B8%81)获得；[主网上线](https://www.irisnet.org/mainnet?lang=CN)后IRIS基金会将开展空投，验证人私募等活动帮助验证人获得IRIS通证。同时，验证人也可通过Rainbow钱包验证人介绍及验证人推广等活动吸引委托人参与。


**浏览器**
IRISplorer：https://testnet.irisplorer.cn

##  IRISnet 测试网与主网

**IRISnet测试网 - Fuxi** 

IRISnet在开发测试阶段建立的测试网取名伏羲（Fuxi），伏羲测试网现在已经迭代到Fuxi-8000，这也是进入主网前的最后一个测试网。有意愿成为验证人的伙伴现在就可以加入Fuxi测试网熟悉IRISnet，IRIS基金会也会在主网上线后向参与和完成测试任务的[验证人提供奖励](https://github.com/irisnet/testnets/blob/master/fuxi/docs_CN/Fuxi_FAQ.md)。测试网相关信息可以通过[测试网主页](https://www.irisnet.org/testnets)查询。 

伏羲测试网在主网上线后将与主网平行运行，在功能上保持与主网一致，用以新功能投入主网前的试运行，和为用户提供一个应用开发的体验、测试平台。 

**IRISnet主网**

IRISnet主网的第一步：Betanet将在2019年2月中旬启动，主网启动的相关信息将在[IRIS基金会官网](https://www.irisnet.org/mainnet?lang=CN)更新。 

##  验证人交流渠道

- Riot chat: #irisvalidators:matrix.org
- IRIShub验证人工作QQ群：834063323




# Fuxi测试网激励计划常见问题


1.如何加入Fuxi测试网?

你可以加入QQ工作群：834063323。团队将在群里及时通知有关测试网的消息。

有两种方式加入测试网：
* 以验证人的身份加入：你可以在自己的服务器上部署一个IRIShub节点。然后将其绑定成为一个验证人节点。如果你暂时没有服务器，你也可以申请免费试用BaaS的服务，我们将提供
Wancloud和Zig-BaaS的免费试用机会。然后你就可以完成测试网的任务了。

* 以委托人的身份加入：
如果你对于部署一个验证人节点感到很困难，你可以只下载一个客户端，然后执行相关的测试网的任务交易。


2.测试网的激励任务在哪里?

每一个测试网迭代中，团队都会发布相关的测试网激励任务。例如，Fuxi-3001测试网激励任务在下面链接中：https://github.com/irisnet/testnets/tree/master/fuxi/fuxi-3001

3.怎样能知道我的任务完成情况?

每当决定切换到下一个测试网的时候，团队会检查参与者的任务完成情况。例如，Fuxi-2000测试网任务完成情况在这里：https://github.com/irisnet/testnets/issues/51

4.如果我获得了测试网奖励，何时才能拿到?

测试网的奖励将在主网上线后发放给社区成员。你需要用你的Keybase签名一个irishub的地址，然后将其发送给团队。