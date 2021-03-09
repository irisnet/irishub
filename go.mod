module github.com/irisnet/irishub

go 1.15

require (
	github.com/cosmos/cosmos-sdk v0.41.4
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.4.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/irisnet/irismod v1.3.2-0.20210308021302-c22d58f6e250
	github.com/olebedev/config v0.0.0-20190528211619-364964f3a8e4
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/tm-db v0.6.4
	github.com/tidwall/gjson v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	google.golang.org/genproto v0.0.0-20210204154452-deb828366460
	google.golang.org/grpc v1.35.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.34.4-0.20210304082514-9a60267940b8
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)
