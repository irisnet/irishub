# nft-transfer test

## Overview

The following repository contains a basic example of an Interchain NFT module and serves as a developer guide for teams that wish to use Interchain NFT functionality.

The Interchain NFT module is now maintained within the `nft-transfer` repository
[here](https://github.com/bianjieai/nft-transfer).

### Developer Documentation

## Setup

1. Clone this repository and build the application binary

    ```bash
    git clone https://github.com/irisnet/irishub.git
    cd irishub

    make install 
    ```

2. Compile and install an IBC relayer.

    ```bash
    git clone https://github.com/cosmos/relayer.git
    cd relayer
    git checkout v2.4.2
    make install
    ```

3. Bootstrap two chains and create an IBC connection and start the relayer

    ```bash
    make init-golang-rly
    ```

## Demo

**NOTE:** For the purposes of this demo the setup scripts have been provided with a set of hardcoded mnemonics that generate deterministic wallet addresses used below.

```bash
# Store the following account addresses within the current shell env
export DEMOWALLET_1=$(iris keys show demowallet1 -a --keyring-backend test --home ./data/test-1) && echo $DEMOWALLET_1;
export DEMOWALLET_2=$(iris keys show demowallet2 -a --keyring-backend test --home ./data/test-2) && echo $DEMOWALLET_2;
```

### Issue an nft class on the `test-1` chain

Issue an nft class using the `iris tx nft issue` cmd.
Here the message signer is used as the account owner.

```bash
# Issue an nft class
iris tx nft issue cat --name xiaopi --symbol pipi --description "my cat" --uri "hhahahh"  --from demowallet1 --chain-id test-1 --keyring-dir ./data/test-1 --fees=1iris --keyring-backend=test --node tcp://127.0.0.1:16657 --mint-restricted=false  --update-restricted=false

# Query the class
iris query nft denom cat --node tcp://127.0.0.1:16657
```

### Mint a nft on the `test-1` chain

```bash
# Mint a nft
iris tx nft mint cat xiaopi --uri="http://wwww.baidu.com" --from demowallet1 --chain-id test-1 --keyring-dir ./data/test-1 --fees=1iris --keyring-backend=test  --node tcp://127.0.0.1:16657

# query the nft
iris query nft token cat xiaopi --node tcp://127.0.0.1:16657
```

### Transfer a nft from chain `test-1` to chain `test-2`

```bash
# Execute the nft tranfer command
iris tx nft-transfer transfer nft-transfer channel-0 iaa10h9stc5v6ntgeygf5xf945njqq5h32r5y7qdwl cat xiaopi --from demowallet1 --chain-id test-1 --keyring-dir ./data/test-1 --fees=1iris --keyring-backend=test --node tcp://127.0.0.1:16657 --packet-timeout-height 2-10000

# Query the newly generated class-id through class-trace
iris query nft-transfer class-hash nft-transfer/channel-0/cat --node tcp://127.0.0.1:26657

# Query nft information on test-2
iris query nft nft ibc/943B966B2B8A53C50A198EDAB7C9A41FCEAF24400A94167846679769D8BF8311 xiaopi --node tcp://127.0.0.1:26657
```

When the nft is transferred out, the nft on the original chain will be locked to the [escrow account](https://github.com/bianjieai/ibc-go/blob/develop/modules/apps/nft-transfer/types/keys.go#L45). You can use the following command to determine whether the transferred nft is escrow

```bash
iris query nft owner cat xiaopi --node tcp://127.0.0.1:16657
```

### Transfer a nft back from chain `test-2` to chain `test-1`

```bash
iris tx nft transfer nft-transfer channel-0 iaa10h9stc5v6ntgeygf5xf945njqq5h32r5y7qdwl ibc/943B966B2B8A53C50A198EDAB7C9A41FCEAF24400A94167846679769D8BF8311 xiaopi --from demowallet2 --chain-id test-2 --keyring-dir ./data/test-2 --fees=1stake --keyring-backend=test --node tcp://127.0.0.1:26657 --packet-timeout-height 1-10000
```
