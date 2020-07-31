module github.com/irisnet/irishub

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200721190130-5d71020270ae
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/irismod/coinswap v0.0.0-20200722055706-deeded9d99b8
	github.com/irismod/htlc v0.0.0-20200722060015-b71f49c9b167
	github.com/irismod/nft v1.1.1-0.20200722060344-38fec5db63a2
	github.com/irismod/record v1.1.1-0.20200731015314-a365015d8218
	github.com/irismod/service v1.1.1-0.20200723031529-6abecb02ceb1
	github.com/irismod/token v1.1.1-0.20200723031618-028bdd6fb30a
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/otiai10/copy v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.33.6
	github.com/tendermint/tm-db v0.5.1
	github.com/tidwall/gjson v1.6.0
	google.golang.org/grpc v1.30.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.19.1-0.20200722022502-e2d6c76ae750
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
)
