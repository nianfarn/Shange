package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// TODO:
//  flag PORT = 3000
//  flag URL_PREFIX = "api/v1/"
// *try to add sll mode to postgres database (docker)

type appConfig struct {
	// Relative migration directory path
	migrationsDir string

	// Data base config
	dbConfig dbConfig

	port string

	// Api prefix
	apiPrefix string
}

type dbConfig struct {
	// DB type (postgres or another)
	dbType string

	// Data source
	dataSource string
}

func main() {
	config := appConfig{}
	configure(&config)

	applyMigrations(config)
	setupDispatcher(config)

	startServer(config)
}

func startServer(config appConfig) {
	if err := http.ListenAndServe(":"+config.port, http.HandlerFunc(entryHandler)); err != nil {
		log.Fatalf("could not listen on port %v with: %v", config.port, err)
	}
}

func setupDispatcher(config appConfig) {
	r := mux.NewRouter()

	//r.Methods().Subrouter()
	r.HandleFunc("/users/{id:[0-9]+}", userHandler).
		Methods("GET", "POST", "DELETE")

	http.Handle("/", r)
}

func configure(config *appConfig) {
	var migrationsDir = flag.String("mdir", "./db/migrations", "Directory where the migration files are located")
	var port = flag.String("p", "3000", "Application deployment port")
	var prefix = flag.String("prefix", "api/v1/", "Special prefix for api paths")

	var dbType = flag.String("db.type", "postgres", "Database type")
	var dbUser = flag.String("db.user", "postuser", "Database user")
	var dbPwd = flag.String("db.pwd", "postpass", "Database user password")
	var dbHost = flag.String("db.host", "localhost", "Database host")
	var dbPort = flag.String("db.port", "5432", "Database port")
	var dbName = flag.String("db.name", "shange-db", "Database name")

	flag.Parse()

	ds := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", *dbType, *dbUser, *dbPwd, *dbHost, *dbPort, *dbName)

	config.migrationsDir = *migrationsDir
	config.dbConfig = dbConfig{dbType: *dbType, dataSource: ds}
	config.port = *port
	config.apiPrefix = *prefix
}

func applyMigrations(config appConfig) {
	db, err := sql.Open(config.dbConfig.dbType, config.dbConfig.dataSource)
	if err != nil {
		log.Fatalf("could not connect to the database... %v", err)
	}
	defer db.Close()

	// Run migrations
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.migrationsDir),
		config.dbConfig.dbType, driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("an error occurred while procceding migration: %v", err)
	}

	log.Println("Database successfully migrated")
}
