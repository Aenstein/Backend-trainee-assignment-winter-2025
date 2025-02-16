package services

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var db *sql.DB
var once sync.Once

var MerchItems = map[string]int64{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

func InitDB() error {
	var err error
	once.Do(func() {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return
		}
		err = db.Ping()
		if err != nil {
			return
		}
		err = createTables()
	})
	return err
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			coins INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS inventory (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			item_type TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			UNIQUE(user_id, item_type)
		);`,
		`CREATE TABLE IF NOT EXISTS coin_transactions (
			id SERIAL PRIMARY KEY,
			from_user TEXT,
			to_user TEXT,
			amount INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}