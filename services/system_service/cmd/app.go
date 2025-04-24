package main

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	pgxpoolmodule "authentication_service/core/lib/external/pgxpool"
	"fmt"
	workerjobs "sveves-team/tmail-mail-backend/system-service/jobs"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		logrus.Errorln("❌ Failed to load config: ", err)
		select {}
	}

	db, err := database.NewModuleDB(&pgxpoolmodule.ConfigConnectPgxPool{
		Host:     cfg.Database.Host,
		Port:     fmt.Sprintf("%d", cfg.Database.Port),
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.DB,
		SSLMode:  cfg.Database.SSL,
	})
	if err != nil {
		logrus.Errorln("❌ Failed to init database: ", err)
		select {}
	}

	// Запуск воркеров
	workerjobs.WorkerJobsStart(cfg, db)

	// Блокируем основной поток
	select {}
}

func getConfig() (*configcore.Config, error) {
	options := &configcore.ConfigLoadOptions{
		Database: true,
	}
	cfg, err := configcore.LoadConfig(options)
	if err != nil {
		logrus.Errorln("❌ Failed to load config: ", err)
		select {}
	}
	return cfg, nil
}
