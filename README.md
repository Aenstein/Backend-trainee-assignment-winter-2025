# Backend-trainee-assignment-winter-2025
Тестовое задание для стажёра Backend-направления (зимняя волна 2025)

## Структура проекта

```bash
avito-shop/
├── cmd/
│   └── main.go              # Точка входа в приложение
├── internal/
│   ├── auth/                # Логика JWT
│   ├── handlers/            # HTTP-обработчики (auth, info, sendCoin, buy)
│   ├── middleware/          # Middleware для проверки JWT
│   └── services/            # Доступ к БД, бизнес-логика
├── tests/
│   └── integration/         # Интеграционные тесты
├── models/
│   └── models/              # Модели данных (User, запросы, ответы, инвентарь)
├── Dockerfile
├── docker-compose.yaml
└── README.md
```

Для запуска юнит тестов:

go test -v ./internal/... -cover

Для запуска интеграционных тестов:

go test -v ./tests/integration/...
