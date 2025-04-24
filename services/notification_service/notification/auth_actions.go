package notificationhandler

import (
	"authentication_service/core/typescore"
	"context"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

func (m *ModuleNotification) checkReqFields(notifyParams *typescore.NotifyParams) error {
	if notifyParams.Text == nil {
		log.Println("🔴 error TemporaryPasswordNotifyCategoryAction: Text is nil")
		return errors.New("text is nil")
	}
	return nil
}

// Новое устройство
func (m *ModuleNotification) DeviceNewNotifyCategoryAction(notifyParams *typescore.NotifyParams) error {
	err := m.checkReqFields(notifyParams)
	if err != nil {
		return err
	}

	t := m.ipc.TemplatesMail.NewDeviceInfoTemplate
	title := fmt.Sprintf("New Device %s", m.ipc.Config.SMTPMailServer.BaseTitle)
	bodyText := fmt.Sprintf("%s %s", "IPAddress", *notifyParams.Text)

	gMail, _ := m.CompareMailBody(t, map[string]interface{}{
		"IPAddress": *notifyParams.Text,
	}, title)

	msgList, err := m.getUsersAuthGetters(notifyParams.UsersIDs, nil, gMail, title, bodyText, notifyParams.Category)
	if err != nil {
		return err
	}
	err = m.DistributionNotify(msgList)
	return err
}

func (m *ModuleNotification) getUsersAuthGetters(systemUserIDs []*string, mailAddress *string, gMail *gomail.Message, title, bodyText string, typeNotify *typescore.NotifyCategory) ([]MsgNotifyStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	msgNotifyList := make([]MsgNotifyStruct, 0)
	if len(systemUserIDs) > 0 {
		for _, userSystemId := range systemUserIDs {
			userObj, _, errW := m.ipc.Database.Users.GetUsersListDB(ctx, typescore.ListDbOptions{Filtering: &typescore.User{
				SystemID: userSystemId,
			}})
			if errW != nil {
				log.Println("🔴 error get userObj")
				continue
			}
			if userObj == nil || len(userObj) == 0 {
				log.Println("🔴 error DeviceNewNotifyCategoryAction: userObj is nil")
				continue
			}

			msgObj := &MsgNotifyStruct{
				User:               *userObj[0],
				TitleText:          title,
				BodyText:           bodyText,
				MailMessage:        gMail,
				OnEmail:            true,
				AlertNotifyAppType: typeNotify,
			}
			msgNotifyList = append(msgNotifyList, *msgObj)
		}

	} else if mailAddress != nil {
		userInfo := typescore.User{
			Email: mailAddress,
		}
		msgObj := &MsgNotifyStruct{
			User:               userInfo,
			TitleText:          title,
			BodyText:           bodyText,
			MailMessage:        gMail,
			OnEmail:            true,
			AlertNotifyAppType: typeNotify,
		}
		msgNotifyList = append(msgNotifyList, *msgObj)
	}

	return msgNotifyList, nil
}
