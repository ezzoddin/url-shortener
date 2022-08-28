package main

import (
	"fmt"
	"log"
	"net/http"
	_ "url-shortner/repository/database"
	"url-shortner/router"
)

func main() {
	r := router.Router()

	fmt.Println("Starting server on the port 9090 ...")
	log.Fatal(http.ListenAndServe(":9090", r))
}
