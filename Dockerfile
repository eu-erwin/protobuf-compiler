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
