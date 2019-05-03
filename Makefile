.PHONY: install install-for-mac format

install: format
	@go install

install-for-mac: format
	@go install -tags mac

format:
	@goimports -w .
