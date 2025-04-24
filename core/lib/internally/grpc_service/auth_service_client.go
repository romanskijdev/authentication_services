package grpcservice

import (
	"authentication_service/core/lib/external/grpccore"
	protoobj "authentication_service/core/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func InitClientAuthServiceProto(opts []grpc.DialOption, internalHost, port string) protoobj.AuthServiceClient {
	var conn *grpc.ClientConn
	var err error
	list := fmt.Sprintf("%s:%s", internalHost, port)
	// Попытка подключения с повторением
	for {
		conn, err = grpccore.CreateClientConnects(opts, list, false)
		if err != nil {
			log.Printf("🔴 Failed to connect to Telegram Invoice server: %v. Retrying...", err)
			time.Sleep(1 * time.Second) // Задержка перед следующей попыткой
			continue
		}
		// Вывод сообщения о подключении только после успешного соединения
		log.Println("🟢 Telegram Invoice server connected... ", list)
		break // Выход из цикла, если подключение успешно
	}
	return protoobj.NewAuthServiceClient(conn)
}
