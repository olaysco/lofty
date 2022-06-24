package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/olaysco/evolve/api"
	"github.com/olaysco/evolve/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("error loading config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatal("error connecting to db:", err)
	}

	migrateDB(config.MigrationURL, config.DBUrl)
	serveHTTP(config, conn)
}

func migrateDB(migrationURL string, dbUrl string) {
	migration, err := migrate.New(migrationURL, dbUrl)
	if err != nil {
		log.Fatal("error migrating DB:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("error running migrate up:", err)
	}

	log.Println("db migrated successfully")
}

func serveHTTP(config util.Config, db *sql.DB) {
	log.Printf("server listening at addreess %s", config.HTTPServerAddress)
	server := api.NewAPi(config, db)
	err := server.Run()

	if err != nil {
		log.Fatalf("error starting server at address %s: %s", config.HTTPServerAddress, err.Error())
	}
}
