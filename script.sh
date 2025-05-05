#!/bin/bash

force=false
ignoreClient=false
languages=()
files=()
PROTO_PATH="./"
SERVER_TARGET="./"
CLIENT_TARGET="./client/"

while getopts l:w:p:c:s:if option
do
    case "${option}"
        in
        l)languages+=("${OPTARG}");;
        w)files+=("${OPTARG}");;
        p)path=${OPTARG};;
        c)client=${OPTARG};;
        s)server=${OPTARG};;
        i)ignoreClient=true;;
        f)force=true;;
        *)force=false;;
    esac
done

if [[ "${#languages[@]}" == 0 ]]; then
  if [[ $ignoreClient == true ]]; then
    echo "Generate golang code only"
    echo "Ignores Php, Typescript, Javascript, Dart clients"
    echo "----------------------"
    languages=("Go")
  else
    echo "All scripts will be generated for following languages:"
    echo "Go, Php, Typescript, Javascript, Dart"
    echo "----------------------"
    languages=("Go" "Php" "Typescript" "Dart")
  fi
else
  # shellcheck disable=SC2145
  echo "Generate for: ${languages[@]}"
  echo "----------------------"
fi

ignoreClientArg=""
if [[ $ignoreClient == true ]]; then
  ignoreClientArg="-i"
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

include_files=""

# Loop through each element and concatenate them with a delimiter
for file in "${files[@]}"; do
    include_files+="-w $file "
done

for language in "${languages[@]}"
do
  echo "$language"
  if [[ "$force" == true ]]; then
    compiler compile -l "$language" -p "$PROTO_PATH" $include_files -s "$SERVER_TARGET" -c "$CLIENT_TARGET" -f "$ignoreClientArg"
    compiler create -l "$language" -f "$ignoreClientArg"
  else
    echo "compiler compile -l $language -p $PROTO_PATH $include_files -s $SERVER_TARGET -c $CLIENT_TARGET $ignoreClientArg"
    compiler compile -l "$language" -p "$PROTO_PATH" $include_files -s "$SERVER_TARGET" -c "$CLIENT_TARGET" "$ignoreClientArg"
    compiler create -l "$language" "$ignoreClientArg"
  fi
done