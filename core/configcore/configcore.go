package configcore

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"strings"
)

//go:embed config.yml
var configYML string

// LoadConfig загружает конфигурацию из встраиваемого YAML файла
func LoadConfig(options *ConfigLoadOptions) (*Config, error) {
	var config Config
	var node yaml.Node

	// Загружаем полный конфиг из YAML используя декодер
	readerConfig := strings.NewReader(configYML)
	decoder := yaml.NewDecoder(readerConfig)

	// Сначала декодируем в yaml.Node для получения информации о структуре
	if err := decoder.Decode(&node); err != nil {
		var yamlErr *yaml.TypeError
		if errors.As(err, &yamlErr) {
			// Если это ошибка типа, получаем подробную информацию
			var details strings.Builder
			details.WriteString("yaml decoding error: %v:\n")
			for _, e := range yamlErr.Errors {
				details.WriteString(fmt.Sprintf("- %s\n", e))
			}
			return nil, fmt.Errorf("%s", details.String())
		}
		return nil, fmt.Errorf("yaml decoding error: %v", err)
	}

	// Теперь декодируем в нашу структуру
	if err := node.Decode(&config); err != nil {
		return nil, fmt.Errorf("yaml to struct conversion error: %v", err)
	}

	// Создаем новый конфиг только с нужными секциями
	result := &Config{}

	// Копируем только нужные поля на основе опций
	copyConfigFields(options, &config, result)

	logrus.Info("✅ Config loaded successfully")
	return result, nil
}

// copyConfigFields копирует поля из source в target на основе опций options
func copyConfigFields(options *ConfigLoadOptions, source, target *Config) {
	// Копируем Database
	if options.Database {
		target.Database = source.Database
	}

	// Копируем RabbitMQ
	if options.RabbitMQConfig {
		target.RabbitMQConfig = source.RabbitMQConfig
	}

	// Копируем Secrets
	if options.Secrets.Admin {
		target.Secrets.AuthJWT.AdminSecret = source.Secrets.AuthJWT.AdminSecret
	}
	if options.Secrets.User {
		target.Secrets.AuthJWT.UserSecret = source.Secrets.AuthJWT.UserSecret
	}

	// Копируем GrpsClients
	if options.GrpsClients.AuthService {
		target.GrpsClients.AuthService = source.GrpsClients.AuthService
	}

	// Копируем ExposedServiceConfig
	if options.ExposedServiceConfig.UserService {
		target.ExposedServiceConfig.UserService = source.ExposedServiceConfig.UserService
	}
	if options.ExposedServiceConfig.AuthService {
		target.ExposedServiceConfig.AuthService = source.ExposedServiceConfig.AuthService
	}
	if options.ExposedServiceConfig.NotifyService {
		target.ExposedServiceConfig.NotifyService = source.ExposedServiceConfig.NotifyService
	}
}
