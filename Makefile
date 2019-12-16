#!/usr/bin/make -f

SUBDIRS := $(shell find . * -mindepth 1 -maxdepth 1 -type d | grep -v "\.")

dirs := $(shell ls)

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*"  | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*"  | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*"  | xargs goimports -w -local github.com/irisnet/modules

test-unit:
	@for dir in $(SUBDIRS);  \
		do  \
			cd $$dir; \
			go test ./...; \
			cd ../../ ; \
		done