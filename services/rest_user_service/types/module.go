package typesm

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	rabbitmqlib "authentication_service/core/lib/external/rabbitmq-lib"
	protoobj "authentication_service/core/proto"
)

type InternalProviderControl struct {
	RabbitMQ               *rabbitmqlib.ConnectionRabitMq
	Config                 *configcore.Config
	DB                     *database.ModuleDB
	ClientAuthServiceProto protoobj.AuthServiceClient
}
