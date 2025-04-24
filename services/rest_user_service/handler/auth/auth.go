package authhandler

import (
	errm "authentication_service/core/errmodule"
	protoobj "authentication_service/core/proto"
	"authentication_service/rest_user_service/handler"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type GetTokensPairReq struct {
	UserID *string `json:"user_id"`
}

// RefreshTokensHandler Обновление токенов
// @Summary Обновление токенов
// @Description Обновляет Access и Refresh токены на основе действующего Refresh токена
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <refresh_token>"
// @Success 200 {object} typescore.TokenPair "Успех"
// @Failure 400 {object} handler.ErrorResponse "Некорректный запрос"
// @Failure 401 {object} handler.ErrorResponse "Невалидный или просроченный токен"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /api/auth/refresh [post]
func (s *AuthReg) RefreshTokensHandler(w http.ResponseWriter, r *http.Request) (interface{}, *errm.Error) {
	logrus.Info("🤍 RefreshTokensHandler")
	ctx := r.Context()

	// Получение Refresh токена из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errm.NewError("refresh_token_not_found", errors.New("refresh token not found"))
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errm.NewError("invalid_authorization_header", errors.New("invalid authorization header"))
	}

	refreshToken := parts[1]

	// Получение IP-адреса клиента
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	// Проверка и обновление токенов
	newTokenPair, err := s.ipc.ClientAuthServiceProto.RefreshTokens(ctx, &protoobj.RefreshTokensRequest{ClientIp: ip, RefreshToken: refreshToken})
	if err != nil {
		return nil, errm.NewError("token_generation_error", err)
	}

	return newTokenPair, nil
}

// IssueTokensHandler Выдача пары токенов (Access + Refresh)
// @Summary Выдача пары токенов
// @Description Выдает Access и Refresh токены для пользователя с указанным GUID
// @Tags auth
// @Accept json
// @Produce json
// @Param guid query string true "GUID пользователя"
// @Success 200 {object} typescore.TokenPair "Успех"
// @Failure 400 {object} handler.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} handler.ErrorResponse "Ошибка сервера"
// @Router /api/auth/issue [post]
func (s *AuthReg) IssueTokensHandler(w http.ResponseWriter, r *http.Request) (interface{}, *errm.Error) {
	logrus.Info("🤍 IssueTokensHandler")
	ctx := r.Context()

	tokenReq := &GetTokensPairReq{}
	if errObj := handler.ParseRequestBodyPost(r, tokenReq); errObj != nil {
		return nil, errObj
	}
	if tokenReq == nil || tokenReq.UserID == nil {
		return nil, errm.NewError("empty_obj", errors.New("empty_obj"))
	}

	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	// Генерация токенов
	accessToken, err := s.ipc.ClientAuthServiceProto.IssueTokens(ctx, &protoobj.IssueTokensRequest{UserId: *tokenReq.UserID, ClientIp: ip})
	if err != nil {
		return nil, errm.NewError("token_generation_error", err)
	}

	return accessToken, nil
}
