#!/usr/bin/env bash

set -eo pipefail

SDK_VERSION=v0.42.9
IRISMOD_VERSION=v1.4.1-0.20210820100038-3ce2e8f7874d

chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/third_party/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto

rm -rf ./tmp-swagger-gen ./tmp && mkdir -p ./tmp-swagger-gen ./tmp/proto ./tmp/third_party

cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto ./tmp && rm -rf ./tmp/proto/cosmos/mint
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/third_party/proto ./tmp/third_party
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto ./tmp
cp -r ./proto ./tmp

proto_dirs=$(find ./tmp/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do

    # generate swagger files (filter query files)
    query_file=$(find "${dir}" -maxdepth 1 -name 'query.proto')
    if [[ $dir =~ "cosmos" ]]; then
        query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
    fi
    if [[ ! -z "$query_file" ]]; then
        protoc \
            -I "tmp/proto" \
            -I "tmp/third_party/proto" \
            "$query_file" \
            --swagger_out=./tmp-swagger-gen \
            --swagger_opt=logtostderr=true --swagger_opt=fqn_for_swagger_name=true --swagger_opt=simple_operation_ids=true
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
sed -i '' 's/Gaia/IRIShub/g' ./lite/swagger-ui/swagger.yaml
sed -i '' 's/gaia/irishub/g' ./lite/swagger-ui/swagger.yaml
sed -i '' 's/cosmoshub/irishub/g' ./lite/swagger-ui/swagger.yaml

# generate proto doc
buf protoc \
    -I "tmp/proto" \
    -I "tmp/third_party/proto" \
    --doc_out=./docs/endpoints \
    --doc_opt=./docs/endpoints/protodoc-markdown.tmpl,proto-docs.md \
    $(find "$(pwd)/tmp/proto" -maxdepth 5 -name '*.proto')
go mod tidy

cp ./docs/endpoints/proto-docs.md ./docs/zh/endpoints/proto-docs.md

# clean swagger files
rm -rf ./tmp-swagger-gen

# clean proto files
rm -rf ./tmp
