#!/usr/bin/env bash

set -eo pipefail

SDK_VERSION=v0.41.0
IRISMOD_VERSION=v1.2.2-0.20210204090329-f93c3bea5aad

chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/third_party/proto
chmod -R 755 ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto

rm -rf tmp && mkdir -p tmp/proto tmp/third_party

cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/proto ./tmp && rm -rf ./tmp/proto/cosmos/mint
cp -r ${GOPATH}/pkg/mod/github.com/cosmos/cosmos-sdk@${SDK_VERSION}/third_party/proto ./tmp/third_party
cp -r ${GOPATH}/pkg/mod/github.com/irisnet/irismod@${IRISMOD_VERSION}/proto ./tmp
cp -r ./proto ./tmp

# command to generate docs using protoc-doc-gen
buf protoc \
    -I "tmp/proto" \
    -I "tmp/third_party/proto" \
    --doc_out=./docs/endpoints \
    --doc_opt=./docs/endpoints/protodoc-markdown.tmpl,proto-docs.md \
    $(find "$(pwd)/tmp/proto" -maxdepth 5 -name '*.proto')
go mod tidy

cp ./docs/endpoints/proto-docs.md ./docs/zh/endpoints/proto-docs.md

# clean proto files
rm -rf ./tmp
