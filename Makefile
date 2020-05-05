.PHONY: build install install-for-mac format

build: format
	@go build

install: format
	@go install

install-for-mac: format
	@go install -tags mac

format:
	@goimports -w .
