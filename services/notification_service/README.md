go mod edit -replace sveves-team/zion-crypto-bank/core=../../core
go mod tidy

go build -o service_app cmd/apgo

go run cmd/app.go

swag init -g cmd/app.go


swag init -g /service/rest_user_service/cmd/app.go


swag init -g ./service/rest_user_service/cmd/app.go -o ./service/rest_user_service/swagger
