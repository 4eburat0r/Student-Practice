<<<<<<< HEAD

# 📁 Структура проекта «Практикуйся» — Микросервисная архитектура

## 🏗️ Общая структура проекта

```
backend/
│
├── services/                          # Все микросервисы
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── vacancy-service/
│   ├── resume-service/
│   └── response-service/
│
├── shared/                            # Общий переиспользуемый код
│
├── deployments/                       # Docker Compose и Kubernetes
│
├── scripts/                           # Скрипты для управления проектом
│
├── docs/                              # Документация
│
├── .gitlab-ci.yml                     # CI/CD пайплайн
├── Makefile                           # Общие команды
├── README.md
└── .gitignore
```

---

## 🔹 API Gateway (`services/api-gateway/`)

**Назначение:** Единая точка входа для всех клиентских запросов. Принимает REST запросы, проксирует их в микросервисы, агрегирует ответы.

### Структура:

```
services/api-gateway/
├── cmd/
│   └── api/
│       └── main.go                    # Точка входа приложения
│
├── internal/
│   ├── app/
│   │   ├── app.go                    # Структура App со всеми зависимостями
│   │   └── dependencies.go           # Инициализация зависимостей
│   │
│   ├── config/
│   │   └── config.go                 # Загрузка конфигурации из env
│   │
│   ├── handler/
│   │   ├── middleware/
│   │   │   ├── cors.go              # CORS для фронтенда
│   │   │   ├── logger.go            # Логирование запросов
│   │   │   ├── recovery.go          # Panic recovery
│   │   │   ├── rate_limiter.go      # Rate limiting
│   │   │   └── auth.go              # Валидация JWT (через auth-service)
│   │   │
│   │   ├── auth.go                  # Проксирование в auth-service
│   │   ├── user.go                  # Проксирование в user-service
│   │   ├── vacancy.go               # Проксирование в vacancy-service
│   │   ├── resume.go                # Проксирование в resume-service
│   │   ├── response.go              # Проксирование в response-service
│   │   ├── handler.go               # Главная структура Handler
│   │   └── routes.go                # Регистрация всех маршрутов
│   │
│   ├── client/                       # HTTP клиенты для микросервисов
│   │   ├── auth_client.go
│   │   ├── user_client.go
│   │   ├── vacancy_client.go
│   │   ├── resume_client.go
│   │   └── response_client.go
│   │
│   ├── resilience/                   # Паттерны отказоустойчивости
│   │   ├── circuit_breaker.go       # Circuit Breaker
│   │   ├── retry.go                 # Retry механизм
│   │   └── fallback.go              # Fallback стратегии
│   │
│   └── server/
│       └── server.go                # Настройка HTTP сервера
│
├── pkg/
│   ├── logger/
│   │   └── logger.go                # Zap logger
│   └── response/
│       ├── success.go               # Стандартные успешные ответы
│       └── error.go                 # Стандартные ошибки
│
├── nginx/
│   ├── nginx.conf                   # Конфигурация Nginx
│   └── Dockerfile                   # Dockerfile для Nginx
│
├── api/
│   └── openapi/
│       └── gateway.yaml             # OpenAPI спецификация
│
├── configs/
│   └── .env.example
│
├── build/
│   └── Dockerfile                   # Dockerfile для Go приложения
│
├── .env
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### Описание ключевых файлов:

**`cmd/api/main.go`**
- Точка входа в приложение
- Загружает конфигурацию
- Инициализирует logger
- Создаёт HTTP клиенты для всех микросервисов
- Настраивает middleware
- Запускает HTTP сервер
- Обрабатывает graceful shutdown

**`internal/handler/auth.go`**
- Обработчики для эндпоинтов `/api/auth/register`, `/api/auth/login`
- Проксирует запросы в auth-service
- Обрабатывает ошибки и возвращает клиенту

**`internal/client/auth_client.go`**
- HTTP клиент для общения с auth-service
- Методы: `Register()`, `Login()`, `ValidateToken()`
- Использует retry механизм

**`internal/resilience/circuit_breaker.go`**
- Реализация Circuit Breaker паттерна
- Защищает от каскадных сбоев
- Отслеживает количество ошибок
- Временно блокирует запросы к упавшему сервису

**`nginx/nginx.conf`**
- Принимает весь входящий трафик на порт 80
- Проксирует запросы в api-gateway (Go приложение)
- Настраивает таймауты, логи, размер тела запроса

---

## 🔹 Auth Service (`services/auth-service/`)

**Назначение:** Аутентификация и авторизация пользователей. Регистрация, вход, генерация JWT токенов.

### Структура:

```
services/auth-service/
├── cmd/
│   └── api/
│       └── main.go                    # Точка входа
│
├── internal/
│   ├── app/
│   │   ├── app.go                    # Структура App
│   │   └── dependencies.go           # Dependency injection
│   │
│   ├── config/
│   │   └── config.go                 # Конфигурация (DB, JWT, Server)
│   │
│   ├── domain/
│   │   ├── entity/
│   │   │   └── user.go              # Сущность User (для БД)
│   │   │
│   │   ├── errors/
│   │   │   └── errors.go            # Доменные ошибки
│   │   │
│   │   └── models/
│   │       ├── request/
│   │       │   ├── register.go      # DTO для регистрации
│   │       │   ├── login.go         # DTO для входа
│   │       │   └── refresh.go       # DTO для refresh token
│   │       │
│   │       └── response/
│   │           ├── auth.go          # Ответ с токенами
│   │           ├── user.go          # Ответ с данными пользователя
│   │           └── error.go         # Стандартизированные ошибки
│   │
│   ├── repository/
│   │   ├── repository.go            # Интерфейс UserRepository
│   │   └── postgres/
│   │       └── user_repository.go   # Реализация для PostgreSQL
│   │
│   ├── service/
│   │   ├── service.go               # Интерфейс AuthService
│   │   └── auth_service.go          # Бизнес-логика аутентификации
│   │
│   ├── handler/
│   │   ├── middleware/
│   │   │   ├── auth.go             # Проверка JWT
│   │   │   ├── cors.go             # CORS
│   │   │   ├── logger.go           # Логирование
│   │   │   └── recovery.go         # Panic recovery
│   │   │
│   │   ├── auth.go                 # HTTP handlers для auth endpoints
│   │   ├── handler.go              # Главная структура Handler
│   │   └── routes.go               # Регистрация маршрутов
│   │
│   └── server/
│       └── server.go                # HTTP сервер
│
├── pkg/
│   ├── jwt/
│   │   ├── manager.go              # Генерация и валидация JWT
│   │   └── claims.go               # JWT Claims структура
│   │
│   ├── hash/
│   │   └── password.go             # Bcrypt для паролей
│   │
│   ├── logger/
│   │   └── logger.go               # Zap logger
│   │
│   ├── validator/
│   │   └── validator.go            # Валидация запросов
│   │
│   └── response/
│       ├── success.go              # Стандартные успешные ответы
│       └── error.go                # Стандартные ошибки
│
├── migrations/
│   ├── 001_create_users_table.up.sql
│   └── 001_create_users_table.down.sql
│
├── api/
│   └── openapi/
│       └── auth.yaml               # OpenAPI спецификация
│
├── configs/
│   └── .env.example
│
├── scripts/
│   ├── migrate-up.sh              # Применить миграции
│   └── migrate-down.sh            # Откатить миграции
│
├── build/
│   └── Dockerfile
│
├── .env
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### Описание ключевых файлов:

