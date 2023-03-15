# Configurable env vars
BUILD_NUMBER ?= dev
REGISTRY ?= ghcr.io/loshz/platform
DOCKER ?= sudo docker
BIN_DIR ?= ${CURDIR}/bin
GO_TEST_FLAGS ?= -failfast
PROTOC_VERSION ?= 3.21.12

.PHONY: docker/build docker/push go/build go/lint go/test proto/check proto/install proto/build

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

proto/check:
	@protoc --version | grep $(PROTOC_VERSION) || (echo "Must use libprotoc $(PROTOC_VERSION)"; exit 1)

proto/install:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.29
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

proto/build:
	@protoc --go_out=pkg/pb/v1 --go_opt=module=github.com/loshz/platform/pkg/pb/v1 \
		--go-grpc_out=pkg/pb/v1 --go-grpc_opt=module=github.com/loshz/platform/pkg/pb/v1 \
		./proto/v1/*.proto
