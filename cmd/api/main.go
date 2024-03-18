package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/famusovsky/VkTestTask/docs"
	"github.com/famusovsky/VkTestTask/internal/filmoteka"
	"github.com/famusovsky/VkTestTask/internal/filmoteka/postgres"
	"github.com/famusovsky/VkTestTask/pkg/database"
	_ "github.com/lib/pq"
)

// TODO testing

// @title			Filemoteka API
// @description	This is a Filmoteka API server, made for Vk Trainee Assignment 2024.
// @securityDefinitions.basic  BasicAuth
func main() {
	addr := flag.String("addr", ":8080", "HTTP address")
	overrideTables := flag.Bool("override_tables", false, "Override tables in database")
	defaultAdmin := flag.Bool("default_admin", false, "Add default admin (admin|admin) to database")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	db, err := database.OpenViaEnvVars("postgres", sql.Open)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	dbHandler, err := postgres.GetHandler(db, *overrideTables)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := filmoteka.CreateApp(*addr, infoLog, errorLog, dbHandler, *defaultAdmin)

	app.Run()
}
