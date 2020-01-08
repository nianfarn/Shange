package main

import (
	"database/sql"
	"flag"
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
//  flag URL_PREFIX = "api/v1/"

type appConfig struct {
	// Relative migration directory path
	migrationsDir string

	// Data base config
	dbConfig dbConfig
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

	//http.Handle(urlPrefix + "users", http.HandlerFunc(userHandler))
	//
	//entryUrl := ":" + strconv.Itoa(port)
	//if err := http.ListenAndServe(entryUrl, http.HandlerFunc(entryHandler)); err != nil {
	//	log.Fatalf("could not listen on port %v with: %v", port, err)
	//}
}

func configure(config *appConfig) {
	var migrationsDir = flag.String("mdir", "./db/migrations", "Directory where the migration files are located")

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
}

func applyMigrations(config appConfig) {
	// todo use newest postgres driver
	//config, err := pgx.ParseConnectionString(dataSourceName)
	//if err != nil {
	//
	//}
	//sql.Register("postgres-pgx", pgx.ParseConnectionString(dataSourceName))
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
