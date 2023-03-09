# Configurable env vars
BUILD_NUMBER ?= dev
REGISTRY ?= ghcr.io/loshz/platform
DOCKER ?= sudo docker
BIN_DIR ?= ${CURDIR}/bin
GO_TEST_FLAGS ?= -failfast

.PHONY: docker/build docker/push go/build go/lint go/test

docker/build:
	$(DOCKER) build \
	  --build-arg BUILD_NUMBER=$(BUILD_NUMBER) \
	  --tag $(REGISTRY):$(BUILD_NUMBER) .

docker/push:
	$(DOCKER) push $(REGISTRY):$(BUILD_NUMBER)

go/build: ./cmd/*
	@for CMD in $^; do \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOBIN=$(BIN_DIR) go install \
		  --ldflags="-X github.com/loshz/platform/pkg/version.Build=$(BUILD_NUMBER)" ./$${CMD}; \
	done

go/lint:
	@golangci-lint run

go/test:
	@go test $(GO_TEST_FLAGS) ./...
