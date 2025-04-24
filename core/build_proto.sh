#!/bin/bash

SERVICE_NAME="core"  # Название сервиса
PROTO_DIR="proto_base"  # Папка с .proto файлами
SERVICE_PATH=$(dirname "$(realpath "$0")")
PROTO_PATH="$SERVICE_PATH/$PROTO_DIR"


set -e
trap 'echo "🔴 Error build $SERVICE_NAME"; exit 1' ERR

echo "🟡 Start Build:  $SERVICE_NAME"
echo "🔵SERVICE_PATH:  $SERVICE_PATH"
echo "🔵PROTO_PATH:  $PROTO_PATH"


# Проверка и удаление папки SERVICE_PATH/proto, если она существует
if [ -d "$SERVICE_PATH/proto" ]; then
    echo "🔵 Removing existing proto directory: $SERVICE_PATH/proto"
    rm -rf "$SERVICE_PATH/proto"
fi

# Использование команды find для поиска всех .proto файлов в указанных директориях
proto_files=($(find "${PROTO_PATH}/messages" "${PROTO_PATH}/service" -name "*.proto"))

protoc \
    "--proto_path=${PROTO_PATH}" \
    --go_out="$SERVICE_PATH" \
    --go-grpc_out="$SERVICE_PATH" \
    "${proto_files[@]}"

cd "$SERVICE_PATH"