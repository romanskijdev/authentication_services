package notificationhandler

import (
	"authentication_service/core/utilscore"
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"strconv"
	"sync"
)

func (m *ModuleNotification) DistributionNotify(msgObj []MsgNotifyStruct) error {

	var wg sync.WaitGroup

	for _, msg := range msgObj {
		wg.Add(1)
		go func(msg MsgNotifyStruct) {
			defer wg.Done()
			m.platformSending(msg)
		}(msg)

	}
	wg.Wait()

	return nil
}

func (m *ModuleNotification) platformSending(msg MsgNotifyStruct) {

	if msg.OnEmail && msg.User.Email != nil && msg.MailMessage != nil {
		go m.emailSendNotification(msg)
	}
}

func (m *ModuleNotification) emailSendNotification(msg MsgNotifyStruct) {
	emailRecipient := *msg.User.Email
	gMessage := msg.MailMessage
	errW := utilscore.ValidateEmailFormat(msg.User.Email)
	if errW != nil {
		log.Println("🔴 error Email send notification: ", errW)
		return
	}
	gMessage.SetHeader("From", m.ipc.Config.SMTPMailServer.BaseMail)
	gMessage.SetHeader("To", emailRecipient)

	port, err := strconv.Atoi(m.ipc.Config.SMTPMailServer.SMTPPort)
	if err != nil {
		log.Println("🔴 error Email send notification: ", err)
		return
	}
	dialer := gomail.NewDialer(m.ipc.Config.SMTPMailServer.SMTPHost,
		port,
		m.ipc.Config.SMTPMailServer.BaseMail,
		m.ipc.Config.SMTPMailServer.SMTPPassword)

	err = dialer.DialAndSend(gMessage)
	if err != nil {
		log.Println("🔴 error Email send notification: ", err)
		return
	}
}

func (m *ModuleNotification) CompareMailBody(t *template.Template, msg map[string]interface{}, title string) (*gomail.Message, error) {
	var body bytes.Buffer
	err := t.Execute(&body, msg)
	if err != nil {
		return nil, err
	}
	message := gomail.NewMessage()
	message.SetHeader("Subject", title)
	message.SetBody("text/html", body.String())

	return message, nil
}

func (m *ModuleNotification) GenerateEmptyMail(title string, body template.HTML) (*gomail.Message, error) {
	message := gomail.NewMessage()
	message.SetHeader("Subject", title)
	message.SetBody("text/html", string(body))

	return message, nil
}
