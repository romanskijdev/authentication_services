package userhandler

import (
	errm "authentication_service/core/errmodule"
	"authentication_service/core/typescore"
	"authentication_service/rest_user_service/handler"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetProfileHandler Получение профиля пользователя
// @Summary Получение профиля пользователя
// @Description Получение профиля пользователя
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} typescore.User "Успех"
// @Failure 400 {object} handler.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /api/users/profile [get]
func (s *UsersReg) GetProfileHandler(w http.ResponseWriter, r *http.Request) (interface{}, *errm.Error) {
	logrus.Info("🤍 GetProfileHandler")
	ctx := r.Context()

	guidUser, err := handler.GetGuidFromContext(ctx)
	if err != nil {
		return nil, errm.NewError("user_address_not_found", err)
	}

	options := typescore.ListDbOptions{Filtering: &typescore.User{
		SystemID: &guidUser,
	}}

	users, _, errW := s.ipc.DB.Users.GetUsersListDB(ctx, options)
	if errW != nil {
		return nil, errW
	}

	if len(users) == 0 {
		return nil, errm.NewError("not_found", errors.New("not_found"))
	}

	return users[0], nil
}
