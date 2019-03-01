.PHONY: install

install:
	@goimports -w cmd main.go
	@go install
