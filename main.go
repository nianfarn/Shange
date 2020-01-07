package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
)

// TODO:
//  try to add sll mode to postgres database (docker)
//  use config var as flags
//  flag PORT = 3000
//  flag PORT = 3000
//  flag URL_PREFIX = "api/v1/"

const migrationsDir = "./db/migrations"
const dataSourceName = "postgres://postuser:postpass@localhost:5432/shange-db?sslmode=disable"

func main() {
	migrations()

	//http.Handle(urlPrefix + "users", http.HandlerFunc(userHandler))
	//
	//entryUrl := ":" + strconv.Itoa(port)
	//if err := http.ListenAndServe(entryUrl, http.HandlerFunc(entryHandler)); err != nil {
	//	log.Fatalf("could not listen on port %v with: %v", port, err)
	//}
}

func migrations() {
	//var migrationDir = flag.String("migration.files", migrationsDir, "Directory where the migration files are located")
	//var dataSource = flag.String("postgres.dsn", dataSourceName, "Postgres data source")
	//flag.Parse()

	// todo use newest postgres driver
	//config, err := pgx.ParseConnectionString(dataSourceName)
	//if err != nil {
	//
	//}
	//sql.Register("postgres-pgx", pgx.ParseConnectionString(dataSourceName))
	db, err := sql.Open("postgres", dataSourceName)
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
		fmt.Sprintf("file://%s", migrationsDir),
		"postgres", driver)
	m.Steps(2)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("an error occurred while procceding migration: %v", err)
	}

	log.Println("Database successfully migrated")
}
