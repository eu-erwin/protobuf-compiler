#!/bin/bash

export PATH=$PATH:/usr/local/go/bin:/go/bin
export GO111MODULE="on"

go env -w GO111MODULE="on"

if [ -e "./env" ]; then
  source ./env
else
  if [ -z ${GENERIC_NAME+x} ]; then
    echo "env file is not provided and no env GENERIC_NAME set"
    exit 1;
  else
    CAMEL_CASE="${GENERIC_NAME}"
    LOWER_UNDER_CASE=$(helper naming -ls "${GENERIC_NAME}");
    GENERIC_DESC="Protobuf definition for ${GENERIC_NAME} service"
    LIBRARY_NAME="${LIBRARY_NAME:-protobuf/${GENERIC_NAME}}"
    LIBRARY_DESCRIPTION="${LIBRARY_DESCRIPTION:-${GENERIC_DESC}}"

    GO_PACKAGE_NAME="${GO_PACKAGE_NAME:-${GENERIC_NAME}}"

    PHP_LIBRARY_NAME="${PHP_LIBRARY_NAME:-protobuf\/${GENERIC_NAME}}"
    PHP_LIBRARY_DESCRIPTION="${PHP_LIBRARY_DESCRIPTION:-${GENERIC_DESC}}"
    PHP_METADATA_NAMESPACE="${PHP_METADATA_NAMESPACE:-${CAMEL_CASE}Metadata}"
    PHP_CLIENT_NAMESPACE="${PHP_CLIENT_NAMESPACE:-${CAMEL_CASE}Client}"

    DART_LIBRARY_NAME="${DART_LIBRARY_NAME:-protobuf_${LOWER_UNDER_CASE}}"
    DART_LIBRARY_DESCRIPTION="${DART_LIBRARY_DESCRIPTION:-${GENERIC_DESC}}"

    TYPESCRIPT_LIBRARY_NAME="${TYPESCRIPT_LIBRARY_NAME:-@coruja-protobuf\/${GENERIC_NAME}}"
    TYPESCRIPT_LIBRARY_DESCRIPTION="${TYPESCRIPT_LIBRARY_DESCRIPTION:-${GENERIC_DESC}}"
  fi
fi

# Set default values
AUTHOR="${AUTHOR:-Erwin Eu}"
EMAIL="${EMAIL:-eu.erwin@gmx.de}"
HOMEPAGE="${HOMEPAGE:-https:\/\/github.com\/eu-erwin}"
ROLE="${ROLE:-Developer}"
REPO_URL_SCHEMA="${REPO_URL_SCHEMA:-https:\/\/github.com\/eu-erwin\/protobuf-${GENERIC_NAME}}"
REPO_URL="${REPO_URL:-github.com\/eu-erwin\/protobuf-${GENERIC_NAME}}"

function genCmd() {
  # Generate protoc server & clients
  if [[ true == "$force" ]]; then
    echo "Generating protoc server & clients"
    # shellcheck disable=SC2091
    $($@)
    echo "Protoc server & clients generated"
  else
    # shellcheck disable=SC2145
    echo "Debug: $@"
  fi
}

function genGo() {
  read -r -d '' CMD <<- EOM
    protoc \
        -I=/go/pkg/mod/github.com/google/protobuf@v5.28.2+incompatible/src \
        -I=/go/pkg/mod/github.com/mwitkow/go-proto-validators@v0.3.2 \
        -I=. \
        --go_out=/code \
        --go-grpc_out=/code \
        --govalidators_out=/code \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        --govalidators_opt=paths=source_relative \
        --grpc-gateway_out=allow_delete_body=true:$SERVER_TARGET \
        $(find "$PROTO_PATH" -maxdepth 1 -iname "*.proto")
EOM
  echo $CMD
  genCmd "$CMD"

}

function genPhp() {
  include_files=$(printf "%s " "${files[@]}")
  read -r -d '' CMD <<- EOM
    protoc \
        -I=/go/pkg/mod/github.com/google/protobuf@v5.28.2+incompatible/src \
        --proto_path=. \
        --plugin=protoc-gen-grpc=/usr/bin/grpc_php_plugin \
        --php_out=$PHP_CLIENT_TARGET \
        --grpc_out=$PHP_CLIENT_TARGET \
        google/protobuf/empty.proto \
        google/protobuf/timestamp.proto \
        $(find "$PROTO_PATH" -maxdepth 1 -iname "*.proto")
EOM
  genCmd "$CMD"
}

