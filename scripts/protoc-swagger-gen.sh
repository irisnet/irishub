#!/usr/bin/env bash

set -eo pipefail

IRISMOD_VERSION=v1.1.1-0.20201015063111-aa6725ffb200
SDK_VERSION=v0.34.4-0.20201014023301-f172e47973d0

chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/cosmos
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/ibc

cp -r ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto ./
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/cosmos ./proto
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/ibc ./proto

mkdir -p ./tmp-swagger-gen

proto_dirs=$(find ./proto -path -prune -o -name 'query.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
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

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./lite/config.json -o ./lite/swagger-ui/swagger.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

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
