# GoMeal

> **English** | [中文](README.md)

A food-delivery (takeout) ordering system written in Go — a learning project.

## Tech stack

- **Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL 8.0
- **Cache**: Redis 7
- **Auth**: JWT (with RBAC)
- **Docs**: Swagger
- **Deployment**: Docker Compose

## Features

- User registration / login and profile management (JWT)
- Restaurants (merchants) and their menus: list, detail, create, update, availability
- Shopping cart: add, update, remove, clear
- Delivery addresses: CRUD per user
- Orders: create, list, detail, cancel, pay, query payment status
- Merchant onboarding and merchant-side menu management
- Middleware: JWT auth, RBAC, CORS, request logging, rate limiting
- Redis-backed stock/cache layer

## Quick start

```bash
# configure passwords / secrets
cp .env.example .env

# start MySQL + Redis + app
docker compose up -d

# API is served at http://localhost:8080
```

## Project layout

```
backend/             # Go backend
├── cmd/server/      # entry point
├── internal/
│   ├── handler/     # HTTP handlers
│   ├── service/     # business logic
│   ├── repository/  # data access
│   ├── model/       # data models
│   ├── middleware/  # auth / rate limiting
│   ├── cache/       # Redis cache
│   └── config/      # configuration
├── docs/            # Swagger docs
├── Dockerfile
└── docker-compose.yml

frontend/            # Vue 3 frontend
```

## API docs

See [API.md](backend/API.md), or open the Swagger UI once the service is running.

## License

Released under [GPL-3.0](LICENSE).
