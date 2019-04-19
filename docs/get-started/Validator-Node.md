# Running a Validator Node

Before setting up your validator node, make sure you've already installed `Iris` by following this [guide](Full-Node.md) and your node is fully synced.

Validators are responsible for committing new blocks to the blockchain through consensus. A validator's stake will be slashed if it becomes unavailable, double-signs a transaction, or doesn't cast their votes. Please read about Sentry Node Architecture to protect your node from DDOS attacks and to ensure high-availability.

## Get IRIS Token

### Create Account

You need to get `iris` and `iriscli` installed first. Then, follow the instructions below to create a new account:

```
iriscli keys add <key_name>
```

Then, you should set a password of at least 8 characters.

The output will look like the following:
```
NAME:	TYPE:	ADDRESS:						PUBKEY:
tom	local	iaa1arlugktm7p64uylcmh6w0g5m09ptvklxrmsz9m	iap1addwnpepqvlmtpv7tke2k93vlyfpy2sxup93jfulll6r3jty695dkh09tekrz37h9q9
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

blast change tumble toddler rival ordinary chicken dirt physical club few language noise oak moment consider enemy claim elephant cruel people adult peanut garden
```

You could see the address and public key of this account. Please notice that account address in IRISnet will start with `iaa1` and public key of account will start with `iap1`.

The seed phrase of this account will also be displayed. You could use these 24 phrases to recover this account in another server. The recover command is:
```
iriscli keys add <key_name> --recover
```


### Claim tokens (Only for Fuxi Testnet)

You can always get some test tokens by using the [Faucet](https://testnet.irisplorer.io/#/faucet). The faucet will send you 10IRIS for every request, please don't abuse it.

Once you have created your own address,  you can use it to stake as a validator. The following command is used to check the balance of your account:
```
iriscli bank account <account_address> --node=http://localhost:26657
```

## Create Validator

### Confirm Your Validator is Synced

Your validator is active if the following command returns anything:

```
iriscli status --node=tcp://localhost:26657 
```

You should also be able to see `catching_up` is `false`. 

You need to get the public key of your node before upgrade your node to a validator node. The public key of your node starts with `icp`, it can be used to create a new validator by staking tokens. To understand more about the address encoding in IRIShub, 
please read this [doc](../features/basic-concepts/bech32-prefix.md)

You can find your validator's pubkey by running:

```
iris tendermint show-validator --home=<iris_home>
```
Example output:
```
icp1zcjduepq9l2svsakh9946n42ljt0lxv0kpwrc4v9c2pnqhn9chnjmlvagans88ltuj
```
Next, use the output as  `<pubkey>` field for `iriscli stake create-validator` command following [this](../cli-client/stake/create-validator.md). :

In this way, to stake 10IRIS and create as a validator, you need to do:

::: warning
**Create-validator need more gas and fee， you need to specify --gas=100000 --fee=0.6iris**
:::

```
iriscli stake create-validator --chain-id=<chain-id> --from=<key name> --gas=100000 --fee=0.6iris --pubkey=<validator public key> --amount=10iris --moniker=<your_custom_name> --commission-rate=0.1 --identity=<identity_string>
```
Please note the `fee` and `amount` can be the **decimal** of IRIS token, like `1.01iris`. And you could also use other coin-type like `iris-milli`, To read more about coin-type in IRIShub, you should read [this](../features/basic-concepts/coin-type.md)

`identity` is an optional field, please refer to [keybase](https://keybase.io/)

To read more about fee mechanism in IRIShub throughout the [doc](../features/basic-concepts/fee.md)

### View Validator Info

View the validator's information with this command:

```
iriscli stake validator <address-validator-operator> --chain-id=<chain-id> --node=tcp://localhost:26657 
```

The `<address-validator-operator>` is your account address that starts with 'iva1'


### Confirm Your Validator is Running

Your validator is active if the following command returns anything:

```
iriscli status --node=tcp://localhost:26657 
```

You should also be able to see your `voting_power` is above 0 if your bonded toke is in top 100. Also, you should see validator on the [Explorer](https://testnet.irisplorer.io).


### Edit Validator Description

You can edit your validator's public description following [this](../cli-client/stake/edit-validator.md). This info is to identify your validator, and will be relied on by delegators to decide which validators to stake to. Make sure to provide input for every flag below, otherwise the field will default to empty (`--moniker`defaults to the machine name).

You should put your name of your team in `details`. 

```
iriscli stake edit-validator --from=<key name> --moniker=<your_custom_name> --website=<your_website> --details=<your_details> --chain-id=<chain-id> --node=tcp://localhost:26657 --fee=0.3iris --identity=<identity_string>
```

`identity` is an optional field.

### Use the Explorer

You should also be able to see your validator on the [Explorer](https://www.irisplorer.io). 