**`internal/domain/entity/user.go`**
- Сущность User для БД
- Поля: id, email, hashed_password, role, is_active, is_verified, created_at, updated_at
- Теги для ORM: `db:"field_name"`

**`internal/repository/postgres/user_repository.go`**
- Методы: `Create()`, `GetByEmail()`, `GetByID()`, `Exists()`, `UpdateLastLogin()`
- Работа с PostgreSQL через database/sql
- Обработка ошибок уникальности, NotFound

**`internal/service/auth_service.go`**
- `Register()` — хеширует пароль, создаёт пользователя
- `Login()` — проверяет пароль, генерирует JWT токены
- `ValidateToken()` — проверяет валидность токена
- `RefreshToken()` — обновляет access token

**`internal/handler/auth.go`**
- `POST /register` — регистрация
- `POST /login` — вход
- `POST /refresh` — обновление токена
- Валидация входных данных, обработка ошибок

**`pkg/jwt/manager.go`**
- `Generate()` — создание JWT с claims (user_id, role, exp)
- `Verify()` — проверка подписи и срока действия
- Использует HMAC-SHA256

**`migrations/001_create_users_table.up.sql`**
- CREATE TABLE users
- Индексы на email, is_active
- Constraints: UNIQUE(email), CHECK(role IN ('student', 'employer'))

---

## 🔹 User Service (`services/user-service/`)

**Назначение:** Управление профилями пользователей (студенты и работодатели).

### Структура:

