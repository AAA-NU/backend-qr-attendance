# QR Attendance Service

Сервис для отметки посещаемости студентов с использованием динамически обновляющихся QR-кодов.

## Особенности

-   Генерация QR-кодов каждые 10-15 секунд
-   Веб-интерфейс с автообновлением
-   API для верификации QR-кодов
-   Интеграция с Telegram ботами
-   PostgreSQL для хранения данных

## Установка и запуск

1. Установите зависимости:

```bash
go mod tidy
```

2. Настройте переменные окружения:

```bash
export DATABASE_URL="postgres://username:password@localhost:5432/qr_attendance?sslmode=disable"
export BOT_USERNAME="your_bot_username"
export QR_LIFETIME_SECONDS="12"
export PORT="8080"
```

3. Запустите сервис:

```bash
go run cmd/main.go
```

## API Endpoints

-   `GET /` - Главная страница с QR-кодом
-   `GET /qr/current` - Получить текущий QR-код (изображение)
-   `POST /api/verify/:uuid` - Проверить валидность QR-кода

## Интеграция с Telegram Bot

QR-код содержит ссылку вида: `https://t.me/your_bot_username?start=UUID`

Пример запроса для верификации из бота:

```bash
POST /api/verify/550e8400-e29b-41d4-a716-446655440000
```

Ответ:

```json
{
    "valid": true,
    "uuid": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Структура проекта

```
qr-attendance/
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   └── qr_handler.go
│   ├── models/
│   │   └── qr.go
│   └── services/
│       └── qr_service.go
├── web/
│   └── templates/
│       └── index.html
├── go.mod
└── README.md
```
