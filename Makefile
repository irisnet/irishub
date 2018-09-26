all: get_vendor_deps install

COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -ldflags "-X github.com/irisnet/irishub/version.GitCommit=${COMMIT_HASH}"

########################################
### Tools & dependencies

check_tools:
	cd deps_tools && $(MAKE) check_tools

check_dev_tools:
	cd deps_tools && $(MAKE) check_dev_tools

update_tools:
	cd deps_tools && $(MAKE) update_tools

update_dev_tools:
	cd deps_tools && $(MAKE) update_dev_tools

get_tools:
	cd deps_tools && $(MAKE) get_tools

get_dev_tools:
	cd deps_tools && $(MAKE) get_dev_tools

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v

draw_deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i github.com/irisnet/irishub/cmd/iris -d 2 | dot -Tpng -o dependency-graph.png

########################################
### Generate swagger docs for irislcd
update_irislcd_swagger_docs:
	@statik -src=client/lcd/swaggerui -dest=client/lcd

########################################
### Compile and Install
install: update_irislcd_swagger_docs
	go install $(BUILD_FLAGS) ./cmd/iris
	go install $(BUILD_FLAGS) ./cmd/iriscli
	go install $(BUILD_FLAGS) ./cmd/irislcd

install_debug:
	go install ./cmd/irisdebug

build_linux: update_irislcd_swagger_docs
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iris ./cmd/iris && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iriscli ./cmd/iriscli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/irislcd ./cmd/irislcd

build_windows: update_irislcd_swagger_docs
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/iris.exe ./cmd/iris && \
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/iriscli.exe ./cmd/iriscli && \
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/irislcd.exe ./cmd/irislcd

build_cur: update_irislcd_swagger_docs
	go build -o build/iris ./cmd/iris  && \
	go build -o build/iriscli ./cmd/iriscli && \
	go build -o build/irislcd ./cmd/irislcd

build_example: update_irislcd_swagger_docs
	go build  -o build/iris1 ./examples/irishub1/cmd/iris1
	go build  -o build/iriscli1 ./examples/irishub1/cmd/iriscli1
	go build  -o build/iris2 ./examples/irishub2/cmd/iris2
	go build  -o build/iriscli2 ./examples/irishub2/cmd/iriscli2

install_examples: update_irislcd_swagger_docs
	go install ./examples/irishub1/cmd/iris1
	go install ./examples/irishub1/cmd/iriscli1
	go install ./examples/irishub2/cmd/iris2
	go install ./examples/irishub2/cmd/iriscli2

build_example_linux: update_irislcd_swagger_docs
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris1 ./examples/irishub1/cmd/iris1
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli1 ./examples/irishub1/cmd/iriscli1
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris2 ./examples/irishub2/cmd/iris2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli2 ./examples/irishub2/cmd/iriscli2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris2-bugfix ./examples/irishub-bugfix-2/cmd/iris-bugfix-2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli2-bugfix ./examples/irishub-bugfix-2/cmd/iriscli-bugfix-2
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iris3-bugfix ./examples/irishub-bugfix-3/cmd/iris-bugfix-3
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o build/iriscli3-bugfix ./examples/irishub-bugfix-3/cmd/iriscli-bugfix-3