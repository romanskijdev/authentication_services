#!/bin/bash

set -e
trap 'echo "🔴 Скрипт остановлен"; exit' SIGINT

# Имя сети Docker
DOCKER_NETWORK_NAME="networklocal1"
# Наименование сервиса
SERVICE_NAME="demo_notifications_service"

# Наименование образа
IMAGE_NAME="${SERVICE_NAME}_i"
# Наименование контейнера
CONTAINER_NAME="${SERVICE_NAME}_c"
# Имя приложения
APP_NAME="${SERVICE_NAME}_app"

# Останавливаем и удаляем старый контейнер, если он существует
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
fi

# Удаляем старый образ, если он существует
if [ "$(docker images -q $IMAGE_NAME)" ]; then
    docker rmi $IMAGE_NAME
fi

# Создаем директорию vendor с зависимостями
go mod tidy
rm -rf vendor
go mod vendor

# Проверяем существование сети и создаем, если не существует
docker network ls | grep $DOCKER_NETWORK_NAME || docker network create $DOCKER_NETWORK_NAME

# Собираем образ
docker build --build-arg APP_NAME=$APP_NAME -t $IMAGE_NAME .

# Запускаем контейнер с новым образом в сетевом режиме host и подключаем к сети по имени
docker run -d \
    --network $DOCKER_NETWORK_NAME \
    --restart always \
    --name $CONTAINER_NAME \
    --add-host=database_host:192.168.0.15 \
    -e APP_NAME=$APP_NAME \
    $IMAGE_NAME