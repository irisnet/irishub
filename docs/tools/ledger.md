---
order: 4
---

# Ledger Nano Support

It is recommended to have a basic understanding of the [IRISnet Key](../concepts/key.md) before using the ledger.

## Ledger Support for account keys

At the core of a Ledger device, there is a mnemonic that is used to generate private keys. When you initialize you Ledger, a mnemonic is generated.

::: danger
**Do not lose or share your 24 words with anyone. To prevent theft or loss of funds, it is best to ensure that you keep multiple copies of your mnemonic, and store it in a safe, secure place and that only you know how to access. If someone is able to gain access to your mnemonic, they will be able to gain access to your private keys and control the accounts associated with them.**
:::

This mnemonic is compatible with IrisNet accounts. The tool used to generate addresses and transactions on the IrisNet network is called `iriscli`, which supports derivation of account keys from a Ledger seed. Note that the Ledger device acts as an enclave of the seed and private keys, and the process of signing transaction takes place within it. No private information ever leaves the Ledger device.

To use `iriscli` with a Ledger device you will need the following(Since IRISnet is based on cosmos-sdk, the COSMOS app is available for IRISnet):

- [A Ledger Nano with the `COSMOS` app installed and an account.](#using-a-ledger-device)
- [A running `iris` instance connected to the network you wish to use.](../get-started/mainnet.md)
- [A `iriscli` instance configured to connect to your chosen `iris` instance.](../cli-client/intro.md)

Now, you are all set to start sending transactions on the network.

At the core of a ledger device, there is a mnemonic used to generate accounts on multiple blockchains (including the IRISnet). Usually, you will create a new mnemonic when you initialize your ledger device.

Next, click [here](#using-a-ledger-device) to learn how to generate an account.

## Creating an account

To create an account, you just need to have `iriscli` installed. Before creating it, you need to know where you intend to store and interact with your private keys. The best options are to store them in an offline dedicated computer or a ledger device. Storing them on your regular online computer involves more risk, since anyone who infiltrates your computer through the internet could exfiltrate your private keys and steal your funds.

### Using a ledger device

::: warning
Only use Ledger devices that you bought factory new or trust fully
:::

When you initialize your ledger, a 24-word mnemonic is generated and stored in the device. This mnemonic is compatible with IRISnet and IRISnet accounts can be derived from it. Therefore, all you have to do is make your ledger compatible with `iriscli`. To do so, you need to go through the following steps:

1. Download the Ledger Live app [here](https://www.ledger.com/pages/ledger-live).
2. Connect your ledger via USB and update to the latest firmware
3. Go to the ledger live app store, and download the `Cosmos` application (this can take a while). **Note: You may have to enable `Dev Mode` in the `Settings` of Ledger Live to be able to download the "Cosmos" application**.
4. Navigate to the Cosmos app on your ledger device

Then, to create an account, use the following command:

```bash
iriscli keys add <yourAccountName> --ledger
```
