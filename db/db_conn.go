package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/pranesh/bitespeed/home"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("postgres", home.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatal("db conn. error:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("db ping fail:", err)
	}

	runMigrations()
}

func runMigrations() {
	query, err := os.ReadFile("migrations/contacts.sql")
	if err != nil {
		log.Fatal("failed to read migration file:", err)
	}
	if _, err := DB.Exec(string(query)); err != nil {
		log.Fatal("migration error:", err)
	}
}
