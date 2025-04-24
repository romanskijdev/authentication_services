package dbmigration

import (
	tablesmigration "authentication_service/core/database/db_migration/tables"
	"github.com/sirupsen/logrus"
	"time"
)

func (m *MigrationFunctions) migrateModels() {
	logrus.Info("🚀 Start db table migration ", time.Now().UTC())

	db := m.gormDB

	// Миграция таблиц пользователей
	err := tablesmigration.UserTableMigrate(db)
	if err != nil {
		logrus.Errorf("failed to migrate user table: %v", err)
		return
	}
}
