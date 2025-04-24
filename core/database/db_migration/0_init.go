package dbmigration

import (
	gormmodule "authentication_service/core/lib/external/gorm"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MigrationFunctions struct {
	pool   *pgxpool.Pool
	gormDB *gorm.DB
}

func NewMigrationFunctions(pool *pgxpool.Pool, gorm *gorm.DB) *MigrationFunctions {
	return &MigrationFunctions{pool: pool, gormDB: gorm}
}

func (m *MigrationFunctions) MigrateDB() error {
	// Миграция функций
	m.FunctionMigrations()

	m.migrateModels()

	// Закрытие соединения с базой данных
	err := gormmodule.GormDatabaseDisconnect()
	if err != nil {
		logrus.Error("❌ error: GormDatabaseDisconnect")
		return err
	}

	return nil
}
