# Build config.
BUILD_NUMBER ?= dev
BIN_DIR ?= ${CURDIR}/bin/
GOARCH ?= amd64
GOOS ?= linux
CGO ?= 0
GO_TEST_FLAGS ?= -failfast -race
PROTOC_VERSION ?= 26.0

# Docker config.
DOCKER ?= sudo docker
DOCKER_IMAGE ?= loshz/platform

# TLS config.
TLS_CERT_DIR ?= ./config/tls

.PHONY: docker/build docker/compose go/build go/lint go/test proto/install proto/lint proto/build tls tls/ca tls/certs

docker/build:
	$(DOCKER) build \
	  --build-arg BUILD_NUMBER=$(BUILD_NUMBER) \
	  --tag $(DOCKER_IMAGE):$(BUILD_NUMBER) .

docker/compose:
	$(DOCKER) compose build --build-arg BUILD_NUMBER=$(BUILD_NUMBER)
	$(DOCKER) compose up

go/build: ./cmd/*
	@for CMD in $^; do \
		CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_DIR) \
		  --ldflags="-X github.com/loshz/platform/internal/version.Build=$(BUILD_NUMBER)" ./$${CMD}; \
	done

go/lint:
	@golangci-lint run --config .golangci.yml

go/test:
	@go test $(GO_TEST_FLAGS) ./...

proto/install:
	@curl -sL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip -o /tmp/protoc.zip
	@sudo unzip -o /tmp/protoc.zip -d /usr/local bin/protoc
	@sudo unzip -o /tmp/protoc.zip -d /usr/local 'include/*'
	@go install github.com/bufbuild/buf/cmd/buf@v1.30.0
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

proto/lint:
	@buf format --diff --exit-code
	@buf lint

proto/build: proto/lint
	@protoc --go_out=internal/api/v1 --go_opt=module=github.com/loshz/platform/internal/api/v1 \
		--go-grpc_out=internal/api/v1 --go-grpc_opt=module=github.com/loshz/platform/internal/api/v1 \
		./proto/v1/*.proto

tls:
	@mkdir $(TLS_CERT_DIR)
	$(MAKE) tls/ca
	$(MAKE) tls/certs

tls/ca:
	@openssl genpkey -algorithm ED25519 -out $(TLS_CERT_DIR)/ca.key.pem
	@openssl req -nodes -new -sha256 -x509 -key $(TLS_CERT_DIR)/ca.key.pem -out $(TLS_CERT_DIR)/ca.crt.pem \
		-subj "/O=Platform/CN=localhost" \
		-addext "subjectAltName = DNS:localhost,IP:0.0.0.0"

tls/certs:
	@echo "Generating server certs..."
	@openssl genpkey -algorithm ED25519 -out $(TLS_CERT_DIR)/server.key.pem
	@openssl req -nodes -new -sha256 -key $(TLS_CERT_DIR)/server.key.pem -out $(TLS_CERT_DIR)/server.csr.pem \
		-subj "/O=Platform/CN=localhost" \
		-addext "subjectAltName = DNS:localhost,IP:0.0.0.0"
	@openssl x509 -req -sha256 -in $(TLS_CERT_DIR)/server.csr.pem \
		-CA $(TLS_CERT_DIR)/ca.crt.pem -CAkey $(TLS_CERT_DIR)/ca.key.pem -CAcreateserial \
		-out $(TLS_CERT_DIR)/server.crt.pem
	@echo "Generating client certs..."
	@openssl genpkey -algorithm ED25519 -out $(TLS_CERT_DIR)/client.key.pem
	@openssl req -nodes -new -sha256 -key $(TLS_CERT_DIR)/client.key.pem -out $(TLS_CERT_DIR)/client.csr.pem \
		-subj "/O=Platform/CN=localhost" \
		-addext "subjectAltName = DNS:localhost,IP:0.0.0.0"
	@openssl x509 -req -sha256 -in $(TLS_CERT_DIR)/client.csr.pem \
		-CA $(TLS_CERT_DIR)/ca.crt.pem -CAkey $(TLS_CERT_DIR)/ca.key.pem -CAcreateserial \
		-out $(TLS_CERT_DIR)/client.crt.pem
