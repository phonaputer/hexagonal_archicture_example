check: compilecheck lint

compilecheck:
	go build -o compilecheckbin cmd/exampleapp/main.go && \
	rm -f compilecheckbin

lint:
	golangci-lint run ./...