```
services/user-service/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── dependencies.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── admin.go            # Сущность Admin
│   │   │   ├── student.go          # Сущность Student
│   │   │   └── employer.go         # Сущность Employer
│   │   │
│   │   ├── errors/
│   │   │   └── errors.go
│   │   │
│   │   └── models/
│   │       ├── request/
│   │       │   ├── create_admin.go
│   │       │   ├── update_admin.go
│   │       │   ├── create_student.go
│   │       │   ├── update_student.go
│   │       │   ├── create_employer.go
│   │       │   └── update_employer.go
│   │       │
│   │       └── response/
│   │           ├── admin.go
│   │           ├── student.go
│   │           └── employer.go
│   │
│   ├── repository/
│   │   ├── repository.go           # Интерфейсы AdminRepository, StudentRepository, EmployerRepository
│   │   └── postgres/
│   │       ├── admin_repository.go
│   │       ├── student_repository.go
│   │       └── employer_repository.go
│   │
│   ├── service/
│   │   ├── service.go              # Интерфейсы сервисов
│   │   ├── admin_service.go        # Логика работы с админами
│   │   ├── student_service.go      # Логика работы со студентами
│   │   └── employer_service.go     # Логика работы с работодателями
│   │
│   ├── handler/
│   │   ├── middleware/
│   │   │   ├── auth.go            # Проверка JWT через auth-service
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   │
│   │   ├── student.go             # Endpoints для студентов
│   │   ├── employer.go            # Endpoints для работодателей
│   │   ├── admin.go               # Endpoints для админов
│   │   ├── handler.go
│   │   └── routes.go
│   │
│   ├── cache/                      # Redis кэш
│   │   ├── admin_cache.go       # Кэширование профилей админов
│   │   ├── student_cache.go       # Кэширование профилей студентов
│   │   └── employer_cache.go      # Кэширование профилей работодателей
│   │
│   └── server/
│       └── server.go
│
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   ├── response/
│   │   ├── success.go
│   │   └── error.go
│   └── storage/                    # Для загрузки фото профилей
│       └── file_storage.go
│
├── migrations/
│   ├── 001_create_students_table.up.sql
│   ├── 001_create_students_table.down.sql
│   ├── 002_create_employers_table.up.sql
│   └── 002_create_employers_table.down.sql
│
├── api/
│   └── openapi/
│       └── user.yaml
│
├── configs/
│   └── .env.example
│
├── scripts/
│   ├── migrate-up.sh
│   └── migrate-down.sh
│
├── build/
│   └── Dockerfile
│
├── .env
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### Описание ключевых файлов:

**`internal/domain/entity/student.go`**
- Поля: id, user_id, full_name, university, major, course, skills, photo_url, created_at, updated_at
- Связь с users таблицей через user_id (FK)

**`internal/domain/entity/employer.go`**
- Поля: id, user_id, company_name, industry, description, website, logo_url, created_at, updated_at

**`internal/cache/student_cache.go`**
- `Set()` — сохранить профиль в Redis (TTL 15 минут)
- `Get()` — получить из кэша
- `Invalidate()` — удалить при обновлении профиля

**`internal/handler/student.go`**
- `POST /students` — создать профиль студента
- `GET /students/:id` — получить профиль (сначала из кэша)
- `PUT /students/:id` — обновить профиль
- `DELETE /students/:id` — удалить профиль

**`pkg/storage/file_storage.go`**
- Загрузка фото профиля
- Валидация типа файла (JPEG, PNG)
- Генерация уникального имени
- Сохранение на диск или S3

---

## 🔹 Vacancy Service (`services/vacancy-service/`)

**Назначение:** Управление вакансиями для практики.

### Структура:

```
services/vacancy-service/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── dependencies.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── entity/
│   │   │   └── vacancy.go          # Сущность Vacancy
│   │   │
│   │   ├── errors/
│   │   │   └── errors.go
│   │   │
│   │   └── models/
│   │       ├── request/
│   │       │   ├── create_vacancy.go
│   │       │   ├── update_vacancy.go
│   │       │   └── search_vacancy.go
│   │       │
│   │       └── response/
│   │           ├── vacancy.go
│   │           └── vacancy_list.go
│   │
│   ├── repository/
│   │   ├── repository.go
│   │   └── postgres/
│   │       └── vacancy_repository.go
│   │
│   ├── service/
│   │   ├── service.go
│   │   └── vacancy_service.go
│   │
│   ├── handler/
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   │
│   │   ├── vacancy.go
│   │   ├── handler.go
│   │   └── routes.go
│   │
│   ├── cache/
│   │   └── vacancy_cache.go        # Кэш популярных вакансий
│   │
│   └── server/
│       └── server.go
│
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   └── response/
│       ├── success.go
│       └── error.go
│
├── migrations/
│   ├── 001_create_vacancies_table.up.sql
│   └── 001_create_vacancies_table.down.sql
│
├── api/
│   └── openapi/
│       └── vacancy.yaml
│
├── configs/
│   └── .env.example
│
├── scripts/
│   ├── migrate-up.sh
│   └── migrate-down.sh
│
├── build/
│   └── Dockerfile
│
├── .env
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### Описание ключевых файлов:

**`internal/domain/entity/vacancy.go`**
- Поля: id, employer_id, title, description, requirements, location, salary_from, salary_to, employment_type, status, published_at, expires_at, created_at, updated_at

