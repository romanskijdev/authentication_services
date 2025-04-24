package grpccore

import (
	"authentication_service/core/variables"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// CreateGRPCSystemCertPoolOptions_x509 создает опции для gRPC клиента, используя сертификаты x509
func CreateGRPCSystemCertPoolOptions_x509(certPathsClient []string, useClientCertificatesTLS bool) ([]grpc.DialOption, error) {
	// logrus.Info("🟨 CreateGRPCSystemCertPoolOptions_x509")
	var err error
	var optionsGrpc []grpc.DialOption

	// Если требуется использовать клиентские сертификаты TLS
	if useClientCertificatesTLS {
		// Если предоставлены пути к сертификатам
		if certPathsClient != nil {
			certPool := x509.NewCertPool() // Создание нового пула сертификатов
			// Цикл по всем путям к сертификатам
			for _, certPath := range certPathsClient {
				var cert []byte
				// Чтение сертификата
				if cert, err = os.ReadFile(certPath); err == nil {
					// Добавление сертификата в пул
					resultAddCertificate := certPool.AppendCertsFromPEM(cert)
					// Если сертификат не добавлен в пул, вывод ошибки
					if !resultAddCertificate {
						fmt.Printf("Error failed to append cert to pool from %s: %s\n", certPath, err)
					}
				}
			}
			// Создание TLS credentials с использованием пула сертификатов
			creds := credentials.NewTLS(&tls.Config{RootCAs: certPool})
			// Добавление TLS credentials в опции gRPC
			optionsGrpc = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
		} else {
			// Если пути к сертификатам не предоставлены, использование системного пула сертификатов
			var systemRoots *x509.CertPool
			if systemRoots, err = x509.SystemCertPool(); err == nil {
				optionsGrpc = []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{RootCAs: systemRoots}))}
			}
		}
	} else {
		// Если не требуется использовать клиентские сертификаты TLS, использование небезопасного соединения
		optionsGrpc = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
	return optionsGrpc, err // Возвращение опций gRPC и возможной ошибки
}

// CreateServerGRPC создает gRPC сервер с опциональными TLS сертификатами
func CreateServerGRPC(certificateBaseCrtFilePath *string, certificateBaseKeyFilePath *string) (*grpc.Server, error) {
	// logrus.Info("🟨 CreateServerGRPC")
	var err error
	var cert tls.Certificate
	var certPool *x509.CertPool

	// Начальные опции сервера с установленным максимальным размером сообщения
	serverOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(variables.MaxMsgGRPCSize),
		grpc.MaxSendMsgSize(variables.MaxMsgGRPCSize),
	}

	// Если требуется использовать TLS сертификаты сервера
	if certificateBaseCrtFilePath != nil && certificateBaseKeyFilePath != nil {
		// Загрузка системного пула сертификатов
		if certPool, err = x509.SystemCertPool(); err == nil {
			// Загрузка пары ключей сертификата
			if cert, err = tls.LoadX509KeyPair(*certificateBaseCrtFilePath, *certificateBaseKeyFilePath); err == nil {
				// Создание конфигурации TLS с загруженными сертификатами
				tlsConfig := &tls.Config{
					Certificates: []tls.Certificate{cert},
					RootCAs:      certPool,
				}
				// Добавление TLS credentials в опции сервера
				serverOptions = append(serverOptions, grpc.Creds(credentials.NewTLS(tlsConfig)))
			}
		}
	}

	kaep := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Minute, // Минимальное время между pings от клиента
		PermitWithoutStream: true,            // Разрешить pings без активных потоков
	}
	kasp := keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Minute, // Максимальное время простоя соединения
		MaxConnectionAge:      30 * time.Minute, // Максимальное время жизни соединения
		MaxConnectionAgeGrace: 5 * time.Minute,  // Дополнительное время после достижения MaxConnectionAge
		Time:                  5 * time.Minute,  // Время между pings от сервера
		Timeout:               20 * time.Second, // Таймаут для ответа на ping
	}

	serverOptions = append(serverOptions,
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.MaxConcurrentStreams(100))

	// Создание нового gRPC сервера с указанными опциями
	return grpc.NewServer(serverOptions...), err
}
