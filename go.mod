module github.com/irisnet/irishub

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200825201020-d9fd4d2ca9a3
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.7
	github.com/irismod/coinswap v0.0.0-20200901103745-f38672ac63ec
	github.com/irismod/htlc v0.0.0-20200901103718-c3bf89708dce
	github.com/irismod/nft v1.1.1-0.20200827095318-d16861212579
	github.com/irismod/record v1.1.1-0.20200827095301-3e27fc43ae73
	github.com/irismod/service v1.1.1-0.20200901115916-d898b826bf10
	github.com/irismod/token v1.1.1-0.20200901121217-d3aa04e760e3
	github.com/olebedev/config v0.0.0-20190528211619-364964f3a8e4
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
	github.com/cosmos/cosmos-sdk => github.com/irisnet/cosmos-sdk v0.34.4-0.20200901030027-1e0963031861
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
)
