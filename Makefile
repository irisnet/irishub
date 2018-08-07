all: get_vendor_deps install

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

install:
	go install ./cmd/iris
	go install ./cmd/iriscli

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iris ./cmd/iris && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iriscli ./cmd/iriscli

build_cur:
	go build -o build/iris ./cmd/iris  && \
	go build -o build/iriscli ./cmd/iriscli

build_example:
	go build  -o build/basecoind ./examples/basecoin/cmd/basecoind
	go build  -o build/basecli ./examples/basecoin/cmd/basecli
	go build  -o build/basecoind1 ./examples/basecoin1/cmd/basecoind1
	go build  -o build/basecli1 ./examples/basecoin1/cmd/basecli1

install_examples:
	go install ./examples/basecoin/cmd/basecoind
	go install ./examples/basecoin/cmd/basecli
	go install ./examples/basecoin1/cmd/basecoind1
	go install ./examples/basecoin1/cmd/basecli1