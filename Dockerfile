#################
# Build stage 0 #
#################
FROM golang:1.21-alpine3.18

# Install build dependencies
RUN apk --no-cache add build-base

# Create work dir and copy sourcecode.
WORKDIR /go/src/github.com/loshz/platform
COPY . .

ARG BUILD_NUMBER

# Build binaries
RUN BUILD_NUMBER=${BUILD_NUMBER} BIN_DIR=/go/bin/platform make go/build

#################
# Build stage 1 #
#################
FROM alpine:3.18

ARG USER=platform

# Create group (-g/-G) and system user (-S)
# with specific UID (-u) and no password (-D)
RUN addgroup -g 2000 -S $USER \
  && adduser -G $USER -S -D -u 2000 $USER

# Copy operator binary from build stage 0
COPY --from=0 --chown=$USER /go/bin/platform/ /bin/
COPY --chown=$USER ./config/certs /home/$USER/certs

WORKDIR /home/$USER

USER $USER
