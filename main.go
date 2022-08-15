package main

import (
	"bufio"
	"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/itchyny/base58-go"
	"net/url"
	"os"
	"strings"
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
	baseUrl := getConsoleUrl()
	validateUrl(baseUrl)

	id := insert(baseUrl)
	idAsString := fmt.Sprint(id)
	generateShortLink(baseUrl.url, idAsString)
	connectDB()
}

func getConsoleUrl() ShortUrl {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter url: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSuffix(url, "\n")

	return ShortUrl{url: url}
}

func validateUrl(inputUrl ShortUrl) {
	_, err := url.ParseRequestURI(inputUrl.url)
	checkErr(err)
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))

	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(encoded)
}

func generateShortLink(baseUrl string, id string) string {
	urlHashBytes := sha256Of(baseUrl + id)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))

	return finalString[:8]
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
