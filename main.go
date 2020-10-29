package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ashmidgley/countries-of-the-world-api/database"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	var err error
	database.DBConnection, err = sql.Open("postgres", database.GetConnectionString())
	if err != nil {
		panic(err)
	}

	err = database.DBConnection.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/leaderboard", GetEntries).Methods("GET")
	router.HandleFunc("/api/leaderboard/{id}", GetEntry).Methods("GET")
	router.HandleFunc("/api/leaderboard", CreateEntry).Methods("POST")
	router.HandleFunc("/api/leaderboard/{id}", UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/leaderboard/{id}", DeleteEntry).Methods("DELETE")
	router.HandleFunc("/api/countries", GetCountries).Methods("GET")
	router.HandleFunc("/api/countries/alternatives", GetAlternativeNamings).Methods("GET")
	router.HandleFunc("/api/countries/prefixes", GetPrefixes).Methods("GET")
	router.HandleFunc("/api/countries/map", GetCountriesMap).Methods("GET")
	router.HandleFunc("/api/codes", GetCodes).Methods("GET")

	handler := cors.Default().Handler(router)
	http.ListenAndServe(":8080", handler)
}
