package main

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	pgxpoolmodule "authentication_service/core/lib/external/pgxpool"
	rabbitmqlib "authentication_service/core/lib/external/rabbitmq-lib"
	"authentication_service/core/variables"
	"authentication_service/notification_service/loader"
	"authentication_service/notification_service/locale"
	notificationhandler "authentication_service/notification_service/notification"
	typesm "authentication_service/notification_service/types"
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

	notifications := notificationhandler.NewModuleNotification(ipc)
	if err := ipc.RabbitMQClient.DeclareQueue(
		variables.RabbitMQNotificationsQueueName,
		variables.RabbitMQExchangeNotifications,
		variables.RabbitMQNotificationsServiceRoute); err != nil {
		logrus.Errorln("🔴 DeclareQueue: failed to declare queue: ", err)
		select {}
	}
	go notifications.InitRabbitMQConsumer(variables.RabbitMQNotificationsQueueName, variables.RabbitMQNotificationsServiceConsumer)

	// Блокировка основной горутины
	select {}
}

func getConfig() (*configcore.Config, error) {
	options := &configcore.ConfigLoadOptions{
		Database:       true,
		RabbitMQConfig: true,
		Telegram:       true,
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

	return &typesm.InternalProviderControl{
		Config:         configObj,
		RabbitMQClient: rabbitMQClient,
		Database:       db,
		TemplatesMail:  loader.LoadMailTemplates(),
		BundleI18n:     locale.I8nInit(),
	}, nil
}
