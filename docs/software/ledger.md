# Ledger Nano Support

## A note on HD wallet

HD Wallets, originally specified in Bitcoin's [BIP32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki), are a special kind of wallet that let users derive any number of accounts from a single seed. To understand what that means, let us first define some terminology:

- **Wallet**: Set of accounts obtained from a given seed.
- **Account**: A pair of public key/private key. 
- **Private Key**: A private key is a secret piece of information used to sign messages. In the blockchain context, a private key identifies the owner of an account. The private key of a user should never be revealed to others.
- **Public Key**: A public key is a piece of information obtained by applying a one-way mathematical function on a private key. From it, an address can be derived. A private key cannot be found from a public key.
- **Address**: An address is a public string with a human-readable prefix that identifies an account. It is obtained by applying mathematical transformations to a public key.
- **Digital Signature**: A digital signature is a piece of cryptographic information that proves the owner of a given private key approved of a given message without revealing the private key.
- **Seed**: Same as Mnemonic.
- **Mnemonic**: A mnemonic is a sequence of words that is used as seed to derive private keys. The mnemonic is at the core of each wallet. NEVER LOSE YOUR MNEMONIC. WRITE IT DOWN ON A PIECE OF PAPER AND STORE IT SOMEWHERE SAFE. IF YOU LOSE IT, THERE IS NO WAY TO RETRIEVE IT. IF SOMEONE GAINS ACCESS TO IT, THEY GAIN ACCESS TO ALL THE ASSOCIATED ACCOUNTS.

At the core of a HD wallet, there is a seed. From this seed, users can deterministically generate accounts. To generate an account from a seed, one-way mathematical transformations are applied. To decide which account to generate, the user specifies a `path`, generally an `integer` (`0`, `1`, `2`, ...). 

By specifying `path` to be `0` for example, the Wallet will generate `Private Key 0` from the seed. Then, `Public Key 0` can be generated from `Private Key 0`.  Finally, `Address 0` can be generated from `Public Key 0`. All these steps are one way only, meaning the `Public Key` cannot be found from the `Address`, the `Private Key` cannot be found from the `Public Key`, ...

```
     Account 0                         Account 1                         Account 2

+------------------+              +------------------+               +------------------+
|                  |              |                  |               |                  |
|    Address 0     |              |    Address 1     |               |    Address 2     |
|        ^         |              |        ^         |               |        ^         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        +         |              |        +         |               |        +         |
|  Public key 0    |              |  Public key 1    |               |  Public key 2    |
|        ^         |              |        ^         |               |        ^         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        |         |              |        |         |               |        |         |
|        +         |              |        +         |               |        +         |
|  Private key 0   |              |  Private key 1   |               |  Private key 2   |
|        ^         |              |        ^         |               |        ^         |
+------------------+              +------------------+               +------------------+
         |                                 |                                  |
         |                                 |                                  |
         |                                 |                                  |
         +--------------------------------------------------------------------+
                                           |
                                           |
                                 +---------+---------+
                                 |                   |
                                 |  Mnemonic (Seed)  |
                                 |                   |
                                 +-------------------+
```

The process of derivating accounts from the seed is deterministic. This means that given the same path, the derived private key will always be the same. 

The funds stored in an account are controlled by the private key. This private key is generated using a one-way function from the mnemonic. If you lose the private key, you can retrieve it using the mnemonic. However, if you lose the mnemonic, you will lose access to all the derived private keys. Likewise, if someone gains access to your mnemonic, they gain access to all the associated accounts. 

## Ledger Support for account keys

At the core of a Ledger device, there is a mnemonic that is used to generate private keys. When you initialize you Ledger, a mnemonic is generated. 

::: danger
**Do not lose or share your 24 words with anyone. To prevent theft or loss of funds, it is best to ensure that you keep multiple copies of your mnemonic, and store it in a safe, secure place and that only you know how to access. If someone is able to gain access to your mnemonic, they will be able to gain access to your private keys and control the accounts associated with them.**
:::

This mnemonic is compatible with IrisNet accounts. The tool used to generate addresses and transactions on the IrisNet network is called `iriscli`, which supports derivation of account keys from a Ledger seed. Note that the Ledger device acts as an enclave of the seed and private keys, and the process of signing transaction takes place within it. No private information ever leaves the Ledger device. 

To use `iriscli` with a Ledger device you will need the following(Since IRISnet is based on cosmos-sdk, the COSMOS app is available for IRISnet):

- [A Ledger Nano with the `COSMOS` app installed and an account](#using-a-ledger-device)
- [A running `iris` instance connected to the network you wish to use.](../get-started/Join-the-Mainnet.md)
- [A `iriscli` instance configured to connect to your chosen `iris` instance.](./cli-client.md)

Now, you are all set to start sending transactions on the network.

At the core of a ledger device, there is a mnemonic used to generate accounts on multiple blockchains (including the IRISnet). Usually, you will create a new mnemonic when you initialize your ledger device.

Next, click [here](#using-a-ledger-device) to learn how to generate an account.

## Creating an account

To create an account, you just need to have `iriscli` installed. Before creating it, you need to know where you intend to store and interract with your private keys. The best options are to store them in an offline dedicated computer or a ledger device. Storing them on your regular online computer involves more risk, since anyone who infiltrates your computer through the internet could exfiltrate your private keys and steal your funds.

#### Using a ledger device

::: warning
**Only use Ledger devices that you bought factory new or trust fully**
:::

When you initialize your ledger, a 24-word mnemonic is generated and stored in the device. This mnemonic is compatible with IRISnet and IRISnet accounts can be derived from it. Therefore, all you have to do is make your ledger compatible with `iriscli``. To do so, you need to go through the following steps:

1. Download the Ledger Live app [here](https://www.ledger.com/pages/ledger-live). 
2. Connect your ledger via USB and update to the latest firmware
3. Go to the ledger live app store, and download the "Cosmos" application (this can take a while). **Note: You may have to enable `Dev Mode` in the `Settings` of Ledger Live to be able to download the "Cosmos" application**. 
4. Navigate to the Cosmos app on your ledger device

Then, to create an account, use the following command:

```bash
iriscli keys add <yourAccountName> --ledger 
```


