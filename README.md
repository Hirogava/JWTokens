# Сервис аутентификации с JWT токенами

Сервис предоставляет REST API для аутентификации с использованием JWT токенов. Реализована система с Access и Refresh токенами, включая механизм ротации Refresh токенов.

## Технологии

- Go
- JWT (HS512 алгоритм)
- PostgreSQL
- bcrypt для хеширования Refresh токенов

## Требования

- Go 1.21 или выше
- PostgreSQL 14 или выше

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/Hirogava/JWTokens.git
cd JWTokens
```

2. Создайте файл `.env` на основе `.env_example`:
```bash
cp .env_example .env
```

3. Настройте переменные окружения в `.env`:
```
DB_CONNECTION_STRING=
SERVER_PORT=
```

4. Запустите миграции базы данных:
```bash
go run main.go
```

## API Endpoints

### 1. Получение пары токенов

**Endpoint:** `POST /token`

**Request Body:**
```json
{
    "id": "3f7e188c-060a-4662-b216-b476dbf1f321",
    "ip": "127.0.0.1",
    "email": "test@example.com"
}
```

**Response:**
```json
{
    "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "dGhpcyBpcyBhIHNlY3JldCB0b2tlbg=="
}
```

### 2. Обновление токенов

**Endpoint:** `POST /token/refresh`

**Request Body:**
```json
{
    "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "dGhpcyBpcyBhIHNlY3JldCB0b2tlbg==",
    "id": "3f7e188c-060a-4662-b216-b476dbf1f321"
}
```

**Response:**
```json
{
    "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "bmV3IHNlY3JldCB0b2tlbg=="
}
```

## Тестирование через Postman

### 1. Получение токенов

1. Создайте новый запрос в Postman:
   - Метод: POST
   - URL: `http://localhost:8080/token`
   - Headers:
     ```
     Content-Type: application/json
     ```
   - Body (raw JSON):
     ```json
     {
         "id": "3f7e188c-060a-4662-b216-b476dbf1f321",
         "ip": "127.0.0.1",
         "email": "test@example.com"
     }
     ```

2. Отправьте запрос и сохраните полученные токены

### 2. Обновление токенов

1. Создайте новый запрос в Postman:
   - Метод: POST
   - URL: `http://localhost:8080/token/refresh`
   - Headers:
     ```
     Content-Type: application/json
     ```
   - Body (raw JSON):
     ```json
     {
         "accessToken": "<ваш_access_token>",
         "refreshToken": "<ваш_refresh_token>",
         "id": "3f7e188c-060a-4662-b216-b476dbf1f321"
     }
     ```

2. Отправьте запрос и получите новую пару токенов

## Особенности реализации

1. **Access Token:**
   - JWT формат
   - Алгоритм подписи HS512
   - Срок действия 15 минут
   - Содержит информацию о пользователе и IP адресе

2. **Refresh Token:**
   - Произвольный формат
   - Передается в base64
   - Хранится в базе в виде bcrypt хеша
   - Может быть использован только один раз
   - Автоматически удаляется через 30 дней

3. **Безопасность:**
   - Проверка IP адреса при обновлении токенов
   - Ротация Refresh токенов
   - Защита от повторного использования токенов
   - Автоматическая очистка старых токенов

## Обработка ошибок

Сервис возвращает следующие HTTP статусы:

- 200 OK - успешное выполнение
- 400 Bad Request - неверный формат запроса
- 401 Unauthorized - неверные или истекшие токены
- 500 Internal Server Error - внутренняя ошибка сервера

## Логирование

Сервис логирует следующие события:
- Получение новых токенов
- Обновление токенов
- Попытки повторного использования Refresh токенов
- Изменение IP адреса
- Ошибки при работе с базой данных

## Тестовые данные

Для тестирования можно использовать следующие данные:
- GUID: `3f7e188c-060a-4662-b216-b476dbf1f321`
- IP: `127.0.0.1`
- Email: `test@example.com`
