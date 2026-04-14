# Wishlist service
Сервис для создания вишлистов с подарками.

[Техническое задание](./docs/Go.md)

### Функционал:

- Регистрация и авторизация по почте

- Создание, редактирование, удаление вишлистов и позиций в нем (только для авторизованных пользователей)

- Возможность делиться вишлистом через публичный токен
- Возможность бронировать подарок в публичном вишлисте по токену и ID позиции.

#### Бизнес модели:

**Wishlist**: Название, Описание, Дата

**Gift**: Название, Описание, Ссылка, Приоритет

### Архитектура

Проект построен по принципам hexogonal architecture
```
cmd/ - точка входа в приложение
internal/ - внутренние пакеты
  adapter/ - адаптеры
     in/ - входящие запросы
       dto/ - data objects
       httpservice/ - http хэндлеры
     out/ - исходящие запросы
       repository/ - репозитории
  app/ 
    app.go - запуск приложения
    config.go - кофиг
  database/ - запуск миграций
  model/ - модели бизнес логики
  service/ - точка входа в бизнес логику
database/ - файлы миграций
tests/ - интеграционные тесты
```

## Запуск

Пример необходимых для запуска переменных окружения задан в [.env.example](.env.example)

Сервис по умолчанию использует порт **8080**, PostgreSQL - **5432**

Документация **swagger** будет доступна по адресу [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Для авторизации в swagger-ui нужно получить jwt-токен и ввести его в поле **Authorize** в виде ```Bearer your-token```

### Docker
Для запуска в Docker-контейнере:

`docker-compose up --build` - поднимает PostgreSQL в контейнере **db** и сервис в **app**

### Локально
Для локальной сборки:

```go run ./cmd/main.go```

Для запуска юнит тестов:

```go test ./...```

Для запуска интеграционных тестов:

```go test -tags=integration ./tests/...```

Интеграционные тесты требуют установленный Docker.

## Github actions

Для проекта настроен **CI pipeline** в GitHub Actions.

Workflow автоматически запускается при `push` и `pull request` в ветку `main`.

Стадии: **lint**, **unit-test**, **integration-test**, **build**

## Подробности реализации

### Стек технологий
- **Gin** - HTTP-фреймворк

- **pgx** - работа с PostgreSQL

- **goose** - миграции базы данных

- **testcontainers-go** - интеграционное тестирование с использованием PostgreSQL в контейнере

### Реализованные доп.требования и бонусы

- **Unit-тесты** для бизнес-логики

- **Интеграционные тесты** сценариев создания вишлиста и бронирования в нем позиции [./tests](./tests)

- **Swagger-документация** API с помощью swaggo

- **Graceful shutdown** через cancel контекста

- **CI pipeline** в Github Actions со стадиями линтера, тестов и сборки приложения
- **golangci-lint** для чистоты кода

### Принятые решения

- Подарок - зависимая от вишлиста сущность, работа с ним через API завязана через wishlist пути
- Бронирование подарка происходит через атомарный update, без использования транзакций, т.к. операция достаточно простая и атомарности SQL-запроса хватает
- Интеграционные тесты используют **testcontainers-go** для поднятия изолированной чистой базы данных, обеспечивая воспроизводимость тестов

## Основные эндпоинты

Все эндпоинты, кроме публичных, требуют JWT-токен в заголовке:

```http
Authorization: Bearer <token>
```

#### Авторизация
`POST /api/auth/register` - регистрация нового пользователя по email и паролю.

`POST /api/auth/login` - аутентификация пользователя и получение JWT-токена.

#### Вишлисты
`POST /api/wishlists` - создание нового вишлиста для авторизованного пользователя.

`GET /api/wishlists` - получение списка всех вишлистов текущего пользователя.

`GET /api/wishlists/details/{id}` - получение одного вишлиста пользователя по ID вместе с позициями.

`PUT /api/wishlists/details/{id}` - обновление вишлиста по ID.

`DELETE /api/wishlists/details/{id}` - удаление вишлиста по ID.

#### Подарки
`POST /api/wishlists/{wishlistId}/gifts` - добавление новой позиции в указанный вишлист.

`PUT /api/wishlists/{wishlistId}/gifts/{id}` - обновление позиции подарка в указанном вишлисте.

`DELETE /api/wishlists/{wishlistId}/gifts/{id}` - удаление позиции подарка из вишлиста.

#### Публичный доступ
`GET /api/public/wishlists/token/{token}` - получение публичного вишлиста по токену без авторизации.

`POST /api/public/wishlists/{token}/gifts/{id}` - бронирование подарка в публичном вишлисте по токену и ID позиции.

## Примеры запросов

Ниже приведены примеры запросов в формате **bash curl**.

#### Регистрация
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
  "email": "user@example.com",
  "password": "secret123"
}'
```
#### Логин
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
  "email": "user@example.com",
  "password": "secret123"
}'
```

Пример ответа:
```json
{
  "token": "your-jwt-token"
}
```

#### Создать wishlist
```bash
curl -X POST http://localhost:8080/api/wishlists \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
  "title": "Birthday",
  "description": "My birthday wishlist",
  "date": "2030-01-02T15:04:05Z"
}'
```

#### Добавить gift в wishlist
```bash
curl -X POST http://localhost:8080/api/wishlists/1/gifts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
  "name": "LEGO",
  "description": "Big set",
  "link": "https://example.com/gift",
  "priority": 5
}'
```

#### Получить публичный wishlist по токену
```bash
curl -X GET http://localhost:8080/api/public/wishlists/token/your-wishlist-token
```
#### Забронировать gift по публичному токену
```bash
curl -X POST http://localhost:8080/api/public/wishlists/your-wishlist-token/gifts/1
```

Пример ответа при конфликте бронирования:
```json
{
  "error": "already_booked",
  "message": "gift already booked"
}
```

Далее приведены примеры не основных CRUD-запросов.

#### Получить список своих wishlist
```bash
curl -X GET http://localhost:8080/api/wishlists \
  -H "Authorization: Bearer your-jwt-token"
```

Пример ответа:
```json
{
  "list": [
    {
      "id": 1,
      "token": "wishlist-public-token",
      "title": "Birthday",
      "description": "My birthday wishlist",
      "date": "2030-01-02T15:04:05Z"
    }
  ]
}
```
#### Получить wishlist по ID
```bash
curl -X GET http://localhost:8080/api/wishlists/details/1 \
  -H "Authorization: Bearer your-jwt-token"
```
#### Обновить wishlist
```bash
curl -X PUT http://localhost:8080/api/wishlists/details/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
  "title": "Updated birthday wishlist"
  }'
```
#### Удалить wishlist
```bash
curl -X DELETE http://localhost:8080/api/wishlists/details/1 \
  -H "Authorization: Bearer your-jwt-token"
```

#### Обновить gift
```bash
curl -X PUT http://localhost:8080/api/wishlists/1/gifts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-jwt-token" \
  -d '{
  "priority": 4
  }'
```
#### Удалить gift
```bash
curl -X DELETE http://localhost:8080/api/wishlists/1/gifts/1 \
  -H "Authorization: Bearer your-jwt-token"
```
