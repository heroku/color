PROJECT_ROOT := $(shell pwd)

.PHONY: test bench

test:
	go test -race -coverprofile c.out $(PROJECT_ROOT)/...

bench:
	go test -bench=. -count=5  -cpu=1,2,4,8 -cpuprofile=cpu.prof github.com/murphybytes/color

cover: test
	go tool cover -html=c.out

lint: $(GOPATH)/bin/golangci-lint
	@echo "--> Running linter with default config"
	golangci-lint run --deadline 3m0s -c $(PROJECT_ROOT)/.golangcli.yml

$(GOPATH)/bin/golangci-lint:
	@echo "--> Installing linter"
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.18.0