**`internal/repository/postgres/vacancy_repository.go`**
- `Create()`, `GetByID()`, `GetAll()`, `Search()`, `Update()`, `Delete()`
- `GetByEmployerID()` — вакансии конкретного работодателя
- Пагинация, фильтры, сортировка

**`internal/cache/vacancy_cache.go`**
- Кэширование списка популярных вакансий (TTL 10 минут)
- Инвалидация при публикации новой вакансии

**`internal/handler/vacancy.go`**
- `POST /vacancies` — создать вакансию (только employer)
- `GET /vacancies` — список вакансий с фильтрами
- `GET /vacancies/:id` — получить вакансию
- `PUT /vacancies/:id` — обновить (только автор)
- `DELETE /vacancies/:id` — удалить (только автор)

---

## 🔹 Resume Service (`services/resume-service/`)

**Назначение:** Управление резюме студентов.

### Структура:

```
services/resume-service/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── dependencies.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── entity/
│   │   │   └── resume.go           # Сущность Resume
│   │   │
│   │   ├── errors/
│   │   │   └── errors.go
│   │   │
│   │   └── models/
│   │       ├── request/
│   │       │   ├── create_resume.go
│   │       │   └── update_resume.go
│   │       │
│   │       └── response/
│   │           ├── resume.go
│   │           └── resume_list.go
│   │
│   ├── repository/
│   │   ├── repository.go
│   │   └── postgres/
│   │       └── resume_repository.go
│   │
│   ├── service/
│   │   ├── service.go
│   │   └── resume_service.go
│   │
│   ├── handler/
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   │
│   │   ├── resume.go
│   │   ├── handler.go
│   │   └── routes.go
│   │
│   ├── cache/
│   │   └── resume_cache.go
│   │
│   └── server/
│       └── server.go
│
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   └── response/
│       ├── success.go
│       └── error.go
│
├── migrations/
│   ├── 001_create_resumes_table.up.sql
│   └── 001_create_resumes_table.down.sql
│
├── api/
│   └── openapi/
│       └── resume.yaml
│
├── configs/
│   └── .env.example
│
├── scripts/
│   ├── migrate-up.sh
│   └── migrate-down.sh
│
├── build/
│   └── Dockerfile
│
├── .env
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### Описание ключевых файлов:

**`internal/domain/entity/resume.go`**
- Поля: id, student_id, title, summary, experience, education, skills, languages, achievements, status, published_at, created_at, updated_at

**`internal/handler/resume.go`**
- `POST /resumes` — создать резюме (только student)
- `GET /resumes` — список резюме
- `GET /resumes/:id` — получить резюме
- `PUT /resumes/:id` — обновить (только автор)
- `DELETE /resumes/:id` — удалить (только автор)

---

## 🔹 Response Service (`services/response-service/`)

**Назначение:** Управление откликами студентов на вакансии. Отправка событий в Kafka.

### Структура:

```
services/response-service/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── dependencies.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── entity/
│   │   │   └── response.go         # Сущность Response (отклик)
│   │   │
│   │   ├── errors/
│   │   │   └── errors.go
│   │   │
│   │   └── models/
│   │       ├── request/
│   │       │   ├── create_response.go
│   │       │   └── update_status.go
│   │       │
│   │       └── response/
│   │           ├── response.go
│   │           └── response_list.go
│   │
│   ├── repository/
│   │   ├── repository.go
│   │   └── postgres/
│   │       └── response_repository.go
│   │
│   ├── service/
│   │   ├── service.go
│   │   └── response_service.go
│   │
│   ├── handler/
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── recovery.go
│   │   │
│   │   ├── response.go
│   │   ├── handler.go
│   │   └── routes.go
│   │
│   ├── events/                      # Kafka producer
│   │   ├── producer.go             # Отправка событий в Kafka
│   │   └── events.go               # Типы событий
│   │
│   └── server/
│       └── server.go
│
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── validator/
│   │   └── validator.go
│   └── response/
│       ├── success.go
│       └── error.go
│
├── migrations/
│   ├── 001_create_responses_table.up.sql
│   └── 001_create_responses_table.down.sql
│
├── api/
│   └── openapi/
│       └── response.yaml
│
├── configs/
│   └── .env.example
│
├── scripts/
│   ├── migrate-up.sh
│   └── migrate-down.sh
│
├── build/
│   └── Dockerfile
│
├── .env
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### Описание ключевых файлов:

**`internal/domain/entity/response.go`**
- Поля: id, vacancy_id, student_id, resume_id, cover_letter, status, created_at, updated_at
- Status: pending, viewed, accepted, rejected

**`internal/events/producer.go`**
- `PublishResponseCreated()` — отправить событие "отклик создан" в Kafka
- `PublishStatusChanged()` — отправить "статус изменён"
- Topic: `response.created`, `response.status_changed`

