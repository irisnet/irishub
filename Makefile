#!/usr/bin/make -f

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*"  | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*"  | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*"  | xargs goimports -w -local github.com/irisnet/modules