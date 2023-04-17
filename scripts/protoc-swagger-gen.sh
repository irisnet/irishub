#!/usr/bin/env bash

set -eo pipefail

rm -rf ./tmp-swagger-gen ./tmp && mkdir -p ./tmp-swagger-gen ./tmp/proto ./tmp/third_party

chmod a+x ./scripts/protoc-swagger-gen-ibc.sh
chmod a+x ./scripts/protoc-swagger-gen-evm.sh
./scripts/protoc-swagger-gen-ibc.sh
./scripts/protoc-swagger-gen-evm.sh

SDK_VERSION=v0.46.9
IRISMOD_VERSION=v1.7.3

go mod download github.com/cosmos/cosmos-sdk@${SDK_VERSION}
go mod download github.com/irisnet/irismod@${IRISMOD_VERSION}

chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto

cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto ./tmp && rm -rf ./tmp/proto/cosmos/mint
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto ./tmp
cp -r ./proto ./tmp

proto_dirs=$(find ./tmp/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
    # generate swagger files (filter query files)
    query_file=$(find "${dir}" -maxdepth 1 -name 'query.proto')
    if [[ $dir =~ "cosmos" ]]; then
        query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
    fi
    if [[ $dir =~ "ibc" ]]; then
        query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
    fi
    if [[ ! -z "$query_file" ]]; then
        buf generate --template buf.gen.swagger.yaml $query_file
    fi
done

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

# TODO
# generate proto doc
# buf protoc \
#     -I "tmp/proto" \
#     -I "third_party/proto" \
#     --doc_out=./docs/endpoints \
#     --doc_opt=./docs/endpoints/protodoc-markdown.tmpl,proto-docs.md \
#     $(find "$(pwd)/tmp/proto" -maxdepth 5 -name '*.proto')
# cp ./docs/endpoints/proto-docs.md ./docs/zh/endpoints/proto-docs.md

# clean swagger files
rm -rf ./tmp-swagger-gen
rm -rf ./github.com
rm -rf ./cosmos

# clean proto files
rm -rf ./tmp
