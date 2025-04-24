package typesm

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
)

type InternalProviderControl struct {
	Config   *configcore.Config
	Database *database.ModuleDB
}
