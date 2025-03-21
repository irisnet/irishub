#!/bin/bash

BINARY=iris
CHAIN_DIR=./data
CHAINID_1=test-1
CHAINID_2=test-2
GRPCPORT_1=8090
GRPCPORT_2=9090
GRPCWEB_1=8091
GRPCWEB_2=9091
GRPC_WEB_1_ENABLE=true
GRPC_WEB_2_ENABLE=true



echo "Starting $CHAINID_1 in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID_1.log"
$BINARY start --log_format json --home $CHAIN_DIR/$CHAINID_1 --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT_1" --grpc-web.enable=$GRPC_WEB_1_ENABLE > $CHAIN_DIR/$CHAINID_1.log 2>&1 &

echo "Starting $CHAINID_2 in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID_2.log"
$BINARY start --log_format json --home $CHAIN_DIR/$CHAINID_2 --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT_2"  --grpc-web.enable=$GRPC_WEB_2_ENABLE > $CHAIN_DIR/$CHAINID_2.log 2>&1 &
