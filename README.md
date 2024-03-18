### Выполнил Степанов Алексей Александрович

Для проекта я использовал Go 1.22.1 *(версия 1.22 необходима для работы приложения, так как в нёи используется новый функционал пакета net/http)* и БД PostgreSQL.

## Запуск:

Запуск dev-среды с помощью docker-compose:

```bash
docker-compose up
```
Среда, в которой происходит запуск, должна иметь переменные окружения:
- DB_HOST - хост базы данных.
- DB_PORT - порт для доступа к базе данных.
- DB_USER - логин пользователя БД.
- DB_PASSWORD - пароль пользователя БД.
- DB_NAME - ИМЯ пользователя БД.

Запуск с помощью go run:

```bash
go run ./cmd/api/main.go
# Флаги:
# -override_tables=true - запуск с автоматическим созданием таблиц в БД
# -addr=:8080 - выбор порта, с которым будет работать сервер
# -default_admin=true - запуск с существованием базового администратора (admin|admin).
```

## PostgreSQL Query для создания таблиц в БД вручную:

```sql
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	is_admin BOOlEAN NOT NULL
);
CREATE TABLE IF NOT EXISTS actors (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	gender TEXT NOT NULL,
	date_of_birth DATE NOT NULL
);
CREATE TABLE IF NOT EXISTS movies (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	release_date DATE NOT NULL,
	rating INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS movie_actors (
	movie_id INTEGER NOT NULL,
	actor_id INTEGER NOT NULL,
	FOREIGN KEY (movie_id) REFERENCES movies(id),
   	FOREIGN KEY (actor_id) REFERENCES actors(id),
   	PRIMARY KEY (movie_id, actor_id)
);
```

## UI Swagger доступен по адресу `/swagger`

Документация располагается в папке [docs](./docs/)

Для генерации документации использовалась утилита Swag:

```bash
# Установка Swag
go install github.com/swaggo/swag/cmd/swag@latest
# Генерация документации
swag init -d ./cmd/api,./internal/filmoteka -o ./docs
```

## Комментарий

Так как в задании было разрешено упростить процесс сопоставления ролей и пользователей (считать, что роли задаются вручную), мною было решено устроить это так:

Новых пользователей через Api может создавать только администратор. При создании указывается логин, пароль и роль.
Также, по умолчанию, пользователь admin admin является администратором. 