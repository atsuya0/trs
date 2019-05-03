.PHONY: install format

install: format
	@go install

format:
	@goimports -w .
