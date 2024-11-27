package models

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_Test string
}

func newTestDB(t *testing.T) *sql.DB {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file!")
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal("Error assigning .env file values!")
	}

	db, err := sql.Open("mysql", cfg.DB_Test)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DB_Test: os.Getenv("TEST_DB"),
	}

	return cfg, nil
}
