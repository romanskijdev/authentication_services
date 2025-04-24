package errm

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Error структура для описания ошибки
type Error struct {
	Code        int    `json:"code"`        // Код ошибки
	Message     string `json:"messages"`    // Сообщение об ошибке
	Description string `json:"description"` // Описание ошибки
	Error       error  `json:"error"`       // Вложенная ошибка
}

func NewError(message string, err error) *Error {
	if err == nil {
		return nil
	}
	code := 0
	description := message
	if err != nil {
		description = fmt.Sprintf("%s: %v", message, err)
	}

	logrus.Errorf("🔴 newError: Error: %s, Description: %s, Error: %v", message, description, err)
	return &Error{
		Code:        code,
		Message:     message,
		Description: description,
		Error:       err,
	}
}