**`internal/service/response_service.go`**
- `Create()` — создать отклик → отправить событие в Kafka
- `UpdateStatus()` — обновить статус → отправить событие

---

## 🔹 Notification Service (`services/notification-service/`)

**Назначение:** Отправка email и push уведомлений. Слушает события из Kafka.

### Структура:

```
services/notification-service/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── dependencies.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── consumer/                    # Kafka consumer
│   │   ├── consumer.go             # Основной consumer
│   │   ├── response_consumer.go    # Обработчик событий откликов
│   │   └── handler.go              # Роутинг событий к обработчикам
│   │
│   ├── service/
│   │   ├── email_service.go        # Отправка email (SMTP)
│   │   └── push_service.go         # Push уведомления (опционально)
│   │
│   └── template/                    # Email шаблоны
│       ├── response_created.html   # Шаблон "новый отклик"
│       └── status_changed.html     # Шаблон "статус изменён"
│
├── pkg/
│   └── logger/
│       └── logger.go
│
├── configs/
│   └── .env.example
│
├── build/
│   └── Dockerfile
│
├── .env
├── Makefile
├── README.md
├── go.mod
└── go.sum
```

### Описание ключевых файлов:

**`internal/consumer/response_consumer.go`**
- Подписывается на топики: `response.created`, `response.status_changed`
- При получении события вызывает email_service

**`internal/service/email_service.go`**
- `SendResponseCreatedEmail()` — отправить письмо работодателю о новом отклике
- `SendStatusChangedEmail()` — отправить письмо студенту об изменении статуса
- Использует SMTP, HTML шаблоны

**`internal/template/response_created.html`**
- HTML шаблон письма с подстановкой данных (имя студента, название вакансии)

---

## 🔹 Shared (`shared/`)

**Назначение:** Общий переиспользуемый код для всех сервисов.

### Структура:

```
shared/
├── pkg/
│   ├── logger/
│   │   └── logger.go               # Общий Zap logger
│   │
│   ├── validator/
│   │   └── validator.go            # Общая валидация (go-playground/validator)
│   │
│   ├── errors/
│   │   ├── errors.go               # Типы ошибок: NotFound, Unauthorized, ValidationError
│   │   └── codes.go                # HTTP коды для ошибок
│   │
│   ├── response/
│   │   ├── success.go              # JSON(200, data)
│   │   ├── error.go                # JSON(code, error)
│   │   └── pagination.go           # Pagination metadata
│   │
│   └── middleware/
│       ├── request_id.go           # Добавление X-Request-ID
│
└── README.md                        # Как использовать shared библиотеки
```

### Описание:

**`pkg/logger/logger.go`**
- Создаёт Zap logger с настройками: JSON формат, уровни, структурированные поля
- Импортируется всеми сервисами

**`pkg/errors/errors.go`**
- Определяет стандартные ошибки: `ErrNotFound`, `ErrUnauthorized`, `ErrValidation`, `ErrInternal`
- Методы `ToHTTPCode()`, `ToJSON()`

---

## 🔹 Deployments (`deployments/`)

**Назначение:** Конфигурации для запуска всей системы через Docker Compose.

### Структура:

```
deployments/
├── docker-compose.yml              # Production compose
├── docker-compose.dev.yml          # Development compose (с hot-reload)
├── .env.example                    # Пример переменных окружения

```

### Описание:

**`docker-compose.yml`**
- Поднимает все сервисы: nginx, api-gateway, auth, user, vacancy, resume, response, notification
- Инфраструктуру: postgres, redis, kafka, zookeeper
- Настраивает networks, volumes, healthchecks
- **Expose порт 80 (только Nginx)**

**`docker-compose.dev.yml`**
- Extends docker-compose.yml
- Добавляет volumes для hot-reload
- Открывает доп. порты для отладки

---

## 🔹 Scripts (`scripts/`)

**Назначение:** Скрипты для автоматизации задач.

### Структура:

```
scripts/
├── init-databases.sh               # Создание БД для всех сервисов
├── migrate-all.sh                  # Запуск миграций для всех сервисов
├── start-dev.sh                    # Запуск в dev режиме
├── stop-all.sh                     # Остановка всех контейнеров
├── test-all.sh                     # Запуск тестов всех сервисов
├── coverage-check.sh               # Проверка покрытия >= 30%
└── seed-data.sh                    # Заполнение БД тестовыми данными
```

### Описание:

**`migrate-all.sh`**
```bash
#!/bin/bash
# Применяет миграции для всех сервисов
for service in auth user vacancy resume response; do
  echo "Migrating $service-service..."
  cd services/$service-service
  make migrate-up
  cd ../..
done
```

**`test-all.sh`**
```bash
#!/bin/bash
# Запускает тесты всех сервисов
for service in services/*/; do
  echo "Testing $service..."
  cd $service
  go test -v -cover ./...
  cd ../..
done
```

