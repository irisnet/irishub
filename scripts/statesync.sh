# UNCOMMENT OR COMMENT ITEMS AS NEEDED TO FIT YOUR SETUP
#!/bin/bash
# microtick and bitcanna contributed significantly here.
set -uxe

# set environment variables
export GOPATH=~/go
export PATH=$PATH:~/go/bin


# Install Iris
go install -tags rocksdb ./...


# MAKE HOME FOLDER AND GET GENESIS
iris init test 
wget -O ~/.iris/config/genesis.json https://github.com/irisnet/mainnet/raw/master/config/genesis.json

INTERVAL=1000

# GET TRUST HASH AND TRUST HEIGHT

LATEST_HEIGHT=$(curl -s http://seed-2.mainnet.irisnet.org:26657/block | jq -r .result.block.header.height);
BLOCK_HEIGHT=$(($LATEST_HEIGHT-$INTERVAL))
TRUST_HASH=$(curl -s "http://seed-2.mainnet.irisnet.org:26657/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)


# TELL USER WHAT WE ARE DOING
echo "TRUST HEIGHT: $BLOCK_HEIGHT"
echo "TRUST HASH: $TRUST_HASH"


# exort state sync vars
export IRIS_STATESYNC_ENABLE=true
export IRIS_P2P_MAX_NUM_OUTBOUND_PEERS=500
export IRIS_STATESYNC_RPC_SERVERS="http://seed-2.mainnet.irisnet.org:26657,http://seed-2.mainnet.irisnet.org:26657"
export IRIS_STATESYNC_TRUST_HEIGHT=$BLOCK_HEIGHT
export IRIS_STATESYNC_TRUST_HASH=$TRUST_HASH
export IRIS_GRPC_WEB_ADDRESS=127.0.0.1:9999
export IRIS_P2P_PERSISTENT_PEERS="a17d7923293203c64ba75723db4d5f28e642f469@seed-2.mainnet.irisnet.org:26656"


# WE ARE USING ROCKSDB BECAUSE OF THE LARGE PERFORMANCE BOOST IT OFFERS
iris start --x-crisis-skip-assert-invariants --db_backend rocksdb 