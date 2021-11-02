# TIBC

> **_提示：_** 本文档中显示的命令仅供说明。有关命令的准确语法，请参阅[cli docs](../cli-client/tibc.md)。
## 简介

### 轻客户端
在TIBC协议中，参与者（可能是直接用户，链下进程或一台机器）需要能够验证另一台机器的共识算法已同意的状态更新，并拒绝另一台机器的共识算法尚未达成共识的任何可能的更新。轻客户端是机器可以执行的算法。该标准规范了轻客户端模型和需求，因此只要提供满足列出要求的相关轻客户端算法，
TIBC协议就可以轻松地与运行新的共识算法的新机器集成。TIBC的轻客户端在hub上采用Gov的方式进行管理，在其他链（例如以太坊）可以采用其他方式，比如多签。

#### 使用场景

1. 创建轻客户端提议
   
    如果在以太坊上，ChainName指以太坊上部署的某个轻客户端的合约地址.

    ```bash
        iris tx gov submit-proposal client-upgrade [chain-name] [path/to/client_state.json] [path/to/consensus_state.json]
    ```   


2. 升级轻客户端提议

   当轻客户端过期或者轻客户端状态不正确时，可以使用升级的方式强制更新轻客户端状态。

    ```bash
    iris tx gov submit-proposal client-upgrade [chain-name] [path/to/client_state.json] [path/to/consensus_state.json]
    ```

3. 中继器注册提议

   中继器主要功能：
    - 更新来源链在目标链上轻客户端状态
    - 中继跨链交易到目标链

   ```bash
    iris tx gov submit-proposal relayer-register [chain-name] [relayers-address]
   ```

4. 更新轻客户端状态

   用户提交跨链数据包的时候，必须先更新轻客户端状态，然后再提交交易。

    ```bash
    iris tx tibc update [chain-name] [path/to/header.json]
    ```
   
### 端口和数据包
Port规定了端口分配系统，通过该系统，模块可以绑定到 TIBC 处理程序分配的唯一命名的端口。然后可以在端口间传递Packet，并且可以通过最初绑定到它们的模块进行传输或之后释放。
Packet定义了链间数据包标准。发送和接收 TIBC 数据包的模块决定如何构造数据包数据以及如何处理传入的数据包数据，并且必须使用自己的应用程序逻辑根据数据包包含的数据来决定应用哪种状态交易。

#### 使用场景
1. 发送清理 Packet

   必须确认数据包才能进行清理。
   定义一个新的状态清理packet用于清理跨链数据包生命周期中产生的数据存储。该packet可以清理自身的存储。

    ```bash
   iris tx tibc packet send-clean-packet [dest-chain-name] [sequence] [flags]
   ```  
   

### 路由
路由模块维护一个Routing Rules作为路由白名单，其功能由管理模块或者Gov模块配置，配置白名单条目格式为src,dest,port，src为起始链chain ID，dest为目标链chain ID，port为模块绑定端口。
同时支持通配符，且仅对 Packet 有效，不拦截 ACK。

#### 使用场景
1. 设置路由规则
   
   所有的rule以string数组传入，拼成一条字符串进行持久化。

   ```bash
   iris tx gov submit-proposal set-rules [path/to/routing_rules.json] [flags]
   ```

### 跨链转NFT
通过 TIBC 协议连接的一组链的用户可能希望利用在另一条链上的一条链上发行的资产，也许是为了利用额外的功能，例如交换或隐私保护，
同时保留与发行链上原始资产的可互换性 . 该应用层标准描述了一种用于在与 TIBC 连接的链之间传输NFT的协议，该协议保留了资产的可替代性，保留了资产所有权，限制了拜占庭故障的影响，并且不需要额外的许可。

#### 使用场景
1. 发送 nft

   当使用跨链协议发送 nft时候，它们开始累积已传输的记录。 该信息被编码到class字段中。
   class字段以{prifix}/{sourceChain}/{destinationChain}{...}/class}形式实现,表示完整的跨链路径, prifix = "nft", 当没有prefix的时候，表示发送链就是NFT的源头，如果带有prefix和sourceChain，则表示是从sourceChain转过来的，比如nft/A/B/C/nftClass, 假如 C 链上存在此NFT:(tibc-hash(nft/A/B/C/nftClass))，则表示nftClass是从 A 转到 B 再转到 C 的, 如果 C 想退回此NFT，那么必须先退回到B，变成tibc-hash(nft/A/B/class) ,再由 B 链退回到 A变成class。
   目前只支持最多两次跳转，即只能从A到B再到C，从C回到B再回到A。
   ```bash
   iris tx tibc-nft-transfer transfer [dest-chain] [receiver] [class] [id] [flags]
   ```
   
### 中继器
在 TIBC 协议中，区块链只能记录将特定数据发送到另一条链的意图，而不能直接访问网络传输层。 物理数据报中继必须由可访问传输层（例如TCP/IP）的链外基础设施执行。 该标准定义了 relayer 算法的概念，该算法可由具有查询链状态能力的链外进程执行，以执行此中继。

#### 使用场景
1. 中继数据包

   可以以基于事件的方式中继无序信道中的数据包。 relayer 应在源链(source chain)中监视每当发送数据包时发出的事件，然后使用事件(event)日志中的数据来组成数据包。 随后， relayer 应通过在数据包的序列号查询是否存在确认来检查目的链是否已接收到该数据包，如果尚未出现， relayer 应中继该数据包。

