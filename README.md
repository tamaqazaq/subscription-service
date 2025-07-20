# Subscription Service

REST-сервис на Go для управления онлайн-подписками пользователей. Реализовано по принципам Clean Architecture с поддержкой Swagger, Docker и миграций.

---

## Функциональность

- Создание подписки
- Получение списка всех подписок
- Получение подписки по ID
- Обновление подписки
- Удаление подписки
- Подсчёт общей стоимости подписок за период (`/subscriptions/total`) с фильтрацией по `user_id`, `service_name`

---

## Архитектура

Проект построен по **чистой архитектуре (Clean Architecture)**:

```

subscription-service/
├── cmd/                        # main.go (входная точка)
├── config/                     # конфигурация из .env
├── internal/
│   ├── adapters/http/          # HTTP handlers (Gin)
│   ├── domain/
│   │   ├── application/        # интерфейсы для usecase
│   │   ├── repository/         # интерфейсы для работы с БД
│   │   ├── subscription.go     # модель Subscription (entity)
│   │   └── dateonly.go         # вспомогательные типы
│   ├── infrastructure/
│   │   └── postgres/           # реализация repository через PostgreSQL
│   └── usecase/                # бизнес-логика (usecases)
├── db/migrations/              # SQL миграции (up/down)
├── docs/                       # Swagger-документация
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── Dockerfile                  # Docker-образ приложения
├── docker-compose.yml          # docker + postgres + migrate
├── .dockerignore               # игнорируемые файлы в docker
├── .env.example                # пример env-файла (без чувствит. данных)
├── go.mod / go.sum             # зависимости Go
└── README.md                   # документация проекта


````

---

## Технологии

- **Go 1.21**
- **Gin** — HTTP-фреймворк
- **PostgreSQL 15**
- **Swaggo** — автогенерация Swagger
- **Docker / Docker Compose**
- **golang-migrate** — миграции базы данных

---

## Установка и запуск

### 1. Клонируйте репозиторий
```bash
git clone https://github.com/yourusername/subscription-service.git
cd subscription-service
````

### 2. Создайте `.env` (или используйте `.env.example`)

```env
DB_HOST=db
DB_PORT=5432
DB_USER=
DB_PASSWORD=
DB_NAME=subscriptions
PORT=8080
```

### 3. Запустите через Docker Compose

```bash
docker-compose up --build
```

База данных будет автоматически проинициализирована из `db/migrations/0001_init_subscriptions.up.sql`

---

## Swagger UI

После запуска Swagger доступен по адресу:

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Примеры API-запросов

### POST `/subscriptions`

```json
{
  "service_name": "Netflix",
  "price": 899,
  "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "start_date": "2025-07-01",
  "end_date": "2025-12-31"
}
```

---

### GET `/subscriptions/total?user_id=...&start=07-2025&end=12-2025`

Параметры:

* `user_id` (необязательно)
* `service_name` (необязательно)
* `start` и `end` — **обязательно**, формат `MM-YYYY`

Пример ответа:

```json
{
  "total": 2697
}
```

---

## Тестирование через Postman

Можно импортировать `swagger.json` в Postman или воспользоваться Swagger UI для ручного тестирования.

---

## .dockerignore

```dockerignore
*.log
*.test
*.out
*.exe
*.mod
*.sum

.git
.gitignore
.vscode/
.idea/
.DS_Store
docs/
.env
```

## Автор

Даулет Ермуханов
Контакт: \[[yermukhanovdaulet@gmail.com](mailto:yermukhanovdaulet@gmail.com)]
Telegram: \[@devletin]

---

## Лицензия

Проект создан для тестового задания Junior Golang Developer. Без лицензии.


