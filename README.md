# Go Banking API

## Описание проекта

Go Banking API — это RESTful API-сервер, реализующий основные функции банковского сервиса. Он предоставляет механизм регистрации и авторизации пользователей, создание и пополнение банковских счетов, переводы между пользователями, отправку email-уведомлений и защиту операций с помощью JWT-аутентификации.
## Основной функционал

* Регистрация пользователей с валидацией и проверкой уникальности email и username
* Аутентификация через JWT (срок действия токена — 24 часа)
* Создание банковских счетов, пополнение баланса
* Переводы средств между пользователями по username
* Email-уведомления о поступлении перевода через SMTP (например, Gmail)
* Защищённые маршруты с middleware для авторизации
* Логирование действий через logrus

## Технологии и библиотеки

* Язык: Go 1.20+
* Роутинг: gorilla/mux
* БД: PostgreSQL 17, lib/pq
* JWT: golang-jwt/jwt/v5
* Email: go-mail/mail
* Хеширование: bcrypt
* Логирование: logrus
* Управление переменными окружения: godotenv

## Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/your-username/go-banking-api.git
cd go-banking-api
```

### 2. Настройка переменных окружения

Создайте файл `.env` на основе `.env.example`:

```env
DB_URL=postgres://postgres:password@localhost:5432/banking?sslmode=disable
JWT_SECRET=your_jwt_secret
PORT=8080

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=youraddress@gmail.com
SMTP_PASS=your_app_password
```

### 3. Подготовка базы данных

Создайте БД `banking` в PostgreSQL:

```sql
CREATE DATABASE banking;
```

Затем выполните SQL-скрипт `init.sql`, находящийся в папке `migrations`, чтобы создать необходимые таблицы:

```bash
psql -U postgres -d banking -f migrations/init.sql
```

Это создаст таблицы `users`, `accounts` и другие необходимые для работы проекта.

### 4. Установка зависимостей

```bash
go mod tidy
```

### 5. Запуск приложения

```bash
go run cmd/main.go
```

Сервер будет запущен на порту, указанном в переменной `PORT` (по умолчанию 8080).

## Примеры работы

Ниже приведены примеры запросов и ответов, которые возвращает API в формате JSON:

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "pass123"
}'
```

**Ответ:**

```json
{
  "id": 1,
  "email": "test@example.com",
  "username": "testuser",
  "createdAt": "2025-05-12T12:34:56Z"
}
```

### Авторизация

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "pass123"
}'
```

**Ответ:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Создание счёта

```bash
curl -X POST http://localhost:8080/accounts \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json"
```

**Ответ:**

```json
{
  "id": 2,
  "balance": 0,
  "createdAt": "2025-05-12T13:00:00Z"
}
```

### Пополнение счёта

```bash
curl -X POST http://localhost:8080/accounts/topup \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 1,
    "amount": 1000
}'
```

**Ответ:**

```json
{
  "status": "ok"
}
```

### Перевод другому пользователю

```bash
curl -X POST http://localhost:8080/transfer/by-usernames \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "from_username": "testuser",
    "to_username": "recipient",
    "amount": 500
}'
```

**Ответ:**

```json
{
  "status": "ok"
}
```

После перевода пользователь `recipient@example.com` получит email-уведомление, если у него указан email в системе.

## Примеры ошибок

### Ошибка: пользователь уже существует

**Ответ:**

```json
{
  "error": "email или username уже используется"
}
```

### Ошибка: неверный email или пароль при логине

**Ответ:**

```json
{
  "error": "неверный email или пароль"
}
```

### Ошибка: попытка перевести себе

**Ответ:**

```json
{
  "error": "нельзя переводить самому себе"
}
```

### Ошибка: недостаточно средств

**Ответ:**

```json
{
  "error": "недостаточно средств на счёте"
}
```

### Ошибка: неавторизованный доступ к защищённому эндпоинту

**Ответ:**

```json
{
  "error": "Authorization header required"
}
```