function genTypescript() {
  include_files=$(printf "%s " "${files[@]}")
  read -r -d '' CMD <<- EOM
    protoc \
        -I=/go/pkg/mod/github.com/google/protobuf@v5.28.2+incompatible/src \
        --proto_path=. \
        --plugin=protoc-gen-ts=/usr/bin/protoc-gen-ts \
        --ts_out=service=grpc-web:$TYPESCRIPT_CLIENT_TARGET \
        --js_out=import_style=commonjs,binary:$TYPESCRIPT_CLIENT_TARGET \
        google/protobuf/empty.proto \
        google/protobuf/timestamp.proto \
        $(find "$PROTO_PATH" -maxdepth 1 -iname "*.proto")
EOM
  genCmd "$CMD"
}

function genDart() {
  include_files=$(printf "%s " "${files[@]}")
  read -r -d '' CMD <<- EOM
    protoc \
        -I=/go/pkg/mod/github.com/google/protobuf@v5.28.2+incompatible/src \
        --proto_path=. \
        --dart_out=grpc:$DART_CLIENT_TARGET/src \
        google/protobuf/empty.proto \
        google/protobuf/timestamp.proto \
        $(find "$PROTO_PATH" -maxdepth 1 -iname "*.proto")
EOM
  genCmd "$CMD"
}

function generate() {
    language=$(echo "$1" | tr '[:upper:]' '[:lower:]')
    echo "Generate for lang: $language"
    case $language in
        go)genGo;;
        php)genPhp;;
        dart)genDart;;
        typescript|javascript)genTypescript;;
    esac
}

function compileProtobuf() {
  if [[ $force ]]; then
    mkdir -p "$SERVER_TARGET"
    if [[ true == "$withClient" ]]; then
      mkdir -p "$CLIENT_TARGET" \
        "$TYPESCRIPT_CLIENT_TARGET" \
        "$DART_CLIENT_TARGET" \
        "$DART_CLIENT_TARGET/src" \
        "$PHP_CLIENT_TARGET"
    fi
    # Iterate the loop to read and print each array element
    for language in "${languages[@]}"
    do
      generate "$language"
    done
  else
    echo "Add --force to compile"
  fi
}

function createLibrary() {
  LANGUAGE=$1
  FILE=$2
  FORCE=$3

  if [ -e "./$FILE" ]; then
    if [[ "true" != "$FORCE" ]]; then
      echo "$FILE exists. Abort"
      return
    fi
  fi

  echo "Generate new $FILE"
  cp "/var/protobuf/template/$LANGUAGE/$FILE.temp" "./$FILE"

  sed -i "s/__AUTHOR__/$AUTHOR/g" "./$FILE"
  sed -i "s/__EMAIL__/$EMAIL/g" "./$FILE"
  sed -i "s/__HOMEPAGE__/$HOMEPAGE/g" "./$FILE"
  sed -i "s/__ROLE__/$ROLE/g" "./$FILE"
  sed -i "s/__REPO_URL__/$REPO_URL/g" "./$FILE"
}

function createFile() {
  FILE=$1
  FORCE=$2

  # Generate naming
	name=$GENERIC_NAME;
	namespace=$(helper naming -p "$name");
	package=$(helper naming -ls "$name");
	msg_filename=$(helper naming -ls "$name");
	capitalize=$(helper naming -t "$name");

  if [ -e "./$FILE" ]; then
    if [[ "true" != "$FORCE" ]]; then
      echo "$FILE exists. Abort"
      return
    fi
  fi

  echo "Generate new $FILE"
  cp "/var/protobuf/template/$FILE.temp" "./$FILE"

  sed -i "s/__NAME__/$GENERIC_NAME/g" "./$FILE"

	sed -i -e "s/__name__/$name/g" "./$FILE"
	sed -i -e "s/__package__/$package/g" "./$FILE"
	sed -i -e "s/__filename__/$msg_filename/g" "./$FILE"
	sed -i -e "s/__namespace__/$namespace/g" "./$FILE"
	sed -i -e "s/__capitalize__/$capitalize/g" "./$FILE"
	sed -i -e "s/__organization__/$organization/g" "./$FILE"
}

