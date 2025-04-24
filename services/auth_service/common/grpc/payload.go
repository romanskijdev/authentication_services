package grpcpayment

import (
	rabbitmqlib "authentication_service/core/lib/external/rabbitmq-lib"
	protoobj "authentication_service/core/proto"
	"authentication_service/core/securecore"
	"authentication_service/core/typescore"
	"authentication_service/core/variables"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	accessTokenLifeTime  = time.Minute * 15
	refreshTokenLifeTime = time.Hour * 24 * 7
)

func (s *AuthServiceServiceProto) IssueTokens(ctx context.Context, req *protoobj.IssueTokensRequest) (*protoobj.IssueTokensResponse, error) {
	if s.ipc == nil {
		logrus.Error("module is nil")
		return nil, errors.New("module is nil")
	}

	// Проверка входных данных
	userID := req.GetUserId()
	clientIP := req.GetClientIp()
	if userID == "" || clientIP == "" {
		logrus.Error("invalid input: user_id or client_ip is empty")
		return nil, status.Error(codes.InvalidArgument, "user_id and client_ip are required")
	}

	// Генерация Access токена
	accessToken, err := securecore.GenerateTokenJWT(
		userID,
		clientIP,
		s.ipc.Config.Secrets.AuthJWT.UserSecret,
		accessTokenLifeTime, // Время жизни Access токена
		jwt.SigningMethodHS512,
	)
	if err != nil {
		logrus.Errorf("failed to generate access token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	// Генерация Refresh токена
	refreshToken, err := securecore.GenerateTokenJWT(
		userID,
		clientIP,
		s.ipc.Config.Secrets.AuthJWT.UserSecret,
		refreshTokenLifeTime,
		jwt.SigningMethodHS512,
	)
	if err != nil {
		logrus.Errorf("failed to generate refresh token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	// Возвращаем ответ
	return &protoobj.IssueTokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServiceProto) RefreshTokens(ctx context.Context, req *protoobj.RefreshTokensRequest) (*protoobj.RefreshTokensResponse, error) {
	if s.ipc == nil {
		logrus.Error("module is nil")
		return nil, errors.New("module is nil")
	}

	// Проверка входных данных
	refreshToken := req.GetRefreshToken()
	clientIP := req.GetClientIp()
	if refreshToken == "" || clientIP == "" {
		logrus.Error("invalid input: refresh_token or client_ip is empty")
		return nil, status.Error(codes.InvalidArgument, "refresh_token and client_ip are required")
	}

	// Проверка Refresh токена
	_, claims, err := securecore.VerifyToken(
		refreshToken,
		s.ipc.Config.Secrets.AuthJWT.UserSecret,
		clientIP,
		jwt.SigningMethodHS512,
	)
	if err != nil {
		logrus.Errorf("failed to verify refresh token: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// Извлечение GUID из payload
	userID, ok := claims["guid"].(string)
	if !ok || userID == "" {
		logrus.Error("failed to extract user_id from refresh token claims")
		return nil, status.Error(codes.Internal, "failed to extract user_id from token")
	}

	// Проверка изменения IP-адреса
	tokenIP, ok := claims["client_ip"].(string)
	if !ok || tokenIP != clientIP {
		// Отправка email warning (если требуется)
		sendTo := []*string{&userID}
		category := typescore.DeviceNewNotifyCategory
		notify := &typescore.NotifyParams{
			IsEmail:   true,
			Emergency: true,
			UsersIDs:  sendTo,
			Category:  &category,
		}

		err := rabbitmqlib.PublishMessage(s.ipc.RabbitMQ,
			variables.RabbitMQExchangeNotifications,
			variables.RabbitMQNotificationsServiceRoute,
			notify)
		if err != nil {
			logrus.Errorf("failed to send notification %v", err)
		}

		return nil, status.Error(codes.PermissionDenied, "client IP mismatch")
	}

	// Генерация новой пары токенов
	newAccessToken, err := securecore.GenerateTokenJWT(
		userID,
		clientIP,
		s.ipc.Config.Secrets.AuthJWT.UserSecret,
		accessTokenLifeTime,
		jwt.SigningMethodHS512,
	)
	if err != nil {
		logrus.Errorf("failed to generate new access token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate new access token")
	}

	newRefreshToken, err := securecore.GenerateTokenJWT(
		userID,
		clientIP,
		s.ipc.Config.Secrets.AuthJWT.UserSecret,
		refreshTokenLifeTime,
		jwt.SigningMethodHS512,
	)
	if err != nil {
		logrus.Errorf("failed to generate new refresh token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate new refresh token")
	}

	// Возвращаем ответ
	return &protoobj.RefreshTokensResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		IpChanged:    tokenIP != clientIP,
	}, nil
}
