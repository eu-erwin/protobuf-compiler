FROM node:18.20.4 AS ts_builder
WORKDIR /code
RUN npm install -g ts-protoc-gen \
    protoc-gen-js

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

# Separated both copy, for cached layer
COPY --from=ts_builder /usr/local/bin/node \
    /usr/bin/
COPY --from=ts_builder /usr/local/lib/node_modules \
    /usr/local/lib/node_modules

RUN ln -s /usr/local/lib/node_modules/ts-protoc-gen/bin/protoc-gen-ts /usr/bin/protoc-gen-ts && \
    ln -s /usr/local/lib/node_modules/protoc-gen-js/bin/protoc-gen-js /usr/bin/protoc-gen-js && \
    ln -s /home/$UNAME/.pub-cache/bin/protoc-gen-dart /usr/bin/protoc-gen-dart

    go get github.com/googleapis/googleapis && \
    go get github.com/mwitkow/go-proto-validators && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.14.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.14.0 && \
    go install github.com/favadi/protoc-go-inject-tag@v1.4.0 && \
    go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
