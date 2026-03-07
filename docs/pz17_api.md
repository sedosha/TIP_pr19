# Практическое занятие №1
## Разделение монолита на 2 микросервиса. Взаимодействие через HTTP

### 1. Auth Service (порт 8089)

#### Переменные окружения
| Переменная | Значение по умолчанию |
|------------|----------------------|
| `AUTH_PORT` | 8089 |

#### Эндпоинты

POST /v1/auth/login - получение токена

curl -X POST http://localhost:8089/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: req-001" \
  -d '{"username":"student","password":"student"}'

Response 200:
{
  "access_token": "demo-token",
  "token_type": "Bearer"
}

Response 401:
{"error":"invalid credentials"}

GET /v1/auth/verify - проверка токена

curl -X GET http://localhost:8089/v1/auth/verify \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-002"

Response 200:
{
  "valid": true,
  "subject": "user"
}

Response 401:
{"valid":false,"error":"invalid token"}

GET /health - проверка сервиса

curl http://localhost:8089/health

Response 200:
{"status":"ok","service":"auth"}

### 2. Tasks Service (порт 8090)

#### Переменные окружения
| Переменная | Значение по умолчанию |
|------------|----------------------|
| `TASKS_PORT` | 8090 |
| `AUTH_BASE_URL` | http://localhost:8089 |

#### Эндпоинты (все требуют Authorization: Bearer demo-token)

GET /health

curl http://localhost:8090/health

Response 200:
{"status":"ok","service":"tasks"}

POST /v1/tasks - создание задачи

curl -X POST http://localhost:8090/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-003" \
  -d '{"title":"Do PZ17","description":"split services","due_date":"2026-01-10"}'

Response 201:
{
  "id": "abc12345",
  "title": "Do PZ17",
  "description": "split services",
  "due_date": "2026-01-10",
  "done": false,
  "created_at": "2026-03-08T01:15:23Z"
}

GET /v1/tasks - список задач

curl http://localhost:8090/v1/tasks \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-004"

Response 200:
[
  {
    "id": "abc12345",
    "title": "Do PZ17",
    "done": false
  }
]

GET /v1/tasks/{id} - задача по ID

curl http://localhost:8090/v1/tasks/abc12345 \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-005"

Response 200: аналогично POST  
Response 404: {"error":"task not found"}

PATCH /v1/tasks/{id} - обновление задачи

curl -X PATCH http://localhost:8090/v1/tasks/abc12345 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-006" \
  -d '{"done":true}'

Response 200: обновленная задача

DELETE /v1/tasks/{id} - удаление задачи

curl -X DELETE http://localhost:8090/v1/tasks/abc12345 \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: req-007"

Response 204 (без тела)

### 3. Коды ответов

| Код | Описание |
|-----|----------|
| 200 | Успех |
| 201 | Создано |
| 204 | Нет содержимого |
| 400 | Неверный запрос |
| 401 | Неавторизован |
| 404 | Не найдено |
| 503 | Сервис недоступен |

### 4. Тестирование

# 1. Получить токен
curl -s -X POST http://localhost:8089/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"student","password":"student"}'

# 2. Создать задачу
curl -s -X POST http://localhost:8090/v1/tasks \
  -H "Authorization: Bearer demo-token" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test task"}'

# 3. Получить список
curl -s http://localhost:8090/v1/tasks \
  -H "Authorization: Bearer demo-token"

### 5. Запуск

Auth Service:

cd services/auth
export AUTH_PORT=8089
go run ./cmd/auth

Tasks Service:

cd services/tasks
export TASKS_PORT=8090
export AUTH_BASE_URL=http://localhost:8089
go run ./cmd/tasks
