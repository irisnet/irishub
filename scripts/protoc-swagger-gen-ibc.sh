
IBC_GO=v7.2.0

go mod download github.com/cosmos/ibc-go/v7@${IBC_GO}

IBC_PATH=${GOPATH}/pkg/mod/github.com/cosmos/ibc-go/v7@${IBC_GO}

proto_dirs=$(find ${IBC_PATH}/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf generate --template buf.gen.swagger.yaml $query_file
  fi
done