package workerjobs

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	"github.com/sirupsen/logrus"
)

type WorkerJobs struct {
	Cfg *configcore.Config
	DB  *database.ModuleDB
}

func WorkerJobsStart(cfg *configcore.Config, database *database.ModuleDB) {
	err := database.Migrate().MigrateDB()
	if err != nil {
		logrus.Errorf("❌ error: failed to migrate database: %v", err)
		return
	}

}
