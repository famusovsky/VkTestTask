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

func (app *App) AddActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to add a new actor")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func (app *App) UpdateActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to update an actor")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		app.infoLog.Println("user trying to update an actor not an admin")
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
func (app *App) DeleteActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to delete an actor")
	isAdmin, err := app.authIsAdmin(r)
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func (app *App) GetActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get an actor")

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
func (app *App) GetActors(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of actors")

	actors, err := app.dbHandler.GetActors()
	if err != nil {
		app.errorLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.sendJson(w, actors)
	app.infoLog.Println("list of actors is getted")
}

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
func (app *App) GetMovies(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of movies")

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
func (app *App) GetMoviesByName(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of movies by searching the name")

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
func (app *App) GetMoviesByActor(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("trying to get list of movies by searching the actor")

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
