package postgres

// SQL запросы для создания таблиц.
const (
	// SQL запрос для создания таблицы пользователей.
	createUsers = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		is_admin BOOlEAN NOT NULL
		);`
	// SQL запрос для создания таблицы актёров.
	createActors = `CREATE TABLE IF NOT EXISTS actors (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		gender TEXT NOT NULL,
		date_of_birth DATE NOT NULL
		);`
	// SQL запрос для создания таблицы фильмов.
	createMovies = `CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		release_date DATE NOT NULL,
		rating INTEGER NOT NULL
		);`
	// SQL запрос для создания таблицы отношений фильмов и актёров.
	createActorMovieRelations = `CREATE TABLE IF NOT EXISTS movie_actors (
		movie_id INTEGER NOT NULL,
		actor_id INTEGER NOT NULL,
		FOREIGN KEY (movie_id) REFERENCES movies(id),
    	FOREIGN KEY (actor_id) REFERENCES actors(id),
    	PRIMARY KEY (movie_id, actor_id)
		);`
)

// SQL запросы для удаления таблиц.
const (
	// SQL запрос для удаления таблицы пользователей.
	dropUsers = `DROP TABLE IF EXISTS users;`
	// SQL запрос для удаления таблицы фильмов.
	dropMovies = `DROP TABLE IF EXISTS movies;`
	// SQL запрос для удаления таблицы актёров.
	dropActors = `DROP TABLE IF EXISTS actors;`
	// SQL запрос для удаления таблицы отношений фильмов и актёров.
	dropMovieActors = `DROP TABLE IF EXISTS movie_actors;`
)

// SQL запросы для добавления данных в БД.
const (
	// SQL запрос для добавления пользователя по name, password, is_admin.
	addUser = `INSERT INTO users (name, password, is_admin) VALUES ($1, $2, $3) RETURNING id;`
	// SQL запрос для добавления актёра по name, gender, date_of_birth.
	addActor = `INSERT INTO actors (name, gender, date_of_birth) VALUES ($1, $2, $3) RETURNING id;`
	// SQL запрос для добавления фильма по name, description, release_date, rating.
	addMovie = `INSERT INTO movies (name, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id;`
	// SQL запрос для добавления актёра в фильм по movie_id, actor_id.
	addActorToMovie = `INSERT INTO movie_actors (movie_id, actor_id) VALUES ($1, $2);`
)

// SQL запросы для удаления данных.
const (
	// SQL запрос для удаления актёра по id.
	removeActor = `DELETE FROM actors WHERE id = $1; DELETE FROM movie_actors WHERE actor_id = $1;`
	// SQL запрос для удаления фильма по id.
	removeMovie = `DELETE FROM movies WHERE id = $1; DELETE FROM movie_actors WHERE movie_id = $1;`
	// SQL запрос для удаления фильма из работ актёров по movie_id.
	removeMovieFromActors = `DELETE FROM movie_actors WHERE movie_id = $1;`
)

// SQL запросы для обновления данных.
const (
	// SQL запрос для обновления фильма по id, name, description, release_date, rating.
	updateMovie = `UPDATE movies SET name = $2, description = $3, release_date = $4, rating = $5 WHERE id = $1;`
	// SQL запрос для обновления актёра по id, name, gender, date_of_birth.
	updateActor = `UPDATE actors SET name = $2, gender = $3, date_of_birth = $4 WHERE id = $1;`
)

// SQL запросы для получения данных.
const (
	// SQL запрос для получения актёра по id.
	getActor = `SELECT * FROM actors where id = $1;`
	// SQL запрос для получения актёров.
	getActors = `SELECT * from actors;`
	// SQL запрос для получения фильмов, в которых играл актёр, по actor_id.
	getActorMovies = `SELECT movie_id FROM movie_actors WHERE actor_id = $1;`
	// SQL запрос для получения фильма по id.
	getMovie = `SELECT * FROM movies WHERE id = $1`
	// SQL запрос для получения актёров, которые играли в фильме, по movie_id.
	getMovieActors = `SELECT actor_id FROM movie_actors WHERE movie_id = $1;`
	// SQL запрос для получения фильмов, отсортированных по рейтингу.
	getMoviesSortByRating = `SELECT * FROM movies ORDER BY rating DESC;`
	// SQL запрос для получения фильмов, отсортированных по дате релиза.
	getMoviesSortByReleaseDate = `SELECT * FROM movies ORDER BY release_date;`
	// SQL запрос для получения фильмов, отсортированных по названию.
	getMoviesSortByName = `SELECT * FROM movies ORDER BY name;`
	// SQL запрос для получения фильмов по фрагменту имени актёра.
	getMoviesByActor = `SELECT * FROM movies WHERE id IN (
		SELECT movie_id FROM movie_actors ma JOIN actors a ON ma.actor_id = a.id
		WHERE a.name LIKE '%$1%'
		);`
	// SQL запрос для получения фильмов по фрагменту названия.
	getMoviesByName = `SELECT * FROM movies WHERE name LIKE '%$1%';`
	// SQL запрос для получения статуса пользователя по name, password.
	checkUserRole = `SELECT is_admin FROM users WHERE name = $1 AND password = $2;`
)
