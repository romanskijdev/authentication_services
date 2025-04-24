# Используем официальный образ Go для сборки приложения
FROM golang:1.24.0 AS builder

# Рабочая директория
WORKDIR /app

# Определяем аргумент сборки
ARG APP_NAME
RUN echo "APP_NAME=${APP_NAME}"

# Копируем код внутрь контейнера
COPY . .

# Устанавливаем Swag CLI
# RUN go install github.com/swaggo/swag/cmd/swag@latest

# Генерация документации Swagger
# RUN swag init -g cmd/app.go --dir ./,../../core

# Сборка приложения
RUN go build -o ${APP_NAME} ./cmd/app.go

# Создаём конечный образ на основе Ubuntu
FROM kartashs/ubuntu:22.04

WORKDIR /app

# Определяем аргумент сборки
ARG APP_NAME

# Копируем собранное приложение из предыдущего образа
COPY --from=builder /app/${APP_NAME} /app/${APP_NAME}

# Запуск приложения
CMD ["sh", "-c", "/app/${APP_NAME}"]