package api

import (
	"flag"
	"log"
	"os"

	"github.com/famusovsky/VkTestTask/internal/filmoteka"
	"github.com/famusovsky/VkTestTask/internal/filmoteka/postgres"
	"github.com/famusovsky/VkTestTask/pkg/database"
)

// @title Filemoteka API
// @description This is a Filmoteka API server, made for Vk Trainee Assignment 2024.
func main() {
	addr := flag.String("addr", ":8080", "HTTP address")
	overrideTables := flag.Bool("override_tables", false, "Override tables in database")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	db, err := database.OpenViaEnvVars("postgres")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	dbHandler, err := postgres.GetHandler(db, *overrideTables)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := filmoteka.CreateApp(*addr, infoLog, errorLog, dbHandler)

	app.Run(*addr)
}
