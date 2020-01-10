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
	. "net/http"
	"os"
	"time"
)

// TODO:
// *try to add sll mode to postgres database (docker)

type appConfig struct {
	// Relative migration directory path
	migrationsDir string

	// Data base config
	dbConfig dbConfig

	host string
	port string

	// Api prefix
	apiPrefix string

	timeout time.Duration
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
	server(config)
}

func server(config appConfig) {
	r := mux.NewRouter()

	sr := r.PathPrefix(config.apiPrefix).Subrouter()
	sr.HandleFunc("/users/{id:[0-9]+}", userHandler).
		Methods(MethodGet, MethodPost, MethodDelete)

	domainUrl := fmt.Sprintf("%s:%s", config.host, config.port)
	srv := &Server{
		Addr:         domainUrl,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}

	// TODO check out graceful shutdown

	//go func() {
	//	if err := srv.ListenAndServe(); err != nil {
	//		log.Println(err)
	//	}
	//}()
	//
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//<-c
	//
	//// Create a deadline to wait for.
	//ctx, cancel := context.WithTimeout(context.Background(), config.timeout)
	//defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	log.Fatal(err)
	//}
	os.Exit(0)
}

func configure(config *appConfig) {
	var migrationsDir = flag.String("mdir", "./db/migrations", "Directory where the migration files are located")
	var host = flag.String("h", "localhost", "Application deployment host")
	var port = flag.String("p", "3000", "Application deployment port")
	var prefix = flag.String("prefix", "/api/v1/", "Special prefix for api paths")
	var timeout = flag.Duration("timeout", time.Second*15, "The duration for which the server gracefully wait for existing connections to finish")

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
	config.host = *host
	config.port = *port
	config.apiPrefix = *prefix
	config.timeout = *timeout
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
