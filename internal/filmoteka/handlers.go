package filmoteka

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/models"

	"golang.org/x/crypto/bcrypt"
)

// TODO comments
// TODO wrap errors

// AddActor - обрабатывает http запрос на добавление актёра в фильмотеку.
//
// @Summary      Adds actor to the System.
// @Description  Add actor to the System and get it's ID. User should be an admin.
// @Tags         Actor
// @Accept       json
// @Produce      json
// @Param        actor body models.Actor true "Actor to be added"
// @Security BasicAuth
// @Success      200 {integer} int "ID of the added actor"
// @Failure      400 {string} string "Bad request"
// @Failure      403 {string} string "User not an admin"
// @Failure      500 {string} string "Internal server error"
// @Router       /actor [post]
func (app *App) AddActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to add a new actor")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if !isAdmin {
		app.infoLog.Println("user trying to add a new actor not an admin")
		http.Error(w, "user not an admin", http.StatusForbidden)
		return
	}

	var actor models.Actor
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&actor); err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := app.dbHandler.AddActor(actor.Name, actor.Gender, actor.DateOfBirth)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, id)
	app.infoLog.Printf("actor %d is added\n", id)
}

// UpdateActor - обрабатывает http запрос на обновление актёра в фильмотеке.
//
// @Summary      Updates actor in the System.
// @Description  Update actor in the System. User should be an admin.
// @Tags         Actor
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of the actor to be updated"
// @Param        actor body models.Actor true "Actor data to be updated"
// @Security BasicAuth
// @Success      200 {string} string "OK"
// @Failure      400 {string} string "Bad request"
// @Failure      403 {string} string "User not an admin"
// @Failure      500 {string} string "Internal server error"
// @Router       /actor/{id} [put]
func (app *App) UpdateActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to update an actor")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if !isAdmin {
		app.infoLog.Println("user trying to update an actor not an admin")
		http.Error(w, "user not an admin", http.StatusForbidden)
		return
	}

	// TODO обновлять только то, что пришло в запросе

	var actor models.Actor
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&actor); err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "id must be an integer", http.StatusBadRequest)
		return
	}
	err = app.dbHandler.UpdateActor(id, actor.Name, actor.Gender, actor.DateOfBirth)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.infoLog.Printf("actor %d is updated\n", id)
}

// DeleteActor - обрабатывает http запрос на удаление актёра из фильмотеки.
//
// @Summary      Deletes actor from the System.
// @Description  Delete actor from the System. User should be an admin.
// @Tags         Actor
// @Produce      json
// @Param        id path int true "ID of the actor to be deleted"
// @Security BasicAuth
// @Success      200 {string} string "OK"
// @Failure      400 {string} string "Bad request"
// @Failure      403 {string} string "User not an admin"
// @Failure      500 {string} string "Internal server error"
// @Router       /actor/{id} [delete]
func (app *App) DeleteActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to delete an actor")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if !isAdmin {
		app.infoLog.Println("user trying to delete an actor not an admin")
		http.Error(w, "user not an admin", http.StatusForbidden)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "id must be an integer", http.StatusBadRequest)
		return
	}
	err = app.dbHandler.DeleteActor(id)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.infoLog.Printf("actor %d is deleted\n", id)
}

// GetActor - обрабатывает http запрос на получение актёра из фильмотеки.
//
// @Summary      Get actor from the System.
// @Description  Get actor from the System.
// @Tags         Actor
// @Produce      json
// @Param        id path int true "ID of the actor to be getted"
// @Security BasicAuth
// @Success      200 {object} models.Actor
// @Failure      400 {string} string "Bad request"
// @Failure      500 {string} string "Internal server error"
// @Failure      403 {string} string "User does not exist"
// @Router       /actor/{id} [get]
func (app *App) GetActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get an actor")
	_, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "id must be an integer", http.StatusBadRequest)
		return
	}
	actor, err := app.dbHandler.GetActor(id)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, actor)
	app.infoLog.Printf("actor %d is getted\n", id)
}

// GetActors - обрабатывает http запрос на получение списка актёров из фильмотеки.
//
// @Summary      Get actors from the System.
// @Description  Get actors from the System.
// @Tags         Actor
// @Produce      json
// @Security BasicAuth
// @Success      200 {array} models.Actor
// @Failure      400 {string} string "Bad request"
// @Failure      500 {string} string "Internal server error"
// @Failure      403 {string} string "User does not exist"
// @Router       /actors [get]
func (app *App) GetActors(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of actors")
	_, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	actors, err := app.dbHandler.GetActors()
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, actors)
	app.infoLog.Println("list of actors is getted")
}

// AddMovie - обрабатывает http запрос на добавление фильма в фильмотеку.
//
// @Summary      Adds movie to the System.
// @Description  Add movie to the System and get it's ID. User should be an admin.
// @Tags         Movie
// @Accept       json
// @Produce      json
// @Param        movie body models.Movie true "Movie to be added"
// @Security BasicAuth
// @Success      200 {integer} int "ID of the added movie"
// @Failure      400 {string} string "Bad request"
// @Failure      403 {string} string "User not an admin"
// @Failure      500 {string} string "Internal server error"
// @Router       /movie [post]
func (app *App) AddMovie(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to add a new movie")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		app.infoLog.Println("user trying to add a new movie not an admin")
		http.Error(w, "user not an admin", http.StatusForbidden)
		return
	}

	var movie models.Movie
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&movie); err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := app.dbHandler.AddMovie(movie.Name, movie.Description, movie.ReleaseDate, movie.Rating, movie.Actors)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, id)
	app.infoLog.Printf("movie %d is added\n", id)
}

