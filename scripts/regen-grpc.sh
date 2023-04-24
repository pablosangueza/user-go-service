#!/bin/bash
root=$(git rev-parse --show-toplevel)
localbin=$root/.bin

GOOUT="$root"/api/
PROTOFILES="$root"/api/dom/user/v1/

os="linux-x86_64"
if [[ "$OSTYPE" == "darwin"* ]]; then
    os="osx-x86_64"
elif [[ "$OSTYPE" == "msys" ]]; then
    os="win64"
    localbin=$(cygpath -wa $localbin)
    GOOUT=$(cygpath -wa $GOOUT)
    PROTOFILES=$(cygpath -wa $PROTOFILES)\\
fi

# check that protoc compiler exists and download it if required
PROTOBUF_VERSION=3.19.0
PROTOC_FILENAME=protoc-${PROTOBUF_VERSION}-${os}.zip
PROTOC_PATH=$localbin/protoc-$PROTOBUF_VERSION
if [ ! -d "$PROTOC_PATH" ] ; then
    mkdir -p "$PROTOC_PATH"
    curl -L "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/${PROTOC_FILENAME}" > "$localbin/$PROTOC_FILENAME"
    mkdir -p "$PROTOC_PATH"
    unzip -o "$localbin/$PROTOC_FILENAME" -d "$localbin/protoc-$PROTOBUF_VERSION"
    rm "$localbin/$PROTOC_FILENAME"
fi
# it gets the version of protoc-gen-go from go.mod file
GOBIN=$localbin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
GOBIN=$localbin go install google.golang.org/protobuf/cmd/protoc-gen-go

PATH=$PATH:$localbin "$PROTOC_PATH/bin/protoc" --go_out="$GOOUT" \
 --go_opt=paths=source_relative \
 --go-grpc_opt=paths=source_relative \
 --go-grpc_out=require_unimplemented_servers=false:$GOOUT \
 -I"$root/api" "$PROTOFILES"*.proto
