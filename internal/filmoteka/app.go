// Пакет filmoteka реализует логику приложения Фильмотеки.
package filmoteka

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/famusovsky/VkTestTask/internal/filmoteka/postgres"
	"golang.org/x/sync/errgroup"
)

// App - модель приложения.
type App struct {
	infoLog   *log.Logger
	errorLog  *log.Logger
	dbHandler postgres.DbHandler
	addr      string
}

// CreateApp - создание приложения.
//
// Принимает: логгер, обработчик БД.
//
// Возвращает: приложение.
func CreateApp(addr string, infoLog *log.Logger, errorLog *log.Logger,
	dbHandler postgres.DbHandler) *App {
	return &App{
		infoLog:   infoLog,
		errorLog:  errorLog,
		dbHandler: dbHandler,
		addr:      addr,
	}
}

// Run - запуск приложения.
//
// Принимает: адрес.
func (app *App) Run() {
	// Создание и запуск сервера.
	srvr := &http.Server{
		Addr:     app.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		s := <-sigQuit
		return fmt.Errorf("captured signal: %v", s)
	})

	go func() {
		app.infoLog.Printf("starting srvr on %s\n", app.addr)
		err := srvr.ListenAndServe()
		app.errorLog.Fatal(err)
	}()

	if err := eg.Wait(); err != nil {
		app.infoLog.Printf("gracefully shutting down the server: %v", err)
	}

	err := srvr.Shutdown(context.Background())
	app.errorLog.Fatal(err)
}
