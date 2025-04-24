# system_service

### Инициализация проекта
```bash
go mod init authentication_service/system_service
go mod tidy
go mod edit -replace authentication_service/core => ../../core
```

### Локальный старт
```bash
go run cmd/app.go
```

### Старт автотестов 
``` bash
go test ./...
```

``` bash
go vet ./...
```

```bash
gofmt -w .
```

```bash
swag init -g cmd/app.go --dir ./,../../core
```