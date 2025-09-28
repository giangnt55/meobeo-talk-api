# Project Structure
```
myapp/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── connection.go
│   │   └── migration.go
│   ├── models/
│   │   └── user.go
│   ├── domain/
│   │   ├── user/
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── handler.go
│   │   └── common/
│   │       ├── response.go
│   │       ├── request.go
│   │       └── messages.go
│   ├── dto/
│   │   ├── request/
│   │   │   ├── user_request.go
│   │   │   └── pagination.go
│   │   └── response/
│   │       ├── user_response.go
│   │       └── pagination_response.go
│   ├── middleware/
│   │   ├── cors.go
│   │   ├── error_handler.go
│   │   └── logger.go
│   ├── router/
│   │   └── router.go
│   └── utils/
│       ├── validator.go
│       └── pagination.go
├── pkg/
│   └── logger/
│       └── logger.go
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```