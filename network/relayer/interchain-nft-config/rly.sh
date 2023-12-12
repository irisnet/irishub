#!/bin/bash

# Configure predefined mnemonic pharses
BINARY=rly
CHAIN_DIR=./data
RELAYER_DIR=./relayer
MNEMONIC_1="alley afraid soup fall idea toss can goose become valve initial strong forward bright dish figure check leopard decide warfare hub unusual join cart"
MNEMONIC_2="record gift you once hip style during joke field prize dust unique length more pencil transfer quit train device arrive energy sort steak upset"

# Ensure rly is installed
if ! [ -x "$(command -v $BINARY)" ]; then
    echo "$BINARY is required to run this script..."
    echo "You can download at https://github.com/cosmos/relayer"
    exit 1
fi

echo "Initializing $BINARY..."
$BINARY config init --home $CHAIN_DIR/$RELAYER_DIR

echo "Adding configurations for both chains..."
$BINARY chains add -f $PWD/network/relayer/interchain-nft-config/chains/test-1.json test-1 --home $CHAIN_DIR/$RELAYER_DIR
$BINARY chains add -f $PWD/network/relayer/interchain-nft-config/chains/test-2.json test-2 --home $CHAIN_DIR/$RELAYER_DIR
$BINARY paths add test-1 test-2 test1-nft-test2 -f $PWD/network/relayer/interchain-nft-config/paths/test1-nft-test2.json --home $CHAIN_DIR/$RELAYER_DIR

echo "Restoring accounts..."
$BINARY keys restore test-1 test-1 "$MNEMONIC_1" --home $CHAIN_DIR/$RELAYER_DIR
$BINARY keys restore test-2 test-2 "$MNEMONIC_2" --home $CHAIN_DIR/$RELAYER_DIR

echo "Linking both chains and starting to listen relayer"
$BINARY transact link test1-nft-test2  --src-port nft-transfer --dst-port nft-transfer --order unordered --version ics721-1 --max-retries 20 -d --home $CHAIN_DIR/$RELAYER_DIR

echo "Starting to listen relayer..."
$BINARY start test1-nft-test2 --home $CHAIN_DIR/$RELAYER_DIR