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

ADD cmd /code/cmd
COPY go.mod go.sum main.go /code/

RUN go build \
    -a -ldflags "-s -w -X main.appCommit=$REVISION -X main.appVersion=$VERSION -X main.appEnv=prod" \
    -o helper \
    main.go

FROM gruebel/upx:latest AS upx
WORKDIR /code
COPY --from=go_builder /code/helper /helper
# Compress the binary and copy it to final image
RUN upx --best --lzma -o /app /helper

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
    org.opencontainers.image.title="Protobuf compiler" \
    org.opencontainers.image.url="erwineu/protobuf-compiler" \
    org.opencontainers.image.created=$TIMESTAMP \
    org.opencontainers.image.revision=$REVISION \
    org.opencontainers.image.version=$VERSION \
    version=$VERSION \
    description="Docker image based for compiling protobuf in four languages golang, typescript, dart and php"

RUN wget -qO- https://dl-ssl.google.com/linux/linux_signing_key.pub | gpg --dearmor -o /usr/share/keyrings/dart.gpg && \
    echo 'deb [signed-by=/usr/share/keyrings/dart.gpg arch=amd64] https://storage.googleapis.com/download.dartlang.org/linux/debian stable main' | tee /etc/apt/sources.list.d/dart_stable.list
RUN apt-get update && \
    apt-get -y install protobuf-compiler ^dart

WORKDIR /code

ARG UNAME=${UNAME:-compiler}
ARG UID=${UID:-1000}
ARG GID=${GID:-1000}

COPY compiler.sh /usr/local/bin/compiler
COPY script.sh /usr/local/bin/compile

COPY --from=go_builder /go/bin/protoc-gen-go \
    /go/bin/protoc-gen-go-grpc \
    /go/bin/protoc-gen-grpc-gateway \
    /go/bin/protoc-gen-openapiv2 \
    /go/bin/protoc-go-inject-tag \
    /go/bin/

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

RUN chmod +x /usr/local/bin/compiler && \
    chmod +x /usr/local/bin/compile

RUN ssh-keyscan -t rsa github.com > ~/.ssh/known_hosts && \
    go env -w GO111MODULE="on" && \
    go mod init example.com/compiler && \
    go get github.com/google/protobuf@v5.28.2+incompatible && \
    go get github.com/googleapis/googleapis && \
    go get github.com/mwitkow/go-proto-validators && \
    go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

RUN groupadd -g $GID -o $UNAME && \
    useradd -m -d /home/$UNAME -u $UID -g $GID -o -s /bin/bash $UNAME

RUN mkdir -p -m 0600 /home/$UNAME/.ssh && \
    chown $UNAME:$UNAME -R /home/$UNAME /home/$UNAME/.ssh /go /code

USER "$UNAME"

RUN dart pub global activate protoc_plugin

ENV HOME=/home/$UNAME
WORKDIR /code

ENV PATH=$PATH:/usr/local/go/bin
ENV GO111MODULE="on"

COPY --from=upx /app /bin/helper

ENTRYPOINT ["compiler"]
CMD ["compile"]
