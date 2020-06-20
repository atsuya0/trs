.PHONY: build install format

build: format
	@go build

install:
	@go install

format:
	@goimports -w .
