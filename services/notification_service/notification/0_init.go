package notificationhandler

import (
	"authentication_service/core/typescore"
	typesm "authentication_service/notification_service/types"
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type ModuleNotification struct {
	ipc *typesm.InternalProviderControl
}

func NewModuleNotification(ipc *typesm.InternalProviderControl) *ModuleNotification {
	return &ModuleNotification{
		ipc: ipc,
	}
}

const maxConcurrentProcessors = 100

// InitRabbitMQConsumer инициализирует потребителя RabbitMQ
func (m *ModuleNotification) InitRabbitMQConsumer(queueName, consumerTag string) {
	consumer := m.ipc.RabbitMQClient.NewConsumer()

	// Семафор для ограничения количества одновременных обработчиков
	sem := make(chan struct{}, maxConcurrentProcessors)

	err := consumer.ConsumeSimple(queueName, consumerTag, func(body []byte) error {
		// Захватываем слот в семафоре
		sem <- struct{}{}

		// Обрабатываем каждое сообщение в отдельной горутине
		go func(messageBody []byte) {
			defer func() {
				<-sem
				// В случае паники логируем ошибку
				if r := recover(); r != nil {
					logrus.Errorln("🔴 Handler: error in handler: ", r)
				}
			}()

			var msg interface{}
			if err := json.Unmarshal(messageBody, &msg); err != nil {
				logrus.Errorln("🔴 Handler: error decoding message: ", err)
				return
			}

			formattedJSONTg, err := json.MarshalIndent(msg, "", "  ")
			if err != nil {
				logrus.Errorln("🛑 error formatting JSON: ", err)
				return
			}

			update := &typescore.NotifyParams{}
			if err := json.Unmarshal(formattedJSONTg, update); err != nil {
				logrus.Error("🛑 error unmarshaling update: ", err)
				return
			}

			// обрабатываем сообщение
			errW := m.processNotification(update)
			if err != nil {
				logrus.Errorln("🔴 Handler: error processing message: ", errW)
				return
			}
		}(body)

		// Сразу подтверждаем получение, так как обработка идёт асинхронно
		return nil
	})

	if err != nil {
		logrus.Errorln("🔴 Handler: error starting consumer: ", err)
		return
	}
}

func (m *ModuleNotification) processNotification(notificationObj *typescore.NotifyParams) error {
	err := m.NotifyRouting(notificationObj)
	if err != nil {
		logrus.Errorln("🔴 error NotifyUser NotifyRouting: ", err)
		return err
	}
	return nil
}
