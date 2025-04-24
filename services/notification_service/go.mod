module authentication_service/notification_service

go 1.24.0

toolchain go1.24.2

replace authentication_service/core => ../../core

require (
	authentication_service/core v0.0.0-00010101000000-000000000000
	github.com/BurntSushi/toml v1.4.0
	github.com/nicksnyder/go-i18n/v2 v2.4.1
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/text v0.22.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require (
	aidanwoods.dev/go-paseto v1.5.4 // indirect
	aidanwoods.dev/go-result v0.3.1 // indirect
	github.com/Masterminds/squirrel v1.5.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.11 // indirect
	gorm.io/gorm v1.25.12 // indirect
)
