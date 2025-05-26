#!/usr/bin/env bash

set -euo pipefail
shopt -s globstar

if ! [[ "$0" =~ scripts/generate_protobuf.sh ]]; then
  echo "The script must be executed from repository root."
  exit 255
fi

PROTOBUF_ROOT="./api/protobuf"

function generate() {
  local output_dir=$1
  local service_name=$2

  mkdir -p "$output_dir/common/protobuf/${service_name}pb"
  find "$output_dir/common/protobuf/${service_name}pb" -type f -name "*.pb.go" -delete

  protoc -I="$(brew --prefix)/include/google/protobuf" -I="$PROTOBUF_ROOT" \
    --go_out="$output_dir" \
    --go-grpc_opt=require_unimplemented_servers=false \
    --go-grpc_out="$output_dir" \
    "$PROTOBUF_ROOT/$service_name/$service_name.proto"
}

# protoc 会在当前目录 . 按照 proto 文件中的 go_package 路径生成代码
generate . stock
generate . order

echo "protobuf code generate success!"
