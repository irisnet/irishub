module github.com/irisnet/irishub

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200825201020-d9fd4d2ca9a3
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.7
	github.com/irismod/coinswap v0.0.0-20200827063026-8276ce7d5e75
	github.com/irismod/htlc v0.0.0-20200827065042-ae68dc82561a
	github.com/irismod/nft v1.1.1-0.20200827025546-49ab1f3ea314
	github.com/irismod/record v1.1.1-0.20200826120326-a9bb1ec595f4
	github.com/irismod/service v1.1.1-0.20200827072807-de292c89b86a
	github.com/irismod/token v1.1.1-0.20200827064712-8b4997bca4fe
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.0-rc3
	github.com/tendermint/tm-db v0.6.1
	github.com/tidwall/gjson v1.6.0
	google.golang.org/grpc v1.31.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.34.4-0.20200827085823-7b1a0adbfd1e
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
