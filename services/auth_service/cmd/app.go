package main

import (
	grpcauth "authentication_service/auth_service/common/grpc"
	typesm "authentication_service/auth_service/types"
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	pgxpoolmodule "authentication_service/core/lib/external/pgxpool"
	rabbitmqlib "authentication_service/core/lib/external/rabbitmq-lib"
	"authentication_service/core/variables"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		logrus.Errorln("❌ Failed to load config: ", err)
		select {}
	}

	ipc, err := initInternalProvider(cfg)
	if err != nil {
		logrus.Errorln("❌ Failed to init internal provider: ", err)
		select {}
	}

	if err := ipc.RabbitMQ.DeclareQueue(variables.RabbitMQAuthQueueName, variables.RabbitMQExchangeAuth, variables.RabbitMQNotificationsServiceRoute); err != nil {
		logrus.Errorln("🔴 DeclareQueue: failed to declare queue: ", err)
		select {}
	}

	// Запуск gRPC сервера
	go grpcauth.StartGrpcServer(ipc)

	// Блокировка основной горутины
	select {}
}

func getConfig() (*configcore.Config, error) {
	options := &configcore.ConfigLoadOptions{
		Database:       true,
		RabbitMQConfig: true,
		ExposedServiceConfig: configcore.ExposedServiceOptions{
			AuthService: true,
		},
	}

	cfg, err := configcore.LoadConfig(options)
	if err != nil {
		logrus.Errorln("❌ Failed to load config: ", err)
		return nil, err
	}
	return cfg, nil
}

func initInternalProvider(configObj *configcore.Config) (*typesm.InternalProviderControl, error) {
	rabbitMQClient, err := rabbitmqlib.ConnectRabitMq(rabbitmqlib.ConnectParams{
		Username: configObj.RabbitMQConfig.User,
		Host:     configObj.RabbitMQConfig.Host,
		Port:     configObj.RabbitMQConfig.Port,
		Password: configObj.RabbitMQConfig.Password,
	})
	if err != nil || rabbitMQClient == nil {
		logrus.Error("❌ Failed to init RabbitMQ client: ", err)
		return nil, err
	}

	db, err := database.NewModuleDB(&pgxpoolmodule.ConfigConnectPgxPool{
		Host:     configObj.Database.Host,
		Port:     fmt.Sprintf("%d", configObj.Database.Port),
		User:     configObj.Database.User,
		Password: configObj.Database.Password,
		Name:     configObj.Database.DB,
		SSLMode:  configObj.Database.SSL,
	})
	if err != nil {
		logrus.Errorln("❌ Failed to init database: ", err)
		return nil, err
	}

	if db == nil {
		logrus.Error("❌ failed to init db pool: nil")
		return nil, errors.New("❌ failed to init db pool: nil")
	}

	return &typesm.InternalProviderControl{
		Config:   configObj,
		RabbitMQ: rabbitMQClient,
		Database: db,
	}, nil
}
