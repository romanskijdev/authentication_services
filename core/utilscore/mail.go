package utilscore

import (
	errm "authentication_service/core/errmodule"
	"net/mail"
)

func ValidateEmailFormat(email *string) *errm.Error {
	if email == nil || *email == "" {
		return errm.NewError("email is empty", nil)
	}
	_, err := mail.ParseAddress(*email)
	if err != nil {
		return errm.NewError("failed to validate email", err)
	}
	return nil
}
