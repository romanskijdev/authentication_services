package database

import (
	dbcore "authentication_service/core/database/db"
	dbmigration "authentication_service/core/database/db_migration"
	gormmodule "authentication_service/core/lib/external/gorm"
	pgxpoolmodule "authentication_service/core/lib/external/pgxpool"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type ModuleDB struct {
	Pool          *pgxpool.Pool
	gormDB        *gorm.DB
	Users         dbcore.UserDBI
	Notifications *dbcore.NotificationDB
}

func NewModuleDB(
	config *pgxpoolmodule.ConfigConnectPgxPool,
) (*ModuleDB, error) {
	pool, err := pgxpoolmodule.ConnectDB(config)
	if err != nil {
		return nil, err
	}

	gormDataBase, _, err := gormmodule.GormDatabaseConnect(&gormmodule.ConfigConnectGorm{
		Host:     config.Host,
		Port:     config.Port,
		User:     config.User,
		Password: config.Password,
		Name:     config.Name,
		SslMode:  config.SSLMode,
	})
	if err != nil {
		return nil, err
	}

	moduleDB := &ModuleDB{Pool: pool, gormDB: gormDataBase}
	moduleDB = initDBModules(moduleDB)

	return moduleDB, nil
}

func initDBModules(modules *ModuleDB) *ModuleDB {
	modules.Users = dbcore.NewUserDB(modules.Pool)
	modules.Notifications = dbcore.NewNotificationDB(modules.Pool)
	return modules
}

func (m *ModuleDB) Migrate() *dbmigration.MigrationFunctions {
	return dbmigration.NewMigrationFunctions(m.Pool, m.gormDB)
}
