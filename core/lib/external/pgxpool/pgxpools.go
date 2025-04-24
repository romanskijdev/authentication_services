package pgxpoolmodule

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigConnectPgxPool struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

var (
	pool *pgxpool.Pool
	once sync.Once
)

func ConnectDB(configObj *ConfigConnectPgxPool) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		if configObj == nil {
			err = errors.New("ConnectDB: configObj is nil")
			return
		}

		log.Println("🚀 ConnectDB")
		sslMode := "disable"
		if configObj.SSLMode != "" {
			sslMode = configObj.SSLMode
		}

		// Поддерживаем только disable и require
		if sslMode != "disable" && sslMode != "require" {
			err = errors.New("ConnectDB: unsupported sslmode, only 'disable' and 'require' are allowed")
			return
		}

		log.Println("🚀 sslMode: ", sslMode)
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s target_session_attrs=read-write",
			configObj.Host,
			configObj.Port,
			configObj.User,
			configObj.Password,
			configObj.Name,
			sslMode)

		configDB, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			err = errors.New("ConnectDB: failed to parse config: " + err.Error())
			return
		}

		// настройки пула
		configDB.MaxConns = 2                         // Максимальное количество открытых соединений в пуле
		configDB.MinConns = 1                         // Минимальное количество открытых соединений, которое пул будет поддерживать
		configDB.MaxConnLifetime = time.Second * 45   // Максимальное время жизни соединения. После этого времени соединение будет закрыто
		configDB.MaxConnIdleTime = time.Second * 45   // Максимальное время простоя соединения. Если соединение не используется в течение этого времени, оно будет закрыто
		configDB.HealthCheckPeriod = time.Second * 15 // Периодичность проверки состояния соединений в пуле

		pool, err = pgxpool.NewWithConfig(context.Background(), configDB)
		if err != nil {
			err = errors.New("ConnectDB: failed to connect to database: " + err.Error())
			return
		}
	})

	if err != nil {
		return nil, err
	}

	if pool == nil {
		return nil, errors.New("ConnectDB: failed to connect to database: pool is nil")
	}

	return pool, nil
}
