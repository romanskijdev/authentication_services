package securecore

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(
	guid string,
	clientIP string,
	secretJWT string,
	expirationTime time.Duration,
	signingMethod *jwt.SigningMethodHMAC,
) (string, error) {
	token := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
		"guid":      guid,
		"client_ip": clientIP,
		"exp":       time.Now().Add(expirationTime).Unix(),
	})

	t, err := token.SignedString([]byte(secretJWT))
	if err != nil {
		return "", err
	}

	return t, nil
}

// VerifyToken верифицирует JWT токен и возвращает данные по токену или ошибку
func VerifyToken(
	tokenString string,
	jwtSecret string,
	clientIP string,
	signingMethod *jwt.SigningMethodHMAC,
) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != signingMethod {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	// Извлекаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, errors.New("invalid token claims")
	}

	// Проверка, что ExpiresAt не истек
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, nil, errors.New("token expired")
		}
	} else {
		return nil, nil, errors.New("missing expiration time in token")
	}

	// Проверка IP-адреса клиента
	if ip, ok := claims["client_ip"].(string); ok {
		if ip != clientIP {
			return nil, nil, errors.New("client IP mismatch")
		}
	} else {
		return nil, nil, errors.New("missing client IP in token")
	}

	return token, claims, nil
}