**`coverage-check.sh`**
```bash
#!/bin/bash
# Проверяет покрытие >= 30%
for service in auth user vacancy resume response; do
  cd services/$service-service
  coverage=$(go test -cover ./... | grep 'total' | awk '{print $3}' | sed 's/%//')
  if (( $(echo "$coverage < 30" | bc -l) )); then
    echo "ERROR: $service coverage is $coverage% (< 30%)"
    exit 1
  fi
  cd ../..
done
```

---

## 🔹 Docs (`docs/`)

**Назначение:** Документация проекта.

### Структура:

```
docs/
├── README.md                       # Оглавление документации
├── architecture.md                 # Архитектура системы
├── deployment.md                   # Гайд по деплою
├── development.md                  # Гайд для разработчиков
│
├── api/                            # API документация
│   ├── api-gateway.md
│   ├── auth-service.md
│   ├── user-service.md
│   ├── vacancy-service.md
│   ├── resume-service.md
│   ├── response-service.md
│   └── notification-service.md

```

### Описание:

**`architecture.md`**
- Описание микросервисной архитектуры
- Диаграммы взаимодействия
- Паттерны: Circuit Breaker, Event-Driven
- Обоснование технологий

**`development.md`**
- Как запустить проект локально
- Структура кода каждого сервиса
- Best practices
- Как добавить новый эндпоинт

---

## 🔹 CI/CD (`.gitlab-ci.yml`)

**Назначение:** Автоматизация сборки, тестов, проверки покрытия.

### Структура:

```yaml
stages:
  - test
  - build
  - deploy

# Тесты для каждого сервиса
test:auth-service:
  stage: test
  script:
    - cd services/auth-service
    - go test -v -cover -coverprofile=coverage.out ./...
    - coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    - if (( $(echo "$coverage < 30" | bc -l) )); then exit 1; fi

test:user-service:
  stage: test
  script:
    - cd services/user-service
    - go test -v -cover -coverprofile=coverage.out ./...
    - coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    - if (( $(echo "$coverage < 30" | bc -l) )); then exit 1; fi

# ... аналогично для других сервисов

# Сборка Docker образов
build:
  stage: build
  script:
    - docker-compose build

# Деплой (manual)
deploy:
  stage: deploy
  script:
    - docker-compose up -d
  when: manual
```

---

## 🔹 Makefile (корневой)

**Назначение:** Упрощение команд для разработки.

```makefile
.PHONY: help build up down logs test migrate clean

help:
	@echo "Available commands:"
	@echo "  make build       - Build all services"
	@echo "  make up          - Start all services"
	@echo "  make down        - Stop all services"
	@echo "  make logs        - Show logs"
	@echo "  make test        - Run tests for all services"
	@echo "  make migrate     - Run migrations for all services"
	@echo "  make clean       - Remove all containers and volumes"

build:
	docker-compose -f deployments/docker-compose.yml build

up:
	docker-compose -f deployments/docker-compose.yml up -d

down:
	docker-compose -f deployments/docker-compose.yml down

logs:
	docker-compose -f deployments/docker-compose.yml logs -f

test:
	./scripts/test-all.sh

migrate:
	./scripts/migrate-all.sh

clean:
	docker-compose -f deployments/docker-compose.yml down -v
	docker system prune -f
```

---

## 📊 Итоговая структура

```
praktikujsya-microservices/
│
├── services/                          # 7 микросервисов (одинаковая структура)
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── vacancy-service/
│   ├── resume-service/
│   ├── response-service/
│   └── notification-service/
│
├── shared/                            # Общий код
│   └── pkg/
│
├── deployments/                       # Docker Compose + K8s
│   ├── docker-compose.yml
│   └── docker-compose.dev.yml
│
├── scripts/                           # Bash скрипты
│   ├── migrate-all.sh
│   ├── test-all.sh
│   └── ...
│
├── docs/                              # Документация
│   ├── architecture.md
│   ├── api/
│   └── diagrams/
│
├── .gitlab-ci.yml                     # CI/CD
├── Makefile                           # Команды
├── README.md
└── .gitignore
```

---

## ✅ Соответствие требованиям задания

| Требование | Реализация |
|-----------|-----------|
| Микросервисная архитектура | ✅ 7 сервисов с чёткими границами |
| Docker + Compose | ✅ Каждый сервис в контейнере |
| API Gateway + Nginx | ✅ Единая точка входа, порт 80 |
| Отказоустойчивость | ✅ Circuit Breaker, Retry, Fallback |
| Redis кэширование | ✅ В user, vacancy, resume сервисах |
| PostgreSQL | ✅ Database per service |
| Kafka (опционально) | ✅ Response → Notification через Kafka |
| Тесты >= 30% | ✅ Юнит + интеграционные тесты |
| GitLab CI/CD | ✅ Автоматические тесты и проверка покрытия |
| Документация | ✅ OpenAPI + Markdown docs |
| Логирование | ✅ Zap structured logging |

