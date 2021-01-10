#!/usr/bin/env bash

set -eo pipefail

SDK_VERSION=v0.40.0-rc5
IRISMOD_VERSION=v1.1.1-0.20201229063925-7d7dad20f951
WASMD_VERSION=v0.13.1-0.20201217131318-53bbf96e9e87
WASMD_PROTO_DIR=x/wasm/internal/types

chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto/cosmos
chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto/ibc
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/!cosm!wasm/wasmd@${WASMD_VERSION}/${WASMD_PROTO_DIR}

mkdir -p ./proto/wasm

cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto/cosmos ./proto
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto/ibc ./proto
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto ./
cp -r ${GOPATH}/pkg/mod/github.com/!cosm!wasm/wasmd@${WASMD_VERSION}/${WASMD_PROTO_DIR}/*.proto ./proto/wasm

sed -i "" "s@${WASMD_PROTO_DIR}@wasm@g" `find ./proto/wasm -type file`

mkdir -p ./tmp-swagger-gen

proto_dirs=$(find ./proto -path './proto/cosmos/base/tendermint*' -prune -o -name 'query.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do

    # generate swagger files (filter query files)
    query_file=$(find "${dir}" -maxdepth 1 -name 'query.proto')
    echo $query_file
    if [[ ! -z "$query_file" ]]; then
        protoc \
            -I "proto" \
            -I "third_party/proto" \
            "$query_file" \
            --swagger_out ./tmp-swagger-gen \
            --swagger_opt logtostderr=true --swagger_opt fqn_for_swagger_name=true --swagger_opt simple_operation_ids=true
    fi
done

# copy cosmos swagger_legacy.yaml
chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/client/docs/swagger_legacy.yaml
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/client/docs/swagger_legacy.yaml ./lite/cosmos_swagger_legacy.yaml

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./lite/config.json -o ./lite/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# replace APIs example
sed -r -i '' 's/cosmos1[a-z,0-9]+/iaa1sltcyjm5k0edlg59t47lsyw8gtgc3nudklntcq/g' ./lite/swagger-ui/swagger.yaml
sed -r -i '' 's/cosmosvaloper1[a-z,0-9]+/iva1sltcyjm5k0edlg59t47lsyw8gtgc3nudrwey98/g' ./lite/swagger-ui/swagger.yaml
sed -r -i '' 's/cosmosvalconspub1[a-z,0-9]+/icp1zcjduepqwhwqn4h5v6mqa7k3kmy7cjzchsx5ptsrqaulwrgfmghy3k9jtdzs6rdddm/g' ./lite/swagger-ui/swagger.yaml
sed -i '' 's/Gaia/IRISHub/g' ./lite/swagger-ui/swagger.yaml
sed -i '' 's/gaia/irishub/g' ./lite/swagger-ui/swagger.yaml
sed -i '' 's/cosmoshub/irishub/g' ./lite/swagger-ui/swagger.yaml
 
# clean swagger files
rm -rf ./tmp-swagger-gen

# clean proto files
rm -rf ./proto/cosmos
rm -rf ./proto/ibc

rm -rf ./proto/coinswap
rm -rf ./proto/htlc
rm -rf ./proto/nft
rm -rf ./proto/oracle
rm -rf ./proto/random
rm -rf ./proto/record
rm -rf ./proto/service
rm -rf ./proto/token
rm -fr ./proto/wasm
