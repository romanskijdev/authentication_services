# demo authentication_service

# Основная информация
### Нобходимое ПО
1) [**postgresql**](https://www.postgresql.org/)
2) [**Redis**](https://redis.io/)
3) [**Docker**](https://www.docker.com/)
4) [**Rabbit MQ**](https://www.rabbitmq.com/tutorials)
5) [**SMTP Service**](https://aws.amazon.com/en/what-is/smtp/) (for notifications)

### Средства разработки
1) [**golang 1.24.0**](https://go.dev/doc/devel/release)
2) [**protobuf**](https://protobuf.dev/)

### Архитектура
| Наименование                                                                              | Назначение | Тип  |
|-------------------------------------------------------------------------------------------|-------------|------|
| <a name="core_name"></a>Ядро (`core`)                                                     |  Общие схемы авторизаций, проверок, взаимодейсвия с базами данных   | lib  |
| <a name="payment_service_name"></a>Сервис уведомлений (`notification_service`)            | Сервис | -    |
| <a name="system_service_name"></a>Сервис Системных действий (`system_service`)            | Сервис | -    |
| <a name="user_service_name"></a>Сервис авторизации (`auth_service`)                       | Сервис | gRPC |
| <a name="rest_user_service_name"></a>Сервис API Пользовательских дейсвий (`user_service`) | Сервис | REST |

#### Первый запуск
1) Необходимо установить все указанные выше продукты настроить и заполнить необходимые данные
2) Настроить .yml из configcore руководствуясь информацией из config.example.yml
3) Провести билды сервисов и ядра
4) Запустить сервисы с помощью скрипта runs.sh или с помощью docker
5) все базы данные и первичные данные баз данных заполнятся/обновятся автоматически


#### Заполнить config.yml следующими параметрами


    - secure_params
```json
{
  "jwt_secret": "********",
  "salt": "********",
  "session_token_hours_life": "********",
  "admin_session_token_hours_life": "********"
}
```
    - telegram
```json
{
  "api_url": "********",
  "bot_token": "********"
}
```
    - storage
```json
{
  "0_host": "********", // хост сервера базы данных
  "0_name": "********", // имя базы данных для подключения
  "0_password": "********", // пароль подключения
  "0_port": "********", // порт базы данных
  "0_user": "********" // пользователь базы данных
}
```
    - redis
```json
{
  "addr": "********", // хост сервера redis
  "username": "********", // имя пользователя
  "password": "********", // пароль 
  "port": "********", // redis server port, по умолчанию 6379
  "db": "********" // redis database index, по умолчанию 0
}
```