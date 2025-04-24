# auth_service

### Инициализация проекта
```bash
go mod init sveves-team/tmail-mail-backend/ton-service
go mod tidy
go mod edit -replace sveves-team/tmail-mail-backend/core=../../core
```


### Локальный старт
```bash
go run cmd/app.go
``` 

### Старт автотестов 
``` bash
go test ./...
```