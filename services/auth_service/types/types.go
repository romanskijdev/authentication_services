package typesm

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	rabbitmqlib "authentication_service/core/lib/external/rabbitmq-lib"
)

type InternalProviderControl struct {
	Config   *configcore.Config
	RabbitMQ *rabbitmqlib.ConnectionRabitMq
	Database *database.ModuleDB
}
