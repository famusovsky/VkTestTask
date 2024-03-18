package filmoteka

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// routes - создание маршрутов.
func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /actor", app.AddActor)
	mux.HandleFunc("PUT /actor/{id}", app.UpdateActor)
	mux.HandleFunc("DELETE /actor/{id}", app.DeleteActor)
	mux.HandleFunc("GET /actor/{id}", app.GetActor)
	mux.HandleFunc("GET /actors", app.GetActors)

	mux.HandleFunc("POST /movie", app.AddMovie)
	mux.HandleFunc("DELETE /movie/{id}", app.DeleteMovie)
	mux.HandleFunc("PUT /movie/{id}", app.UpdateMovie)

	mux.HandleFunc("GET /movies", app.GetMovies)
	mux.HandleFunc("GET /movies/name/{name}", app.GetMoviesByName)
	mux.HandleFunc("GET /movies/actor/{actor}", app.GetMoviesByActor)

	mux.HandleFunc("POST /users", app.AddUser)

	return mux
}