// DeleteMovie - обрабатывает http запрос на удаление фильма из фильмотеки.
//
// @Summary      Deletes movie from the System.
// @Description  Delete movie from the System. User should be an admin.
// @Tags         Movie
// @Produce      json
// @Param        id path int true "ID of the movie to be deleted"
// @Security BasicAuth
// @Success      200 {string} string "OK"
// @Failure      400 {string} string "Bad request"
// @Failure      403 {string} string "User not an admin"
// @Failure      500 {string} string "Internal server error"
// @Router       /movie/{id} [delete]
func (app *App) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to delete a movie")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		app.infoLog.Println("user trying to delete an movie not an admin")
		http.Error(w, "user not an admin", http.StatusForbidden)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "id must be an integer", http.StatusBadRequest)
		return
	}
	err = app.dbHandler.DeleteMovie(id)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.infoLog.Printf("movie %d is deleted\n", id)
}

// UpdateMovie - обрабатывает http запрос на обновление фильма в фильмотеке.
//
// @Summary      Updates movie in the System.
// @Description  Update movie in the System. User should be an admin.
// @Tags         Movie
// @Accept       json
// @Produce      json
// @Param        id path int true "ID of the movie to be updated"
// @Param        movie body models.Movie true "Movie data to be updated"
// @Security BasicAuth
// @Success      200 {string} string "OK"
// @Failure      400 {string} string "Bad request"
// @Failure      403 {string} string "User not an admin"
// @Failure      500 {string} string "Internal server error"
// @Router       /movie/{id} [put]
func (app *App) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to update a movie")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		app.infoLog.Println("user trying to update a movie not an admin")
		http.Error(w, "user not an admin", http.StatusForbidden)
		return
	}

	var movie models.Movie
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&movie); err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, "id must be an integer", http.StatusBadRequest)
		return
	}
	err = app.dbHandler.UpdateMovie(id, movie.Name, movie.Description, movie.ReleaseDate, movie.Rating, movie.Actors)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.infoLog.Printf("movie %d is updated\n", id)
}

// GetMovies - обрабатывает http запрос на получение списка фильмов из фильмотеки.
//
// @Summary      Get movies from the System.
// @Description  Get movies from the System.
// @Tags         Movie
// @Produce      json
// @Param        sort query string false "Sort movies by name, release date or rating"
// @Security BasicAuth
// @Success      200 {array} models.Movie
// @Failure      500 {string} string "Internal server error"
// @Failure      403 {string} string "User does not exist"
// @Router       /movies [get]
func (app *App) GetMovies(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of movies")
	_, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var sortBy int
	sorts, ok := r.URL.Query()["sort"]
	if !ok {
		sortBy = models.SortByRating
	} else {
		switch sorts[0] {
		case "name":
			sortBy = models.SortByName
		case "release":
			sortBy = models.SortByReleaseDate
		case "rating":
			fallthrough
		default:
			sortBy = models.SortByRating
		}
	}
	movies, err := app.dbHandler.GetMovies(sortBy)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, movies)
	app.infoLog.Println("list of movies is getted")
}

// GetMoviesByName - обрабатывает http запрос на получение списка фильмов из фильмотеки по имени.
//
// @Summary      Get movies from the System by name.
// @Description  Get movies from the System by name.
// @Tags         Movie
// @Produce      json
// @Param        name path string true "Name of the movie to be getted"
// @Security BasicAuth
// @Success      200 {array} models.Movie
// @Failure      500 {string} string "Internal server error"
// @Failure      403 {string} string "User does not exist"
// @Router       /movies/name/{name} [get]
func (app *App) GetMoviesByName(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of movies by searching the name")
	_, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	name := r.PathValue("name")
	movies, err := app.dbHandler.GetMoviesByName(name)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, movies)
	app.infoLog.Println("list of movies is getted by searching the name")
}

// GetMoviesByActor - обрабатывает http запрос на получение списка фильмов из фильмотеки по актёру.
//
// @Summary      Get movies from the System by actor.
// @Description  Get movies from the System by actor.
// @Tags         Movie
// @Produce      json
// @Param        actor path string true "Name of the actor to be getted"
// @Security BasicAuth
// @Success      200 {array} models.Movie
// @Failure      500 {string} string "Internal server error"
// @Failure      403 {string} string "User does not exist"
// @Router       /movies/actor/{actor} [get]
func (app *App) GetMoviesByActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of movies by searching the actor")
	_, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	actor := r.PathValue("actor")
	movies, err := app.dbHandler.GetMoviesByActor(actor)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, movies)
	app.infoLog.Println("list of movies is getted by searching the actor")
}

// AddUser - обрабатывает http запрос на добавление пользователя в фильмотеку.
//
// @Summary      Adds user to the System.
// @Description  Add user to the System and get it's ID. User should be an admin.
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User to be added"
// @Security BasicAuth
// @Success      200 {integer} int "ID of the added user"
// @Failure      400 {string} string "Bad request"
// @Failure      403 {string} string "User not an admin"
// @Failure      500 {string} string "Internal server error"
// @Router       /users [post]
func (app *App) AddUser(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to add a new user")

	var user models.User
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&user); err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := app.dbHandler.AddUser(user.Nickname, string(hashedPassword), user.IsAdmin)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, id)
	app.infoLog.Printf("user %d is added\n", id)
}
