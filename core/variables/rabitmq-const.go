package variables

// Переменные для работы с RabbitMQ
var (
	RabbitMQExchangeAuth        = "demo-exchange_auth"         // название обменника для Auth-сообщений
	RabbitMQAuthQueueName       = "demo-queue_auth"            // название очереди для платежей в Auth
	RabbitMQAuthServiceRoute    = "demo-auth_service"          // наименование платежного сервиса
	RabbitMQAuthServiceConsumer = "demo-auth_service_consumer" // наименование потребителя платежного сервиса

	RabbitMQExchangeNotifications        = "demo-exchange_notifications"         // название обменника для уведомлений
	RabbitMQNotificationsQueueName       = "demo-queue_notifications"            // название очереди уведомлений
	RabbitMQNotificationsServiceRoute    = "demo-notifications_service"          // наименование сервиса уведомлений
	RabbitMQNotificationsServiceConsumer = "demo-notifications_service_consumer" // наименование потребителя сервиса уведомлений

)
