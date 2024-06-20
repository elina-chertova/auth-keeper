## Запуск тестов
```shell
go test ./...
```

## Запуск сервера
Необходимо создать .env файл (аналогично .env.example) в корне проекта.
Пример:
```shell
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_NAME=db
DB_PASSWORD=password

APP_ADDRESS=localhost:8080
SECRET_KEY=key
```

```shell
go build cmd/server/main.go
go run cmd/server/main.go
```

# Usage
## Registration
```shell
go run cmd/client/main.go register --username testuser --password testpass --email test@example.com 
```
## Login
```shell
go run cmd/client/main.go login --username testuser --password testpass --email test@example.com 
```


## Add Card
```shell
go run cmd/client/main.go add-card --card_number 1234567890123456 --expiry_date 12/25 --cvv 123  --card_holder "John Doe" --metadata "Some metadata" --token token

```
## Get Card
```shell
go run cmd/client/main.go get-card --token token
```

## Add Text data
```shell
go run cmd/client/main.go add-text-data --content "hello data" --metadata "Some metadata" --token
```

## Get Text data
```shell
go run cmd/client/main.go get-text-data --token token
```

## Add Binary data
```shell
go run cmd/client/main.go add-binary-data --file_path /path/to/data/1.jpg --token token

```
## Get Binary data
```shell
go run cmd/client/main.go get-binary-data --token token
```

## Add Login Password
```shell
go run cmd/client/main.go add-login-password --username rocketman --password 123qwe --token 
```
## Get Login Password
```shell
go run cmd/client/main.go get-login-password --token token
```
