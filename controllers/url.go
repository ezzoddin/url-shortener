package controllers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/itchyny/base58-go"
	_ "github.com/lib/pq" // postgres golang driver
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"url-shortner/models"
	"url-shortner/repository/database"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var url models.Url
	err := json.NewDecoder(r.Body).Decode(&url)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	url.CreatedAt = time.Now()
	url.UpdatedAt = time.Now()

	insertID := insertUrl(url)

	res := response{
		ID:      insertID,
		Message: "Url created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	url, err := getUrl(int64(id))

	if err != nil {
		log.Fatalf("Unable to get url. %v", err)
	}

	json.NewEncoder(w).Encode(url)
}

func GetAllUrls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	urls, err := getAllUrls()

	if err != nil {
		log.Fatalf("Unable to get all urls. %v", err)
	}

	json.NewEncoder(w).Encode(urls)
}

func UpdateUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var url models.Url

	err = json.NewDecoder(r.Body).Decode(&url)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateUrl(int64(id), url)

	msg := fmt.Sprintf("Url updated successfully. Total rows/record affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteUrl(int64(id))

	msg := fmt.Sprintf("Url updated successfully. Total rows/record affected %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertUrl(url models.Url) int64 {

	db := database.Connect()

	sqlStatement := `INSERT INTO ` + models.GetTable() + ` (title, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id`

	err := db.QueryRow(sqlStatement, url.Title, url.CreatedAt, url.UpdatedAt).Scan(&url.Id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", url.Id)

	return url.Id
}

func getUrl(id int64) (models.Url, error) {

	db := database.Connect()

	defer db.Close()

	var url models.Url

	sqlStatement := `SELECT * FROM ` + models.GetTable() + ` WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&url.Id, &url.Title, &url.CreatedAt, &url.UpdatedAt)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return url, nil
	case nil:
		return url, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return url, err
}

func getAllUrls() ([]models.Url, error) {
	db := database.Connect()

	var urls []models.Url

	sqlStatement := `SELECT * FROM ` + models.GetTable()

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var url models.Url

		err = rows.Scan(&url.Id, &url.Title, &url.CreatedAt, &url.UpdatedAt)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		urls = append(urls, url)
	}

	return urls, err
}

func updateUrl(id int64, url models.Url) int64 {

	db := database.Connect()

	defer db.Close()

	sqlStatement := `UPDATE ` + models.GetTable() + ` SET title=$2 WHERE id=$1`

	res, err := db.Exec(sqlStatement, id, url.Title)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func deleteUrl(id int64) int64 {
	db := database.Connect()

	defer db.Close()

	sqlStatement := `DELETE FROM ` + models.GetTable() + ` WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

func isUrlValid(input string) bool {
	validUrl, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	switch validUrl.Scheme {
	case "https":
	case "http":
	default:
		return false
	}

	_, err = net.LookupHost(validUrl.Host)
	if err != nil {
		return false
	}

	return true
}

func GenerateShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var url models.Url
	err := json.NewDecoder(r.Body).Decode(&url)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	url.CreatedAt = time.Now()
	url.UpdatedAt = time.Now()

	insertID := insertUrl(url)

	shortLink := ShortLink(url.Title, insertID)

	res := response{
		ID:      insertID,
		Message: "short url" + shortLink,
	}

	json.NewEncoder(w).Encode(res)
}

func ShortLink(initialLink string, urlId int64) string {
	id := strconv.FormatInt(urlId, 16)
	urlHashBytes := sha256Of(initialLink + id)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
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
