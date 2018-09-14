all: get_vendor_deps install

DEP_CHECK := $(shell command -v dep 2> /dev/null)
STATIK_CHECK := $(shell command -v statik 2> /dev/null)

get_tools:
ifdef DEP_CHECK
	@echo "Dep is already installed.  Run 'make update_tools' to update."
else
	@echo "Installing dep"
	go get -v $(DEP)
endif
ifdef STATIK_CHECK
	@echo "Statik is already installed.  Run 'make update_tools' to update."
else
	@echo "Installing statik"
	go version
	go get -v $(STATIK)
endif

update_gaia_lite_docs:
	@statik -src=client/lcd/swaggerui -dest=client/lcd

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

install: update_gaia_lite_docs
	go install ./cmd/iris
	go install ./cmd/iriscli
	go install ./cmd/irislcd

build_linux: update_gaia_lite_docs
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iris ./cmd/iris && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iriscli ./cmd/iriscli \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/irislcd ./cmd/irislcd

build_cur: update_gaia_lite_docs
	go build -o build/iris ./cmd/iris  && \
	go build -o build/iriscli ./cmd/iriscli

build_example: update_gaia_lite_docs
	go build  -o build/iris1 ./examples/irishub1/cmd/iris1
	go build  -o build/iriscli1 ./examples/irishub1/cmd/iriscli1
	go build  -o build/iris2 ./examples/irishub2/cmd/iris2
	go build  -o build/iriscli2 ./examples/irishub2/cmd/iriscli2

install_examples: update_gaia_lite_docs
	go install ./examples/irishub1/cmd/iris1
	go install ./examples/irishub1/cmd/iriscli1
	go install ./examples/irishub2/cmd/iris2
	go install ./examples/irishub2/cmd/iriscli2

build_example_linux: update_gaia_lite_docs
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris1 ./examples/irishub1/cmd/iris1
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli1 ./examples/irishub1/cmd/iriscli1
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris2 ./examples/irishub2/cmd/iris2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli2 ./examples/irishub2/cmd/iriscli2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris2-bugfix ./examples/irishub-bugfix-2/cmd/iris-bugfix-2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli2-bugfix ./examples/irishub-bugfix-2/cmd/iriscli-bugfix-2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris3-bugfix ./examples/irishub-bugfix-3/cmd/iris-bugfix-3
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli3-bugfix ./examples/irishub-bugfix-3/cmd/iriscli-bugfix-3