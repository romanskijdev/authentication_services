package tablesmigration

import (
	dbcoretablenames "authentication_service/core/database/table_names"
	"authentication_service/core/typescore"
	"gorm.io/gorm"
	"time"
)

type LocalUserProvider typescore.User

func (LocalUserProvider) TableName() string {
	return dbcoretablenames.TableNameUsers.ToString()
}

func UserTableMigrate(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&LocalUserProvider{})

	// Выполняем автосоздание таблицы
	err := db.AutoMigrate(&LocalUserProvider{})
	if err != nil {
		return err
	}

	if !hasTable {
		time.Sleep(5 * time.Second) // Добавляем задержку

		// Add comments to table and columns
		db.Exec(`
            COMMENT ON TABLE users IS 'Таблица для хранения данных о пользователях';
        `)
	}
	return nil
}
