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
