# 🐱 Meobeo Talk API
Modern REST API cho ứng dụng mèo béoooo
---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### 1. Clone repository
```bash
git clone https://github.com/yourusername/meobeo-talk-api.git
cd meobeo-talk-api
```

### 2. Setup environment
```bash
cp .env.example .env
# Edit .env với config của bạn
```

### 3. Start database
```bash
make docker-up
```

### 4. Run migrations
```bash
make migrate-up
```

### 5. Run application
```bash
make run
```

## 🏗️ Cấu trúc dự án
```
meobeo-talk-api/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── database.go
│   │
│   ├── domain/
│   │   ├── entity/  
│   │   │   ├── user.go
│   │   │   ├── conversation.go
│   │   │   └── message.go
│   │   ├── repository/  
│   │   │   ├── user_repository.go
│   │   │   └── message_repository.go
│   │   ├── service/ 
│   │   │   └── user_service.go
│   │   └── errors.go
│   │
│   ├── dto/
│   │   ├── request/
│   │   │   ├── user_request.go
│   │   │   └── pagination_request.go
│   │   └── response/
│   │       ├── user_response.go
│   │       └── pagination_response.go
│   │
│   ├── infrastructure/
│   │   ├── persistence/ 
│   │   │   ├── postgres/
│   │   │   │   ├── user_repository.go
│   │   │   │   └── transaction.go
│   │   │   └── redis/  
│   │   │       └── user_cache.go
│   │   └── http/  
│   │
│   ├── application/ 
│   │   ├── user_service.go
│   │   └── auth_service.go
│   │
│   ├── interfaces/   
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   ├── user_handler.go
│   │   │   │   └── auth_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── transaction.go 
│   │   │   │   └── rate_limit.go
│   │   │   └── routes/
│   │   │       ├── router.go
│   │   │       └── api_v1.go
│   │   └── grpc/
│   │
│   └── pkg/    
│       ├── pagination/
│       │   └── paginator.go
│       ├── response/
│       │   └── json.go
│       └── validator/
│           └── validator.go
│
├── pkg/              
│   ├── logger/
│   ├── jwt/
│   └── errors/               
│
├── migrations/
├── scripts/
├── tests/
├── docs/
└── ...
```