function createMessageService() {
  FILE=$1
  TARGET=$2

  # Generate naming
	name=$GENERIC_NAME;
	organization=$ORGANIZATION;
	namespace=$(helper naming -p "$name");
	package=$(helper naming -ls "$name");
	msg_filename=$(helper naming -ls "$name");
	capitalize=$(helper naming -t "$name");

  if [ -e "./$FILE" ]; then
    echo "$FILE exists. Abort"
    return
  fi

  echo "Generate new $FILE"
  cp "/var/protobuf/template/$FILE.temp" "./$TARGET"

  sed -i "s/__NAME__/$GENERIC_NAME/g" "./$TARGET"

	sed -i -e "s/__name__/$name/g" "./$TARGET"
	sed -i -e "s/__package__/$package/g" "./$TARGET"
	sed -i -e "s/__filename__/$msg_filename/g" "./$TARGET"
	sed -i -e "s/__namespace__/$namespace/g" "./$TARGET"
	sed -i -e "s/__capitalize__/$capitalize/g" "./$TARGET"
	sed -i -e "s/__organization__/$organization/g" "./$TARGET"
}

function createDartLibraryPackage() {
  echo "Generate Dart Library Public Api"
  libFile="./lib/$DART_LIBRARY_NAME.dart"
  createLibrary dart package.dart true
  mv ./package.dart "$libFile"
  sed -i "s/__NAME__/$DART_LIBRARY_NAME/g" "./$libFile"
  for name in $(find ./lib/src -name "*.dart"); do
    sanitizedName="${name/\/lib/}"
    echo "export '$sanitizedName';" >> "$libFile"
  done
}

function createTypeScriptPublicApi() {
  echo "Generate TypeScript Public Api"
  createLibrary typescript public-api.ts true
  for name in $(find ./client/typescript -name "*.d.ts"); do
    totalLines=$(wc -l "$name" | awk '{ print $1 }')
    if [[ $totalLines == 3 ]]; then
      continue
    fi
    sanitizedName="${name/.d.ts/}"
    echo "export * from '$sanitizedName';" >> public-api.ts
  done
}

function cleanPackages() {
    echo "Clean packages"
    cleanPackage go
    cleanPackage php
    cleanPackage dart
    cleanPackage typescript
    rm -rf "client"
    exit 0
}

function cleanPackage() {
    language=$(echo "$1" | tr '[:upper:]' '[:lower:]')
    echo "Clean package for lang: $language"
    FILES=()
    case $language in
        go)
          FILES+=(./go.mod)
          FILES+=("$SERVER_TARGET")
          FILES+=("$SERVER_TARGET/*.go")
        ;;
        php)
          FILES+=(./composer.json)
          FILES+=("$PHP_CLIENT_TARGET")
        ;;
        dart)
          FILES+=("$DART_CLIENT_TARGET")
        ;;
        typescript|javascript)
          FILES+=(./package.json)
          FILES+=(./public-api.ts)
          FILES+=(./tsconfig.json)
          FILES+=("$TYPESCRIPT_CLIENT_TARGET")
        ;;
    esac

    if [[ "${#FILES[@]}" == 0 ]]; then
      echo "No files to removed. Exit"
    fi

    for FILE in "${FILES[@]}"
    do
      echo "Removing $FILE"
      rm -rf "$FILE"
    done
}

