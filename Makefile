.PHONY: all test vuln

GOTESTSUM_VERSION := v1.13.0
REVIVE_VERSION := v1.16.0
GOVULNCHECK_VERSION := v1.7.0

all: build

proto:
	@./buildall.sh

build:
	@go build ./...


# {{{ test

PACKAGES := $(shell go list ./... | grep -v "/examples/" | grep -v "/persistence" | grep -v "/scheduler")




test:
	@go test $(PACKAGES) -timeout=30s

test2:
	@go run gotest.tools/gotestsum@$(GOTESTSUM_VERSION) --format testname $(PACKAGES)

test-short:
	@go test $(PACKAGES) -timeout=30s -short

test-race:
	@go test $(PACKAGES) -timeout=2m -race

lint:
	@go run github.com/mgechev/revive@$(REVIVE_VERSION) -formatter friendly $(PACKAGES)

vet:
	@go vet $(PACKAGES)

vuln:
	@go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

bench:
	@go test $(PACKAGES) -run=^$$ -bench=.

# }}} test

# {{{ benchmark

packages_benchmark := $(shell go list ./... | grep -v "/log")

benchmark:
	go test -benchmem -run=^$$ $(packages_benchmark) -bench ^Benchmark$(t).*$$
# }}}

# {{{ docker-env
root_dir := $(abspath $(CURDIR)/)
docker-env:
	sudo docker run -it --rm \
		-v $(root_dir)/:/go/src/AsncronIT/protoactor-go \
		-w /go/src/AsncronIT/protoactor-go \
		-e GOPATH=/go \
		--entrypoint /bin/bash \
		cupen/protoc:3.9.1-1
# }}}
