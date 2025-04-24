#!/bin/bash

SERVICE_NAMES=(
    "tmail_telegram_service"
    "tmail_payment_service"
    "tmail_w_smtp_service"
    "tmail_w_imap_service"
    "tmail_user_service"
    "tmail_system_service"
    "tmail_tariff_service"
    "tmail_admin_service"
    "tmail_notification_service"
)

# Останавливаем и удаляем контейнеры
for CONTAINER_NAME in "${SERVICE_NAMES[@]}"; do
    FULL_CONTAINER_NAME="${CONTAINER_NAME}_c"
    if [ "$(docker ps -q -f name=$FULL_CONTAINER_NAME)" ]; then
        echo "🟡 Остановка контейнера $FULL_CONTAINER_NAME"
        docker stop $FULL_CONTAINER_NAME
        echo "🟡 Удаление контейнера $FULL_CONTAINER_NAME"
        docker rm $FULL_CONTAINER_NAME
    else
        echo "🔴 Контейнер $FULL_CONTAINER_NAME не найден"
    fi
done

# Удаляем образы
for IMAGE_NAME in "${SERVICE_NAMES[@]}"; do
    FULL_IMAGE_NAME="${IMAGE_NAME}_i"
    if [ "$(docker images -q $FULL_IMAGE_NAME)" ]; then
        echo "🟡 Удаление образа $FULL_IMAGE_NAME"
        docker rmi $FULL_IMAGE_NAME
    else
        echo "🔴 Образ $FULL_IMAGE_NAME не найден"
    fi
done

echo "🟢 Все контейнеры и образы удалены"