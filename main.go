package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

const connectionString = "./database/leaderboard.db"

func main() {
	database, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	statement := "CREATE TABLE IF NOT EXISTS leaderboard (id INTEGER PRIMARY KEY, name TEXT, country TEXT, countries INTEGER, time INTEGER)"
	_, err = database.Exec(statement)
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
		return
	}

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
