BUILD_NUMBER ?= dev
DOCKER ?= sudo docker
DOCKER_IMAGE ?= loshz/platform
BIN_DIR ?= ${CURDIR}/bin
GO_TEST_FLAGS ?= -failfast -race
PROTOC_VERSION ?= 3.21.12

.PHONY: docker/build docker/compose go/build go/lint go/test proto/check proto/install proto/build

docker/build:
	$(DOCKER) build \
	  --build-arg BUILD_NUMBER=$(BUILD_NUMBER) \
	  --tag $(DOCKER_IMAGE):$(BUILD_NUMBER) .

docker/compose:
	$(DOCKER) compose build --build-arg BUILD_NUMBER=$(BUILD_NUMBER)
	$(DOCKER) compose up

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

proto/build: proto/check
	@protoc --go_out=pkg/pb/v1 --go_opt=module=github.com/loshz/platform/pkg/pb/v1 \
		--go-grpc_out=pkg/pb/v1 --go-grpc_opt=module=github.com/loshz/platform/pkg/pb/v1 \
		./proto/v1/*.proto
