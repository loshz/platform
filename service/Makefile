# Configurable env vars
BUILD_NUMBER ?= dev
REGISTRY ?= ghcr.io/loshz/platform
DOCKER ?= sudo docker
GO_TEST_FLAGS ?= -v -failfast

.PHONY: docker/build docker/push go/test

docker/build:
	$(DOCKER) build \
	  --build-arg BUILD_NUMBER=$(BUILD_NUMBER) \
	  --tag $(REGISTRY):$(BUILD_NUMBER) .

docker/push:
	$(DOCKER) push $(REGISTRY):$(BUILD_NUMBER)

go/test:
	go test $(GO_TEST_FLAGS) ./...
