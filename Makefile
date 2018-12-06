.PHONY: build test clean prepare update docker

GO=CGO_ENABLED=0 go

MICROSERVICES=cmd/device-mqtt

.PHONY: $(MICROSERVICES)

DOCKERS=docker_device_mqtt_go

.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION)
GIT_SHA=$(shell git rev-parse HEAD)

GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-mqtt-go.Version=$(VERSION)"

build: $(MICROSERVICES)
	go build ./...

cmd/device-mqtt:
	$(GO) build $(GOFLAGS) -o $@ ./cmd

test:
	go test ./... -cover

clean:
	rm -f $(MICROSERVICES)

prepare:
	glide install

update:
	glide update

run:
	cd bin && ./edgex-launch.sh

docker: $(DOCKERS)

docker_device_mqtt_go:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/docker-device-mqtt-go:$(GIT_SHA) \
		-t edgexfoundry/docker-device-mqtt-go:$(VERSION)-dev \
		.
