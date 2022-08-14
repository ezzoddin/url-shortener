package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	DB_USER     = "mobina"
	DB_PASSWORD = "++"
	DB_NAME     = "golang"
)

type ShortUrl struct {
	url       string
	createdAt time.Time
}

var database *sql.DB

func main() {
	connectDB()
	var sampleUrl = ShortUrl{url: "https://mezz.ir", createdAt: time.Now()}
	insert(sampleUrl)
}

func connectDB() {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", psqlInfo)
	checkErr(err)

	database = db
}

func insert(shortUrl ShortUrl) int {
	id := 0
	sqlStatement := `
					INSERT INTO short_urls (url, created_at)
					VALUES ($1, $2)
					RETURNING id`
	err := database.QueryRow(sqlStatement, shortUrl.url, time.Now()).Scan(&id)
	checkErr(err)

	return id
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
