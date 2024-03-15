#################
# Build stage 0 #
#################
FROM golang:1.22-alpine3.19

# Install build dependencies
RUN apk --no-cache add build-base

# Create work dir and copy sourcecode.
WORKDIR /go/src/github.com/loshz/platform
COPY . .

ARG BUILD_NUMBER

# Build binaries
RUN BUILD_NUMBER=${BUILD_NUMBER} BIN_DIR=/go/bin/platform/ make go/build

#################
# Build stage 1 #
#################
FROM alpine:3.19

ARG USER=platform

# Install runtime dependencies
RUN apk --no-cache add ca-certificates \
  && apk --no-cache add curl

# Create group (-g/-G) and system user (-S)
# with specific UID (-u) and no password (-D)
RUN addgroup -g 2000 -S $USER \
  && adduser -G $USER -S -D -u 2000 $USER

# Copy binaries from build stage 0
COPY --from=0 /go/bin/platform/ /usr/local/bin/

# Copy TLS certs/keys
COPY --chown=$USER ./config/tls/*.crt.pem /usr/local/share/ca-certificates/
COPY --chown=$USER ./config/tls/*.key.pem /usr/local/share/ca-certificates/

WORKDIR /home/$USER

USER $USER
