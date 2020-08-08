package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "./controllers"
)

func main() {
    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
		log.Fatal(err)
	}
    defer database.Close()

    statement := "CREATE TABLE IF NOT EXISTS leaderboard (id INTEGER PRIMARY KEY, name TEXT, country TEXT, countries INTEGER, time TEXT)"
    _, err = database.Exec(statement)
    if err != nil {
		log.Printf("%q: %s\n", err, statement)
		return
	}

    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/api/leaderboard", controllers.GetEntries).Methods("GET")
    router.HandleFunc("/api/leaderboard/{id}", controllers.GetEntry).Methods("GET")
    router.HandleFunc("/api/leaderboard", controllers.CreateEntry).Methods("POST")
    router.HandleFunc("/api/leaderboard/{id}", controllers.UpdateEntry).Methods("PUT")
    router.HandleFunc("/api/leaderboard/{id}", controllers.DeleteEntry).Methods("DELETE")
    router.HandleFunc("/api/countries", controllers.GetCountries).Methods("GET")
    router.HandleFunc("/api/countries/map", controllers.GetCountriesMap).Methods("GET")
    router.HandleFunc("/api/codes", controllers.GetCodes).Methods("GET")

    handler := cors.Default().Handler(router)
    http.ListenAndServe(":8080", handler)
}
