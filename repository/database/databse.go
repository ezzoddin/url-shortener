package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func Connect() *sql.DB {
	config, err := LoadConfig(".")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Pass, config.Name, config.Sslmode)

	db, err := sql.Open(config.Driver, url)

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully connected!")

	return db
}
