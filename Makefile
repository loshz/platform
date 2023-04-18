# Build config.
BUILD_NUMBER ?= dev
BIN_DIR ?= ${CURDIR}/bin
GO_TEST_FLAGS ?= -failfast -race
PROTOC_VERSION ?= 3.21.12

# Docker config.
DOCKER ?= sudo docker
DOCKER_IMAGE ?= loshz/platform

# TLS config.
TLS_CERT_DIR ?= ./config/certs
TLS_SUBJ ?= /O=Platform/CN=localhost

.PHONY: docker/build docker/compose go/build go/lint go/test proto/check proto/install proto/build tls/ca tls/certs

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
	@golangci-lint run --config .golangci.yml

go/test:
	@go test $(GO_TEST_FLAGS) ./...

proto/check:
	@protoc --version | grep $(PROTOC_VERSION) || (echo "Must use libprotoc $(PROTOC_VERSION)"; exit 1)

proto/install:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

proto/build: proto/check
	@echo "Generating proto..."
	@protoc --go_out=pkg/pb/v1 --go_opt=module=github.com/loshz/platform/pkg/pb/v1 \
		--go-grpc_out=pkg/pb/v1 --go-grpc_opt=module=github.com/loshz/platform/pkg/pb/v1 \
		./proto/v1/*.proto

tls/ca:
	@openssl genpkey -algorithm ED25519 -out $(TLS_CERT_DIR)/ca.key.pem
	@openssl req -nodes -new -sha256 -x509 -key $(TLS_CERT_DIR)/ca.key.pem -out $(TLS_CERT_DIR)/ca.crt.pem \
		-subj "$(TLS_SUBJ)" \
		-addext "subjectAltName = DNS:localhost,IP:0.0.0.0"

tls/certs:
	@echo "Generating server certs..."
	@openssl genpkey -algorithm ED25519 -out $(TLS_CERT_DIR)/server.key.pem
	@openssl req -nodes -new -sha256 -key $(TLS_CERT_DIR)/server.key.pem -out $(TLS_CERT_DIR)/server.csr.pem \
		-subj "$(TLS_SUBJ)" \
		-addext "subjectAltName = DNS:localhost,IP:0.0.0.0"
	@openssl x509 -req -sha256 -in $(TLS_CERT_DIR)/server.csr.pem \
		-CA $(TLS_CERT_DIR)/ca.crt.pem -CAkey $(TLS_CERT_DIR)/ca.key.pem -CAcreateserial \
		-out $(TLS_CERT_DIR)/server.crt.pem
	@echo "Generating client certs..."
	@openssl genpkey -algorithm ED25519 -out $(TLS_CERT_DIR)/client.key.pem
	@openssl req -nodes -new -sha256 -key $(TLS_CERT_DIR)/client.key.pem -out $(TLS_CERT_DIR)/client.csr.pem \
		-subj "$(TLS_SUBJ)" \
		-addext "subjectAltName = DNS:localhost,IP:0.0.0.0"
	@openssl x509 -req -sha256 -in $(TLS_CERT_DIR)/client.csr.pem \
		-CA $(TLS_CERT_DIR)/ca.crt.pem -CAkey $(TLS_CERT_DIR)/ca.key.pem -CAcreateserial \
		-out $(TLS_CERT_DIR)/client.crt.pem
