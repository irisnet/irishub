module github.com/irisnet/irishub

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d
	github.com/cosmos/go-bip39 v0.0.0-20180618194314-52158e4697b8
	github.com/cosmos/ledger-cosmos-go v0.10.3
	github.com/emicklei/proto v1.6.5
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-kit/kit v0.6.0
	github.com/gogo/protobuf v1.1.1
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/mattn/go-isatty v0.0.4
	github.com/mitchellh/go-homedir v1.0.0
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.1
	github.com/rakyll/statik v0.1.6
	github.com/spf13/cobra v0.0.1
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.0.0
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.1-0.20190318030020-c3a204f8e965
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/iavl v0.12.1
	github.com/tendermint/tendermint v0.31.0
	github.com/tendermint/tm-db v0.1.1
	github.com/tendermint/tmlibs v0.9.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
)

replace (
	github.com/tendermint/iavl => github.com/irisnet/iavl v0.8.2-0.20190905023710-abef11e7f66b
	github.com/tendermint/tendermint => github.com/irisnet/tendermint v0.28.1-0.20191030084935-cc25237a944f
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
)

go 1.13
