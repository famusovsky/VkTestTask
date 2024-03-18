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
	defAdmin  bool
}

// CreateApp - создание приложения.
//
// Принимает: адрес, логгер информации, логгер ошибок, обработчик БД, указатель на существование базового администратора.
//
// Возвращает: приложение.
func CreateApp(addr string, infoLog *log.Logger, errorLog *log.Logger,
	dbHandler postgres.DbHandler, defAdmin bool) *App {
	return &App{
		infoLog:   infoLog,
		errorLog:  errorLog,
		dbHandler: dbHandler,
		addr:      addr,
		defAdmin:  defAdmin,
	}
}

// server - интерфейс сервера.
type server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

// Run - запуск приложения.
func (app *App) Run() {
	// Создание и запуск сервера.
	srvr := &http.Server{
		Addr:     app.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(),
	}

	GraceRun(srvr, app)
}

// GraceRun - грациозный запуск сервера.
//
// Принимает: сервер, приложение.
//
// Запускает сервер и ожидает сигнала завершения.
func GraceRun(srvr server, app *App) {
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
		if err != nil {
			app.errorLog.Fatal(err)
		}
	}()

	if err := eg.Wait(); err != nil {
		app.infoLog.Printf("gracefully shutting down the server: %v", err)
	}

	err := srvr.Shutdown(context.Background())
	if err != nil {
		app.errorLog.Fatal(err)
	}
}
