all: get_vendor_deps install test

get_vendor_deps:
	go get github.com/Masterminds/glide
	glide install
	
install:
	go install ./cmd/iris

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/iris ./cmd/iris

test:
	@go test `glide novendor`

test_cli:
	bash ./cmd/iris/sh_tests/stake.sh
