package main

import (
	"authentication_service/core/configcore"
	"authentication_service/core/database"
	errm "authentication_service/core/errmodule"
	"authentication_service/core/lib/external/grpccore"
	pgxpoolmodule "authentication_service/core/lib/external/pgxpool"
	rabbitmqlib "authentication_service/core/lib/external/rabbitmq-lib"
	grpcservice "authentication_service/core/lib/internally/grpc_service"
	_ "authentication_service/rest_user_service/docs"
	"authentication_service/rest_user_service/handler"
	authhandler "authentication_service/rest_user_service/handler/auth"
	userhandler "authentication_service/rest_user_service/handler/profile"
	typesm "authentication_service/rest_user_service/types"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"time"
)

const (
	readTimeout  = 15 * time.Minute // Таймаут чтения
	writeTimeout = 15 * time.Minute // Таймаут записи
	idleTimeout  = 15 * time.Minute // Таймаут бездействия
	rateLimit    = 25               // Лимит запросов
	rateWindow   = 1 * time.Second  // Окно времени для лимита запросов
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		logrus.Errorln("❌ Failed to load config: ", err)
		select {}
	}

	ipc, err := initInternalProvider(cfg)
	if err != nil {
		logrus.Errorln("❌ Failed to init internally provider: ", err)
		select {}
	}

	ipc, err = registerGrpcServices(ipc)
	if err != nil {
		logrus.Errorln("❌ Failed to register grpc lib: ", err)
		return
	}

	router, err := initBaseApiRouter(ipc)
	if err != nil {
		logrus.Errorln("❌ Failed to init base api router: ", err)
		select {}
	}

	if err := startRestApiServer(ipc, router); err != nil {
		logrus.Errorln("❌ Failed to start server: ", err)
		select {}
	}
}

func getConfig() (*configcore.Config, error) {
	options := &configcore.ConfigLoadOptions{
		Database: true,
		GrpsClients: configcore.GrpsClientsOptions{
			AuthService: true,
		},
		ExposedServiceConfig: configcore.ExposedServiceOptions{
			UserService: true,
		},
		RabbitMQConfig: true,
		Secrets: configcore.SecretsOptions{
			User: true,
		},
	}

	cfg, err := configcore.LoadConfig(options)
	if err != nil {
		logrus.Errorln("❌ Failed to load config: ", err)
		return nil, err
	}
	return cfg, nil
}

func initInternalProvider(appConfig *configcore.Config) (*typesm.InternalProviderControl, error) {
	db, err := database.NewModuleDB(&pgxpoolmodule.ConfigConnectPgxPool{
		Host:     appConfig.Database.Host,
		Port:     fmt.Sprintf("%d", appConfig.Database.Port),
		User:     appConfig.Database.User,
		Password: appConfig.Database.Password,
		Name:     appConfig.Database.DB,
		SSLMode:  appConfig.Database.SSL,
	})
	if err != nil {
		logrus.Errorln("❌ Failed to init database: ", err)
		return nil, err
	}

	rabbitMQClient, err := rabbitmqlib.ConnectRabitMq(rabbitmqlib.ConnectParams{
		Username: appConfig.RabbitMQConfig.User,
		Host:     appConfig.RabbitMQConfig.Host,
		Port:     appConfig.RabbitMQConfig.Port,
		Password: appConfig.RabbitMQConfig.Password,
	})
	if err != nil || rabbitMQClient == nil {
		logrus.Error("❌ Failed to init RabbitMQ client: ", err)
		return nil, err
	}

	return &typesm.InternalProviderControl{
		Config:   appConfig,
		DB:       db,
		RabbitMQ: rabbitMQClient,
	}, nil
}

func startRestApiServer(ipc *typesm.InternalProviderControl, router *chi.Mux) error {
	serverPort := fmt.Sprintf("%d", ipc.Config.ExposedServiceConfig.UserService.PortRest)

	serverAddr := ":" + serverPort
	httpServer := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Printf("✅ Starting server on %s", serverAddr)
	return httpServer.ListenAndServe()
}

func initBaseApiRouter(ipc *typesm.InternalProviderControl) (*chi.Mux, error) {
	router := chi.NewRouter()
	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   ipc.Config.ExposedServiceConfig.UserService.Cors.AllowedOrigins, // Список разрешенных origin
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Cache-Control", "X-Requested-With"},
		ExposedHeaders:   []string{"Link", "Cache-Control"},
		AllowCredentials: false,
		MaxAge:           300, // Максимальное время жизни предварительных запросов в секундах
	})

	router.Use(corsOptions.Handler)
	router.Use(middleware.RealIP)

	// Добавляем ограничение по IP
	router.Use(httprate.LimitByIP(rateLimit, rateWindow))

	if ipc.DB == nil {
		logrus.Errorln("❌ Failed to connect to database")
		return nil, errors.New("❌ Failed to connect to database")
	}

	err := registerRoutes(router, ipc)
	if err != nil {
		logrus.Errorln("❌ Failed to register routes")
		return nil, err
	}

	return router, nil
}

func registerGrpcServices(ipc *typesm.InternalProviderControl) (*typesm.InternalProviderControl, error) {
	protoOpt := grpccore.CreateDialOptionsProto()

	clientTelegramInvoiceServiceProto := grpcservice.InitClientAuthServiceProto(
		protoOpt,
		ipc.Config.GrpsClients.AuthService.Host,
		fmt.Sprintf("%d", ipc.Config.GrpsClients.AuthService.Port),
	)

	ipc.ClientAuthServiceProto = clientTelegramInvoiceServiceProto
	return ipc, nil
}

// регистрирует маршруты
func registerRoutes(router chi.Router, ipc *typesm.InternalProviderControl) error {
	// Swagger endpoint
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Применение middleware для базовой аутентификации к маршруту Swagger UI
	router.Group(func(r chi.Router) {
		logrus.Infof("🔒 Swagger UI is protected by basic auth: %s %s",
			ipc.Config.ExposedServiceConfig.UserService.Swagger.User,
			ipc.Config.ExposedServiceConfig.UserService.Swagger.Pass)
		r.Use(handler.BasicAuthMiddleware(ipc.Config.ExposedServiceConfig.UserService.Swagger.User, ipc.Config.ExposedServiceConfig.UserService.Swagger.Pass))
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	// Слайс обработчиков маршрутов
	routeHandlers := []struct {
		name    string
		handler func(chi.Router, *typesm.InternalProviderControl) *errm.Error
	}{
		{"auth", authhandler.RegisterAuthRoutes},
		{"user", userhandler.RegisterUsersRoutes},
	}

	// Регистрация всех маршрутов
	for _, rh := range routeHandlers {
		if errObj := rh.handler(router, ipc); errObj != nil {
			logrus.Errorln("❌ Failed to register", rh.name, "routes:", errObj)
			return fmt.Errorf("❌ Failed to register %s routes", rh.name)
		}
	}

	// Вывод зарегистрированных маршрутов
	if err := printRoutes(router); err != nil {
		logrus.Errorln("❌ Failed to print routes:", err)
		return fmt.Errorf("❌ Failed to print routes: %w", err)
	}

	return nil
}

func printRoutes(r chi.Router) error {
	err := chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("🩵 [ %s ] %s\n", method, route)
		return nil
	})
	if err != nil {
		logrus.Errorln("failed to walk routes: ", err)
		return err
	}
	return nil
}
