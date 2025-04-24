package configcore

import (
	"aidanwoods.dev/go-paseto"
)

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DB       string `yaml:"db" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	SSL      string `yaml:"ssl" env-required:"true" default:"disable"`
}

// RabbitMQConfig конфигурация RabbitMQ
type RabbitMQConfig struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
}

// BucketStorageConfig конфигурация хранилища
type BucketStorageConfig struct {
	Key      string `yaml:"key" env-required:"true"`
	Secret   string `yaml:"secret" env-required:"true"`
	Region   string `yaml:"region" env-required:"true"`
	Bucket   string `yaml:"bucket" env-required:"true"`
	Endpoint string `yaml:"endpoint" env-required:"true"`
}

// TelegramConfig конфигурация Telegram
type TelegramConfig struct {
	BotToken  string `yaml:"bot_token"`
	WebAppUrl string `yaml:"web_app_url"`
}

// AuthJWTConfig конфигурация JWT аутентификации
type AuthJWTConfig struct {
	AdminSecret string `yaml:"admin_secret"`
	UserSecret  string `yaml:"user_secret"`
}

// SecretsConfig конфигурация секретов
type SecretsConfig struct {
	AESBucketKey string        `yaml:"aes_bucket_key"`
	AuthJWT      AuthJWTConfig `yaml:"auth_jwt"`
}

// ServiceConfig конфигурация сервиса
type ServiceConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
}

// GrpsClientsConfig конфигурация gRPC клиентов
type GrpsClientsConfig struct {
	AuthService ServiceConfig `yaml:"auth_service"`
}

// CorsConfig конфигурация CORS
type CorsConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// SwaggerConfig конфигурация Swagger
type SwaggerConfig struct {
	User string `yaml:"user" env-required:"true"`
	Pass string `yaml:"pass" env-required:"true"`
}

// RestServiceConfig конфигурация REST сервиса
type RestServiceConfig struct {
	PortRest int           `yaml:"port_rest" env-required:"true"`
	Cors     CorsConfig    `yaml:"cors"`
	Swagger  SwaggerConfig `yaml:"swagger"`
}

// PASETOConfig конфигурация PASETO
type PASETOConfig struct {
	SymmetricKey  paseto.V4SymmetricKey
	ImplicitBytes []byte
}

// GrpcServiceConfig конфигурация gRPC сервиса
type GrpcServiceConfig struct {
	GrpcPort     int `yaml:"grpc_port" env-required:"true"`
	CertFileName string
	PASETO       PASETOConfig
}

// ExposedServiceConfig конфигурация всех сервисов
type ExposedServiceConfig struct {
	UserService   RestServiceConfig `yaml:"user_service"`
	AuthService   GrpcServiceConfig `yaml:"auth_service"`
	NotifyService GrpcServiceConfig `yaml:"notify_service"`
}

type SMTPMailServer struct {
	BaseMail     string `yaml:"base_mail" env-required:"true"`
	BaseTitle    string `yaml:"base_title" env-required:"true"`
	SMTPPassword string `yaml:"smtp_password" env-required:"true"`
	SMTPHost     string `yaml:"smtp_host" env-required:"true"`
	SMTPPort     string `yaml:"smtp_port" env-required:"true"`
}

// Config основной конфиг
type Config struct {
	SMTPMailServer       SMTPMailServer       `yaml:"smtp_mail_server"`
	Database             DatabaseConfig       `yaml:"database"`
	RabbitMQConfig       RabbitMQConfig       `yaml:"rabbitmq"`
	Secrets              SecretsConfig        `yaml:"secrets"`
	GrpsClients          GrpsClientsConfig    `yaml:"grps_clients"`
	ExposedServiceConfig ExposedServiceConfig `yaml:"exposed_service_config"`
}
