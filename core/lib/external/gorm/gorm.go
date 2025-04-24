package gormmodule

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ConfigConnectGorm struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
}

var (
	gormDB *gorm.DB
	once   sync.Once
)

func connectGormDB(configObj *ConfigConnectGorm) (*gorm.DB, error) {
	if configObj == nil {
		return nil, errors.New("connectGormDB: configObj is nil")
	}

	log.Println("🚀 ConnectGormDB")
	sslMode := "disable"
	if configObj.SslMode != "" {
		sslMode = configObj.SslMode
	}

	// Поддерживаем только disable и require
	if sslMode != "disable" && sslMode != "require" {
		return nil, errors.New("ConnectDB: unsupported sslmode, only 'disable' and 'require' are allowed")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configObj.Host,
		configObj.Port,
		configObj.User,
		configObj.Password,
		configObj.Name,
		sslMode)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // медленный SQL порог
			LogLevel:                  logger.Silent, // уровень логирования
			IgnoreRecordNotFoundError: true,          // игнорировать ошибки записи не найдены
			Colorful:                  false,         // отключить цвет
		},
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		logrus.Errorf("🔴 error: %s: %+v", "connectGormDB-Open", err)
		return nil, err
	}

	return db, nil
}

// Подключение к базе данных GORM с использованием Singleton
func GormDatabaseConnect(configObj *ConfigConnectGorm) (*gorm.DB, *sql.DB, error) {
	logrus.Info("🚀 GormDatabaseConnect")
	if configObj == nil {
		return nil, nil, errors.New("configObj is nil")
	}

	var err error
	once.Do(func() {
		gormDB, err = connectGormDB(configObj)
	})

	if err != nil {
		return nil, nil, err
	}

	if gormDB == nil {
		return nil, nil, errors.New("failed to connect database: GormDB instance is nil")
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, errors.New("failed to get DB from GORM")
	}

	// SetMaxIdleConns устанавливает максимальное количество соединений в пуле простаивающих соединений.
	sqlDB.SetMaxIdleConns(2)

	// SetMaxOpenConns устанавливает максимальное количество открытых соединений с базой данных.
	sqlDB.SetMaxOpenConns(10)

	// SetConnMaxLifetime устанавливает максимальное время, в течение которого может быть повторно использовано соединение.
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	gormDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	return gormDB, sqlDB, nil
}

// Отключение от базы данных GORM
func GormDatabaseDisconnect() error {
	if gormDB == nil {
		return nil
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	// Закрываем все открытые соединения
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connections: %w", err)
	}

	// Очищаем пул соединений
	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetMaxOpenConns(0)

	// Обнуляем глобальную переменную
	gormDB = nil
	return nil
}
