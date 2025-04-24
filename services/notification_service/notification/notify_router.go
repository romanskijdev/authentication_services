package notificationhandler

import (
	"authentication_service/core/typescore"
	"errors"

	"gopkg.in/gomail.v2"
)

type MsgNotifyStruct struct {
	User      typescore.User
	TitleText string
	BodyText  string

	MailMessage        *gomail.Message
	OnEmail            bool
	AlertNotifyAppType *typescore.NotifyCategory
}

func (m *ModuleNotification) NotifyRouting(notifyParams *typescore.NotifyParams) error {
	if notifyParams == nil {
		return errors.New("notifyParams is nil")
	}

	if notifyParams.Category == nil {
		return errors.New("notifyParams.Category is nil")
	}

	switch *notifyParams.Category {
	case typescore.DeviceNewNotifyCategory: // Новое устройство
		return m.DeviceNewNotifyCategoryAction(notifyParams)
	}
	return nil
}