---

## 🚀 Быстрый старт

```bash
# 1. Клонировать проект
git clone <repo-url>
cd backend

# 2. Скопировать .env
cp deployments/.env.example deployments/.env

# 3. Запустить все сервисы
make up

# 4. Применить миграции
make migrate

# 5. Открыть в браузере
http://localhost
```

---
=======
# Backend сервиса "Практикуйся" - Документация

## 📁 Обзор структуры проекта

Этот проект представляет собой бэкенд-часть платформы для поиска практики студентами. Проект написан на Go и следует принципам чистой архитектуры.

## 🗂️ Детальное описание структуры файлов

### **cmd/api/main.go** - **Точка входа в приложение**
- **Назначение:** Главный файл, который запускает весь сервер
- **Что делает:**
  - Загружает конфигурацию из переменных окружения
  - Инициализирует логгер
  - Подключается к базе данных PostgreSQL
  - Настраивает все зависимости (репозитории, сервисы, хендлеры)
  - Регистрирует все маршруты API
  - Запускает HTTP-сервер на указанном порту
- **Важно:** Изменяйте этот файл только для добавления новой инициализации

### **internal/config/config.go** - **Конфигурация приложения**
- **Назначение:** Хранение и загрузка всех настроек приложения
- **Что содержит:**
  - Параметры подключения к БД (хост, порт, имя БД, логин, пароль)
  - Настройки JWT (секретный ключ, время жизни токена)
  - Настройки почтового сервиса (для верификации email)
  - Настройки файлового хранилища (пути для загрузки фото)
  - Порт сервера и уровень логирования
- **Пример использования:** `cfg.DB.Host`, `cfg.JWT.SecretKey`

### **internal/domain/entities/** - **Сущности базы данных**
- **Назначение:** Структуры, которые точно соответствуют таблицам в БД
- **Файлы:**
  - `user.go` - пользователь системы (студент/работодатель)
  - `student.go` - расширенная информация о студенте
  - `employer.go` - информация о компании-работодателе
  - `vacancy.go` - вакансии для практики
  - `resume.go` - резюме студентов
  - `response.go` - отклики студентов на вакансии
  - `report.go` - отчеты в техподдержку
- **Важно:** Эти структуры используются ТОЛЬКО для работы с БД

