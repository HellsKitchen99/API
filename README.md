# apiWithDataBase

REST API на Go для управления задачами с использованием PostgreSQL.

## Описание

Этот проект реализует API для работы с задачами (CRUD), а также регистрацию и аутентификацию пользователей с помощью JWT. Данные хранятся в базе данных PostgreSQL.

## Запуск через Docker

1. Скопируйте `.env` с нужными переменными окружения (см. пример ниже).
2. Соберите и запустите контейнеры:
   ```sh
   make up
   ```
3. API будет доступен на `localhost:8352`.

## Основные эндпоинты

- `POST /register` — регистрация пользователя
- `POST /login` — вход пользователя (выдаёт JWT)
- `POST /refresh` — обновление access-токена
- `GET /tasks` — получить все задачи
- `POST /tasks` — добавить задачу (требует авторизации)
- `PATCH /tasks/{id}` — изменить задачу (требует авторизации)
- `DELETE /tasks/{id}` — удалить задачу (требует авторизации)
- `GET /task/{id}` — получить задачу по id
- `GET /tasks/overdue` — просроченные задачи
- `GET /me` — информация о текущем пользователе (требует авторизации)

## Пример .env

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=ваш_пароль
DB_NAME=apishka
secretKeyAccess=ваш_секрет_для_access
secretKeyRefresh=ваш_секрет_для_refresh
```

## Тесты

Для запуска тестов:
```sh
go test ./...
```

## Структура проекта

- [main.go](main.go) — точка входа
- [handlers/](handlers/) — HTTP-обработчики
- [database/](database/) — работа с БД
- [tasks/](tasks/tasks.go) — структуры данных
- [consts/](consts/envKey.go) — константы

---