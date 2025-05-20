#!/usr/bin/env bash

set -euo pipefail
shopt -s globstar

if ! [[ "$0" =~ scripts/generate_openapi.sh ]]; then
  echo "The script must be executed from repository root."
  exit 255
fi

OPENAPI_ROOT="./api/openapi"
SERVER_TYPE="gin-server"

function generate() {
  local output_dir=$1
  local package_name=$2
  local service_name=$3

  mkdir -p "$output_dir"
  find "$output_dir" -type f -name "*.gen.go" -delete

  # code for server
  oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package_name" "$OPENAPI_ROOT/$service_name/openapi.yaml"
  oapi-codegen -generate "$SERVER_TYPE" -o "$output_dir/openapi_server.gen.go" -package "$package_name" "$OPENAPI_ROOT/$service_name/openapi.yaml"
}

generate order/ports ports order

echo "openapi code generate success!"
