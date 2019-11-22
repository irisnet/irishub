---
order: 4
---

# Ledger Nano支持

建议在使用ledger之前对[IRISnet Key](../concepts/key.md)进行一个基本的了解。

## Ledger支持账户密钥

在Ledger设备的核心，有一个用于生成私钥的助记词。初始化Ledger时，会生成助记词。

::: danger
**不要丢失或与任何人分享你的24个单词。为防止盗窃或资金损失，最好确保保留多个助记词副本，并将其存放在安全可靠的地方，并且只有你知道如何访问。如果有人能够访问你的助记词，他们将能够访问你的私钥并控制与其关联的帐户。**
:::

该助记词与IRISnet兼容，用于在IRISnet网络上生成地址和交易的工具是`iriscli`，它支持从Ledger种子派生帐户密钥。请注意，Ledger设备充当种子和私钥的保护，并且签名交易发生在其中。没有任何私人信息会在Ledger设备之外被获取。

要通过Ledger设备使用`iriscli`，你需要采取下面步骤（因为IRISnet基于cosmos-sdk，COSMOS应用程序可用于IRISnet）：

- [在Ledger Nano中安装`COSMOS`应用程序并且创建一个账户。](#using-a-ledger-device)
- [一个连接在区块链网络上的`iris`全节点](../get-started/mainnet.md)
- [一个已配置`iriscli`（且已经连接到`iris`全节点）的实例。](../cli-client/intro.md)

至此可以在网络上发送交易。

在ledger设备的核心，有一个助记词用于在多个区块链（包括IRISnet）上生成帐户。通常，在初始化ledger设备时，将创建新的助记词。

接下来，点击[这里](#using-a-ledger-device)以了解如何生成帐户。

## 创建账户

要创建一个帐户，只需要安装`iriscli`。在创建它之前，需要知道你打算存储私钥的位置以及与私钥的交互。最好的选择是将它们存储在离线专用计算机或ledger设备中。将它们存储在你的常规联网的计算机上会带来更大的风险，因为任何通过互联网侵入你的计算机的人都可以获取你的私钥并窃取你的资金。

### 使用ledger设备

::: warning
仅使用购买的新设备或完全信任的Ledger设备
:::

初始化ledger时，会生成一个24个单词的助记词并存储在设备中。此助记词与IRISnet兼容，可以从中派生出IRISnet账号。因此，你所要做的就是让你的ledger与`iriscli`兼容。为此，你需要执行以下步骤：

1. 在[这里](https://www.ledger.com/pages/ledger-live)下载Ledger Live app。
2. 通过USB连接ledger并更新到最新的firmware。
3. 在ledger live应用商店中下载`Cosmos`application（这可能需要一段时间）。**你可能需要在Ledger Live的`Settings`中启用`Dev Mode`才能下载`Cosmos`应用程序。**
4. 切换到ledger设备上的COSMOS应用。

然后，要创建帐户，请使用以下命令：

```bash
iriscli keys add <yourAccountName> --ledger
```
