#!/usr/bin/env bash

set -eo pipefail

rm -rf ./tmp-swagger-gen ./tmp && mkdir -p ./tmp-swagger-gen ./tmp/proto ./tmp/third_party

SDK_VERSION=v0.47.4
IRISMOD_VERSION=v1.7.4-0.20231020012541-015d2c01fd9e
IBC_GO=v7.3.0
EVM_VERSION=v0.22.0

go mod download github.com/cosmos/cosmos-sdk@${SDK_VERSION}
go mod download github.com/irisnet/irismod@${IRISMOD_VERSION}
go mod download github.com/cosmos/ibc-go/v7@${IBC_GO}
go mod download github.com/evmos/ethermint@${EVM_VERSION}

chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/ibc-go/v7@${IBC_GO}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/evmos/ethermint@${EVM_VERSION}/proto

cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto/amino ./tmp/proto/amino && rm -rf ./tmp/proto/cosmos/mint
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto/cosmos ./tmp/proto/cosmos
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto/tendermint ./tmp/proto/tendermint
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto/irismod ./tmp/proto/irismod
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/ibc-go/v7@${IBC_GO}/proto/ibc ./tmp/proto/ibc
cp -r ${GOPATH}/pkg/mod/github.com/evmos/ethermint@${EVM_VERSION}/proto/ethermint ./tmp/proto/ethermint
cp -r ./proto/* ./tmp/proto/
cp buf.work.yaml ./tmp/

cd ./tmp
proto_dirs=$(find ./proto  -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
    # generate swagger files (filter query files)
    query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
    if [[ -n "$query_file" ]]; then
        buf generate --template proto/buf.gen.swagger.yaml "$query_file"
    fi
done
cd ../

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./lite/config.json -o ./lite/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# replace APIs example
sed -r -i 's/cosmos1[a-z,0-9]+/iaa1sltcyjm5k0edlg59t47lsyw8gtgc3nudklntcq/g' ./lite/swagger-ui/swagger.yaml
sed -r -i 's/cosmosvaloper1[a-z,0-9]+/iva1sltcyjm5k0edlg59t47lsyw8gtgc3nudrwey98/g' ./lite/swagger-ui/swagger.yaml
sed -r -i 's/cosmosvalconspub1[a-z,0-9]+/icp1zcjduepqwhwqn4h5v6mqa7k3kmy7cjzchsx5ptsrqaulwrgfmghy3k9jtdzs6rdddm/g' ./lite/swagger-ui/swagger.yaml
sed -i 's/Gaia/IRIShub/g' ./lite/swagger-ui/swagger.yaml
sed -i 's/gaia/irishub/g' ./lite/swagger-ui/swagger.yaml
sed -i 's/cosmoshub/irishub/g' ./lite/swagger-ui/swagger.yaml

# clean swagger files
rm -rf ./tmp-swagger-gen
rm -rf ./github.com
rm -rf ./cosmos

# clean proto files
rm -rf ./tmp
