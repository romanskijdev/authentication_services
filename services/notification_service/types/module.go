package typesm

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	rabbitmqlib "authentication_service/core/lib/external/rabbitmq-lib"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"html/template"
)

type TemplatesMailSystem struct {
	NewDeviceInfoTemplate *template.Template
}

type InternalProviderControl struct {
	TemplatesMail  *TemplatesMailSystem
	Config         *configcore.Config
	RabbitMQClient *rabbitmqlib.ConnectionRabitMq
	Database       *database.ModuleDB
	BundleI18n     *i18n.Bundle
}
