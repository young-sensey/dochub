# DocFlow

Система управления документами с REST API на Golang, PostgreSQL, React.js и Nginx.

## Структура проекта

```
DocFlow/
├── compose.yaml              # Docker Compose конфигурация
├── Makefile                  # Команды для сборки и тестирования
├── backend/                  # Golang backend
│   ├── main.go              # Точка входа приложения
│   ├── go.mod               # Go модули и зависимости
│   ├── go.sum               # Хеши зависимостей
│   ├── Dockerfile           # Docker образ для backend
│   ├── handlers/            # HTTP обработчики
│   │   ├── document_handlers.go
│   │   └── document_handlers_test.go
│   ├── database/            # Работа с базой данных
│   │   └── database.go
│   ├── routes/              # Настройка маршрутов
│   │   └── routes.go
│   └── integration_test.go  # Интеграционные тесты
├── frontend/                # React.js frontend
│   ├── Dockerfile           # Docker образ для frontend
│   ├── package.json         # Node.js зависимости
│   ├── public/              # Статические файлы
│   │   └── index.html
│   └── src/                 # Исходный код React
│       ├── App.js           # Основной компонент
│       ├── App.css          # Стили приложения
│       ├── index.js         # Точка входа
│       └── index.css        # Глобальные стили
├── nginx/                   # Nginx конфигурация
│   └── nginx.conf
└── test_api.sh             # Тестовый скрипт для API
```

## API Endpoints

### Документы (`/dock`)

- `GET /dock` - Получить список всех документов
- `POST /dock` - Создать новый документ
- `GET /dock/{id}` - Получить документ по ID
- `PUT /dock/{id}` - Обновить документ по ID
- `DELETE /dock/{id}` - Удалить документ по ID

### Health Check

- `GET /health` - Проверка состояния сервиса

## Запуск проекта

1. **Запуск всех сервисов:**
   ```bash
   make run
   # или
   docker compose up --build
   ```

2. **Запуск в фоновом режиме:**
   ```bash
   make run-detached
   # или
   docker compose up -d --build
   ```

3. **Остановка сервисов:**
   ```bash
   make stop
   # или
   docker compose down
   ```

## Доступ к приложению

После запуска приложение будет доступно по адресам:

- **Frontend (React.js)**: http://localhost:80 или http://localhost:3000
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432

## Функциональность Frontend

React.js интерфейс предоставляет:

- **Создание документов** - форма для добавления новых документов
- **Просмотр документов** - список всех документов с информацией
- **Редактирование документов** - возможность изменения существующих документов
- **Удаление документов** - удаление документов с подтверждением
- **Уведомления** - отображение успешных операций и ошибок

## Тестирование

### Unit тесты
```bash
make test-unit
# или
cd backend && go test ./handlers -v
```

### Интеграционные тесты
```bash
# Сначала запустите проект
make run-detached

# Затем запустите интеграционные тесты
make test-integration
# или
cd backend && go test -v -run "TestFullCRUDWorkflow|TestHealthEndpoint|TestInvalidRequests"
```

### Все тесты
```bash
make test
```

### Тестирование API через curl
```bash
make test-api
# или
chmod +x test_api.sh && ./test_api.sh
```

## Полезные команды

```bash
# Сборка проекта
make build

# Просмотр логов
make logs

# Просмотр логов конкретного сервиса
make logs-backend
make logs-frontend
make logs-db
make logs-nginx

# Очистка
make clean
```

## Переменные окружения

### Backend
- `DB_HOST` - Хост PostgreSQL (по умолчанию: localhost)
- `DB_USER` - Пользователь PostgreSQL (по умолчанию: docflow)
- `DB_PASSWORD` - Пароль PostgreSQL (по умолчанию: docflow_pass)
- `DB_NAME` - Имя базы данных (по умолчанию: docflow_db)

### Frontend
- `REACT_APP_API_URL` - URL API backend (по умолчанию: http://localhost:8080)

## Порты

- `80` - Nginx (прокси на frontend и backend)
- `3000` - Frontend (React.js)
- `8080` - Backend API
- `5432` - PostgreSQL

## Структура базы данных

### Таблица `documents`

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL PRIMARY KEY | Уникальный идентификатор |
| title | VARCHAR(255) | Заголовок документа |
| content | TEXT | Содержимое документа |
| author | VARCHAR(255) | Автор документа |
| created_at | TIMESTAMP | Дата создания |
| updated_at | TIMESTAMP | Дата обновления | 