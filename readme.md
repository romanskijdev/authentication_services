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
