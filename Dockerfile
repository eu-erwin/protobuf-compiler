ARG GO_VERSION=${GO_VERSION:-latest}
ARG VERSION=${VERSION:-latest}
ARG REVISION=${REVISION:-latest}
FROM golang:$GO_VERSION AS go_builder
WORKDIR /code
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.14.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.14.0 && \
    go install github.com/favadi/protoc-go-inject-tag@v1.4.0

COPY go.mod go.sum naming.go /code/
RUN go build \
    -a -ldflags "-s -w -X main.appCommit=$REVISION -X main.appVersion=$VERSION -X main.appEnv=prod" \
    -o naming \
    naming.go

FROM gruebel/upx:latest AS upx
WORKDIR /code
COPY --from=go_builder /code/naming /naming
# Compress the binary and copy it to final image
RUN upx --best --lzma -o /app /naming

FROM node:18.20.4 AS ts_builder
WORKDIR /code
RUN npm install -g ts-protoc-gen \
    protoc-gen-js

ARG GO_VERSION=${GO_VERSION:-latest}
FROM golang:$GO_VERSION AS php_builder
WORKDIR /code
RUN apt-get update && \
    apt-get install -y autoconf cmake && \
    git clone -b v1.34.1 --depth 1 https://github.com/grpc/grpc && \
    cd grpc && \
    git submodule update --init && \
    mkdir -p cmake/build && \
    cd cmake/build && \
    cmake ../.. && \
    make protoc grpc_php_plugin && \
    cp grpc_php_plugin /usr/bin/grpc_php_plugin && \
    rm -rf /tmp/grpc

ARG GO_VERSION=${GO_VERSION:-latest}
FROM golang:$GO_VERSION

LABEL org.opencontainers.image.authors="eu.erwin@gmx.de" \
    org.opencontainers.image.description="Docker image based for compiling protobuf in golang" \
    org.opencontainers.image.url="https://github.com/eu-erwin/protobuf-compiler" \
    org.opencontainers.image.documentation="https://github.com/eu-erwin/protobuf-compiler" \
    org.opencontainers.image.version="latest" \
    org.opencontainers.image.license="MIT" \
    org.opencontainers.image.revision=$REVISION \
    org.opencontainers.image.created=$TIMESTAMP

WORKDIR /code

RUN go mod init __MODULE__ && \
    go get github.com/google/protobuf@v4.23.1+incompatible && \
    go get github.com/googleapis/googleapis && \
    go get github.com/mwitkow/go-proto-validators && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.14.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.14.0 && \
    go install github.com/favadi/protoc-go-inject-tag@v1.4.0 && \
    go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

COPY --from=php_builder /usr/bin/grpc_php_plugin \
    /usr/bin/

# Separated both copy, for cached layer
COPY --from=ts_builder /usr/local/bin/node \
    /usr/bin/
COPY --from=ts_builder /usr/local/lib/node_modules \
    /usr/local/lib/node_modules

RUN ln -s /usr/local/lib/node_modules/ts-protoc-gen/bin/protoc-gen-ts /usr/bin/protoc-gen-ts && \
    ln -s /usr/local/lib/node_modules/protoc-gen-js/bin/protoc-gen-js /usr/bin/protoc-gen-js && \
    ln -s /home/$UNAME/.pub-cache/bin/protoc-gen-dart /usr/bin/protoc-gen-dart

RUN mkdir ~/.ssh && \
    mkdir /var/protobuf && \
    mkdir /var/protobuf/template && \
    chmod ugo+rw -R /var/protobuf

COPY template /var/protobuf/template/

COPY --from=upx /app /bin/naming
