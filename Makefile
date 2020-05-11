.PHONY: build install install-for-mac format

build: format
	@go build

install:
	@go install

install-for-mac:
	@go install -tags mac

format:
	@goimports -w .
