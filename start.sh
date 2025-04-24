#!/bin/bash

set -e
trap 'echo "🔴 Скрипт остановлен"; exit' SIGINT

./stop_all.sh

# Определяем абсолютные пути к директориям
SCRIPT_DIR=$(realpath "$(dirname "$0")")


SERVICES_DIR="$SCRIPT_DIR/services"
CORE_DIR="$SCRIPT_DIR/core"
echo "🟡 Директория скрипта: $SCRIPT_DIR"

# Переходим в директорию core и запускаем go vet ./... и go fmt ./...
if [ -d "$CORE_DIR" ]; then
    echo "🟡 Переход в директорию $CORE_DIR и запуск go vet ./... и go fmt ./..."
    cd "$CORE_DIR"
    if [ -f "go.mod" ] && [ -f "go.sum" ]; then
        go vet ./...
        go fmt ./...
        echo "🟢 Команды go vet ./... и go fmt ./... выполнены"
    else
        echo "🔴 Файлы go.mod или go.sum не найдены в $CORE_DIR, пропуск"
    fi
    cd - > /dev/null
else
    echo "🔴 Директория $CORE_DIR не найдена, пропуск"
fi

# Проходим по всем поддиректориям в директории lib
for service_dir in "$SERVICES_DIR"/*/; do
    # Проверяем, существуют ли файлы go.mod и go.sum в текущей поддиректории
    if [ -f "${service_dir}go.mod" ] && [ -f "${service_dir}go.sum" ] && [ -f "${service_dir}start.sh" ]; then
        echo "🟡 Переход в директорию ${service_dir}"
        cd "${service_dir}"
        echo "Текущая директория: $(pwd)"
        go mod tidy
        go mod vendor
        go vet ./...
        go fmt ./...
        echo "🟢 Команды go vet ./... и go fmt ./... выполнены в ${service_dir}"
        # Проверяем, существует ли скрипт start.sh в текущей поддиректории
        if [ -f "${service_dir}start.sh" ]; then
            echo "🟡 Запуск скрипта ${service_dir}start.sh"
            ls 
            # Делаем скрипт исполняемым, если это еще не сделано
            chmod +x "${service_dir}start.sh"
            # Запускаем скрипт start.sh
            "${service_dir}start.sh"
            echo "🟢 Скрипт ${service_dir}start.sh завершен"
        else
            echo "🔴 Скрипт ${service_dir}start.sh не найден, пропуск"
        fi
    else
        echo "🔴 Файлы go.mod или go.sum  или start.sh не найдены в ${service_dir}, пропуск"
    fi
done

echo "🟢 Все скрипты start.sh выполнены"