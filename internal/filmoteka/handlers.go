package filmoteka

import (
	"net/http"
)

// TODO
// routes - создание маршрутов.
func (app *App) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// TODO
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {})

	return mux
}
