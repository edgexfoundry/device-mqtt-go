.PHONY: build test unittest lint clean prepare update docker

GO=CGO_ENABLED=0 GO111MODULE=on go
GOCGO=GCO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-mqtt

ARCH=$(shell uname -m)

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_mqtt_go

.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)

GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-mqtt-go.Version=$(VERSION)"

tidy:
	go mod tidy

build: $(MICROSERVICES)

cmd/device-mqtt:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd

unittest:
	$(GOCGO) test ./... -coverprofile=coverage.out ./...


lint:
	@which golangci-lint >/dev/null || echo "WARNING: go linter not installed. To install, run\n  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.42.1"
	@if [ "z${ARCH}" = "zx86_64" ] && which golangci-lint >/dev/null ; then golangci-lint run --config .golangci.yml ; else echo "WARNING: Linting skipped (not on x86_64 or linter not installed)"; fi

test: unittest lint
	$(GOCGO) vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]
	./bin/test-attribution-txt.sh

clean:
	rm -f $(MICROSERVICES)

run:
	cd bin && ./edgex-launch.sh

docker: $(DOCKERS)

docker_device_mqtt_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/device-mqtt:$(GIT_SHA) \
		-t edgexfoundry/device-mqtt:$(VERSION)-dev \
		.

vendor:
	$(GO) mod vendor