function createPackage() {
    language=$(echo "$1" | tr '[:upper:]' '[:lower:]')
    echo "Create package for lang: $language"
    case $language in
        go)
          FILE=go.mod
          createLibrary go $FILE
          sed -i "s/__MODULE__/github.com\/eu-erwin\/protobuf-$GO_PACKAGE_NAME/g" "./$FILE"
          go get -u
          go mod tidy
        ;;
        php)
          if [[ true == "$IGNORE_CLIENT" ]]; then
            echo "Php Client ignored..."
          fi
          FILE=composer.json
          createLibrary php $FILE
          sed -i "s/__NAME__/$PHP_LIBRARY_NAME/g" "./$FILE"
          sed -i "s/__DESCRIPTION__/$PHP_LIBRARY_DESCRIPTION/g" "./$FILE"
          sed -i "s/__METADATA_NAMESPACE__/$PHP_METADATA_NAMESPACE/g" "./$FILE"
          sed -i "s/__CLIENT_NAMESPACE__/$PHP_CLIENT_NAMESPACE/g" "./$FILE"
        ;;
        dart)
          if [[ true == "$IGNORE_CLIENT" ]]; then
            echo "Dart Client ignored..."
          fi
          FILE=pubspec.yaml
          createLibrary dart $FILE
          sed -i "s/__NAME__/$DART_LIBRARY_NAME/g" "./$FILE"
          sed -i "s/__DESCRIPTION__/$DART_LIBRARY_DESCRIPTION/g" "./$FILE"
          createDartLibraryPackage
          createLibrary dart analysis_options.yaml
        ;;
        typescript|javascript)
          if [[ true == "$IGNORE_CLIENT" ]]; then
            echo "Typescript Client ignored..."
          fi
          FILE=package.json
          createLibrary typescript $FILE
          sed -i "s/__NAME__/$TYPESCRIPT_LIBRARY_NAME/g" "./$FILE"
          sed -i "s/__DESCRIPTION__/$TYPESCRIPT_LIBRARY_DESCRIPTION/g" "./$FILE"
          createTypeScriptPublicApi
          createLibrary typescript tsconfig.json
        ;;
    esac
}

SUB_COMMAND=$1
shift

PROTO_PATH="./"
SERVER_TARGET="./"
CLIENT_TARGET="./client/"

force=false
withClient=true
languages=()
files=()

while getopts l:w:p:c:s:if option
do
    case "${option}"
        in
        l)languages+=("${OPTARG}");;
        w)files+=("${OPTARG}");;
        p)path=${OPTARG};;
        c)client=${OPTARG};;
        s)server=${OPTARG};;
        i)withClient=false;;
        f)force=true;;
        *)force=false;;
    esac
done

if [[ "${#languages[@]}" == 0 ]]; then
  echo "All scripts will be generated for following languages:"
  echo "Go, Php, Typescript, Javascript, Dart"
  echo "----------------------"
  languages=("Go" "Php" "Typescript" "Dart")
else
  # shellcheck disable=SC2145
  echo "Generate for: ${languages[@]}"
  echo "----------------------"
fi

if [[ "$client" != "" ]]; then
  CLIENT_TARGET=$client
fi

if [[ "$server" != "" ]]; then
  SERVER_TARGET=$server
fi

if [[ "$path" != "" ]]; then
  PROTO_PATH=$path
fi

DART_CLIENT_TARGET="lib"
PHP_CLIENT_TARGET="${CLIENT_TARGET}php"
TYPESCRIPT_CLIENT_TARGET="${CLIENT_TARGET}typescript"

language=${languages[0]}
if [[ $SUB_COMMAND == "compile" ]]; then
  echo "Directories:"
  echo "Server: $SERVER_TARGET"
  echo "Client: $CLIENT_TARGET"
  echo "$PHP_CLIENT_TARGET"
  echo "$DART_CLIENT_TARGET"
  echo "$TYPESCRIPT_CLIENT_TARGET"
  echo "----------------------"
  compileProtobuf
elif [[ $SUB_COMMAND == "clean" ]]; then
  echo "Language: $language"
  cleanPackage "$language"
elif [[ $SUB_COMMAND == "clean-all" ]]; then
  cleanPackages
elif [[ $SUB_COMMAND == "create" ]]; then
  echo "Language: $language"
  createPackage "$language"
elif [[ $SUB_COMMAND == "init" ]]; then
  echo "Initialize"
  createFile "README.md"
  createFile ".gitlab-ci.yml"
  createMessageService "message.proto" "${GENERIC_NAME}.proto"
  createMessageService "services.proto" "${GENERIC_NAME}_service.proto"
fi