### **internal/domain/models/** - **Модели для API**
- **Назначение:** Структуры для входящих/исходящих данных API
- **request/** - **Данные, которые приходят от клиента:**
  - `auth.go` - запросы на регистрацию/вход
  - `student.go` - создание/обновление профиля студента
  - `employer.go` - создание/обновление профиля работодателя
  - `vacancy.go` - создание/редактирование вакансии
  - `resume.go` - создание/редактирование резюме
  - `response.go` - создание отклика на вакансию
- **response/** - **Данные, которые отправляются клиенту:**
  - Аналогичные файлы с структурами ответов
  - `auth.go` - содержит JWT токен после успешного входа

### **internal/repository/** - **Слой доступа к данным**
- **Назначение:** Работа с базой данных, выполнение SQL-запросов
- **interfaces/** - **Интерфейсы репозиториев:**
  - Определяют методы, которые ДОЛЖНЫ быть реализованы
  - Например: `CreateUser`, `GetUserByEmail`, `UpdateUser`
  - Позволяют легко менять реализацию (PostgreSQL → MySQL)
- **postgres/** - **Реализация для PostgreSQL:**
  - `user_repository.go` - работа с таблицей users
  - `student_repository.go` - работа с таблицей students
  - `employer_repository.go` - работа с таблицей employers
  - `vacancy_repository.go` - работа с таблицей vacancies
  - `resume_repository.go` - работа с таблицей resumes
  - `response_repository.go` - работа с таблицей responses
- **repository.go** - **Агрегатор всех репозиториев:**
  - Собирает все репозитории в одну структуру
  - Упрощает передачу зависимостей

### **internal/service/** - **Бизнес-логика приложения**
- **Назначение:** Вся логика работы приложения, проверки, бизнес-правила
- **interfaces/** - **Интерфейсы сервисов:**
  - Определяют контракты для бизнес-логики
  - Например: `RegisterStudent`, `CreateVacancy`, `ApplyToVacancy`
- **impl/** - **Реализация бизнес-логики:**
  - `auth_service.go` - регистрация, вход, верификация email
  - `student_service.go` - управление профилем студента
  - `employer_service.go` - управление профилем работодателя
  - `vacancy_service.go` - создание/публикация вакансий
  - `resume_service.go` - создание/публикация резюме
  - `response_service.go` - обработка откликов
- **service.go** - **Агрегатор всех сервисов**
  - Собирает все сервисы в одну структуру

### **internal/handler/** - **HTTP обработчики (контроллеры)**
- **Назначение:** Прием HTTP запросов и возврат ответов
- **middleware/** - **Промежуточное ПО:**
  - `auth.go` - проверка JWT токена, извлечение ID пользователя
  - `cors.go` - настройка CORS для фронтенда
  - `logger.go` - логирование всех входящих запросов
- **{entity}/** - **Обработчики по сущностям:**
  - `auth/` - регистрация и аутентификация
  - `student/` - работа с профилем студента
  - `employer/` - работа с профилем работодателя
  - `vacancy/` - управление вакансиями
  - `resume/` - управление резюме
  - `response/` - обработка откликов
  - `report/` - техподдержка
- **В каждом каталоге:**
  - `handler.go` - обработчики конкретных эндпоинтов
  - `routes.go` - регистрация маршрутов для этой сущности

### **internal/pkg/** - **Вспомогательные утилиты**
- **Назначение:** Переиспользуемый код, не связанный с бизнес-логикой
- **database/postgres.go** - **Подключение к БД:**
  - Создает пул соединений с PostgreSQL
  - Настраивает параметры подключения
- **logger/logger.go** - **Логирование:**
  - Настройка формата логов
  - Уровни логирования (debug, info, error)
- **jwt/jwt.go** - **Работа с JWT токенами:**
  - Генерация токенов при входе
  - Валидация токенов в middleware
  - Извлечение данных из токена (user_id, role)
- **email/email.go** - **Отправка email:**
  - Отправка писем для верификации email
  - Шаблоны писем
- **storage/file_storage.go** - **Работа с файлами:**
  - Загрузка фото профилей
  - Валидация типов и размеров файлов
  - Генерация уникальных имен файлов
- **validator/validator.go** - **Валидация данных:**
  - Проверка email, паролей, телефонных номеров
  - Валидация входных данных от пользователей

### **migrations/** - **Миграции базы данных**
- **Назначение:** Создание и обновление структуры БД
- `001_init_schema.sql` - **Первая миграция:**
  - Создание всех таблиц из схемы
  - Создание индексов для ускорения поиска
  - Добавление тестовых данных (опционально)

### **docker/** - **Контейнеризация**
- **Назначение:** Запуск приложения в Docker
- `Dockerfile` - **Сборка образа приложения:**
  - Многоэтапная сборка для уменьшения размера
  - Копирование бинарного файла в минимальный образ
- `docker-compose.yml` - **Оркестрация:**
  - Запуск PostgreSQL вместе с приложением
  - Настройка сети между контейнерами
  - Монтирование томов для данных

## 🔗 Взаимодействие между слоями

```
HTTP Request → Handler → Service → Repository → Database
                                   ↑
HTTP Response ← Handler ← Service ← Repository ← Database
```

1. **Handler** получает HTTP запрос, валидирует входные данные
2. **Handler** вызывает соответствующий метод **Service**
3. **Service** выполняет бизнес-логику, вызывает **Repository**
4. **Repository** выполняет SQL запрос к базе данных
5. Данные возвращаются по цепочке обратно к **Handler**
6. **Handler** формирует и отправляет HTTP ответ

## 🚀 Быстрый старт

### 1. Установка зависимостей
```bash
cd backend
go mod download
```

### 2. Настройка окружения
```bash
cp .env.example .env
# отредактируйте .env файл
```

### 3. Запуск миграций
```bash
make migrate-up
```

### 4. Запуск приложения
```bash
go run cmd/api/main.go
```

### 5. Запуск в Docker
```bash
docker-compose up --build
```

## 📝 Рекомендации по разработке

### Добавление нового эндпоинта:
1. Добавить структуру в `domain/models/request/`
2. Добавить структуру в `domain/models/response/`
3. Реализовать метод в `repository/postgres/`
4. Реализовать бизнес-логику в `service/impl/`
5. Создать обработчик в `handler/{entity}/`
6. Зарегистрировать маршрут в `handler/{entity}/routes.go`

### Добавление новой сущности:
1. Создать файл в `domain/entities/`
2. Создать интерфейс в `repository/interfaces/`
3. Реализовать репозиторий в `repository/postgres/`
4. Создать интерфейс в `service/interfaces/`
5. Реализовать сервис в `service/impl/`
6. Создать обработчики в `handler/{entity}/`

## 🔧 Технологический стек

- **Язык:** Go 1.21+
- **Фреймворк:** Gin 
- **База данных:** PostgreSQL 14+
- **Миграции:** golang-migrate
- **Аутентификация:** JWT
- **Документация:** Swagger/OpenAPI
- **Контейнеризация:** Docker, Docker Compose
>>>>>>> 35da1feff9e0cb7f6cd58e44cccdb47f150ee12d
