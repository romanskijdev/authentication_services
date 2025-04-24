package configcore

// SecretsOptions определяет, какие секреты нужно загружать
type SecretsOptions struct {
	Admin bool
	User  bool
}

// GrpsClientsOptions определяет, какие gRPC клиенты нужно загружать
type GrpsClientsOptions struct {
	AuthService bool
}

// ExposedServiceOptions определяет, какие сервисы нужно загружать
type ExposedServiceOptions struct {
	UserService   bool
	AuthService   bool
	NotifyService bool
}

// ConfigLoadOptions определяет, какие группы конфигурации нужно загружать
type ConfigLoadOptions struct {
	Database             bool
	RabbitMQConfig       bool
	Telegram             bool
	Secrets              SecretsOptions
	GrpsClients          GrpsClientsOptions
	ExposedServiceConfig ExposedServiceOptions
}
