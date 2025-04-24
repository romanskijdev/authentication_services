package handler

import (
	"authentication_service/core/configcore"
	errm "authentication_service/core/errmodule"
	"authentication_service/core/securecore"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net"
	"net/http"
	"strings"
)

func BasicAuthMiddleware(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || !validateCredentials(user, pass, username, password) {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func validateCredentials(username, password, validUsername, validPassword string) bool {
	return username == validUsername && password == validPassword
}

// Определить пользовательский тип для контекстных ключей
type contextKey string

const guidContextKey contextKey = "guid"

// JWTVerifier middleware для проверки JWT токена
func JWTVerifier(cfg *configcore.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Secrets.AuthJWT.UserSecret == "" {
				errm.NewError("jwt_secret_not_found", errors.New("jwt secret not found"))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errm.NewError("jwt_token_not_found", errors.New("jwt token not found"))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				errm.NewError("jwt_bearer_not_found", errors.New("invalid authorization header"))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Получение IP-адреса клиента
			clientIP := getClientIP(r)
			if clientIP == "" {
				errm.NewError("client_ip_not_found", errors.New("unable to determine client IP"))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Проверка токена
			_, claims, err := securecore.VerifyToken(tokenString, cfg.Secrets.AuthJWT.UserSecret, clientIP, jwt.SigningMethodHS512)
			if err != nil {
				errm.NewError("jwt_token_verification_error", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Сохранение GUID в контексте
			ctx := context.WithValue(r.Context(), guidContextKey, claims["guid"].(string))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetGuidFromContext извлекает GUID из контекста
func GetGuidFromContext(ctx context.Context) (string, error) {
	guid, ok := ctx.Value(guidContextKey).(string)
	if !ok {
		return "", errors.New("no GUID found in context")
	}
	return guid, nil
}

// getClientIP получает IP-адрес клиента из заголовков или RemoteAddr
func getClientIP(r *http.Request) string {
	// Проверка заголовка X-Real-IP
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// Проверка заголовка X-Forwarded-For
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Используем RemoteAddr как fallback
	if ip := r.RemoteAddr; ip != "" {
		// Удаляем порт из адреса
		host, _, err := net.SplitHostPort(ip)
		if err == nil {
			return host
		}
		return ip
	}

	return ""
}
