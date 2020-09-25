#!/usr/bin/env bash

set -eo pipefail

IRISMOD_VERSION=v0.0.0-20200925031428-cad7dc2a03fa
SDK_VERSION=v0.34.4-0.20200918054421-c8b3462ab7a2

chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/cosmos
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/ibc

cp -r ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto ./
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/cosmos ./proto
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/cosmos-sdk@${SDK_VERSION}/proto/ibc ./proto

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
            --swagger_out=logtostderr=true,stderrthreshold=1000,fqn_for_swagger_name=true,simple_operation_ids=true:.
    fi
done

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
swagger-combine ./lite/grpc-gateway/config.json -o ./lite/grpc-gateway/swagger.json --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files & empty folder
find ./ -name 'query.swagger.json' -exec rm {} \;
find ./ -type d -empty | xargs -n 1 rm -rf
rm -r ./cosmos

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
