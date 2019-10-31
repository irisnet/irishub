module github.com/irisnet/irishub

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973 // indirect
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.0.0-20181013004428-67e573d211ac
	github.com/btcsuite/btcutil v0.0.0-20180524032703-d4cc87b86016 // indirect
	github.com/cosmos/go-bip39 v0.0.0-20180618194314-52158e4697b8
	github.com/emicklei/proto v1.6.5
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-kit/kit v0.6.0
	github.com/go-logfmt/logfmt v0.3.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.1.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/protobuf v1.2.0
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/gorilla/websocket v1.2.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmhodges/levigo v0.0.0-20161115193449-c42d9e0ca023 // indirect
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/magiconair/properties v1.8.0 // indirect
	github.com/mattn/go-isatty v0.0.4
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.0.0
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/onsi/ginkgo v1.10.2 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.9.1
	github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910 // indirect
	github.com/prometheus/common v0.0.0-20181020173914-7e9e6cabbd39 // indirect
	github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d // indirect
	github.com/rakyll/statik v0.1.6
	github.com/rcrowley/go-metrics v0.0.0-20180503174638-e2704e165165 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/cobra v0.0.1
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.0.0
	github.com/stretchr/testify v1.2.2
	github.com/syndtr/goleveldb v0.0.0-20180708030551-c4c61651e9e3
	github.com/tendermint/btcd v0.1.0
	github.com/tendermint/go-amino v0.14.0
	github.com/tendermint/iavl v0.12.1
	github.com/tendermint/tendermint v0.31.0
	github.com/tendermint/tmlibs v0.9.0
	github.com/zondax/hid v0.9.0 // indirect
	github.com/zondax/ledger-cosmos-go v0.9.7
	github.com/zondax/ledger-go v0.8.0 // indirect
	golang.org/x/crypto v0.0.0-20180904163835-0709b304e793
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20181031143558-9b800f95dbbc // indirect
	google.golang.org/genproto v0.0.0-20180808183934-383e8b2c3b9e // indirect
	google.golang.org/grpc v1.13.0 // indirect
)

replace (
	github.com/tendermint/iavl => github.com/irisnet/iavl v0.12.2
	github.com/tendermint/tendermint => github.com/irisnet/tendermint v0.31.3
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
)

go 1.13
