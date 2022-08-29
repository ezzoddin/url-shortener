package router

import (
	"github.com/gorilla/mux"
	"url-shortner/controllers"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/url/{id}", controllers.GetUrl).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/url", controllers.GetAllUrls).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/add", controllers.CreateUrl).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/url/{id}", controllers.UpdateUrl).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deleteurl/{id}", controllers.DeleteUrl).Methods("DELETE", "OPTIONS")

	router.HandleFunc("/api/short", controllers.GenerateShortLink).Methods("POST", "OPTIONS")

	return router
}
