package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "./controllers"
)

func main() {
    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/api/leaderboard", controllers.GetEntries).Methods("GET")
    router.HandleFunc("/api/leaderboard/{id}", controllers.GetEntry).Methods("GET")
    router.HandleFunc("/api/leaderboard", controllers.CreateEntry).Methods("POST")
    router.HandleFunc("/api/leaderboard/{id}", controllers.UpdateEntry).Methods("PATCH")
	router.HandleFunc("/api/leaderboard/{id}", controllers.DeleteEntry).Methods("DELETE")
    router.HandleFunc("/api/countries", controllers.GetCountries).Methods("GET")
    router.HandleFunc("/api/countries/map", controllers.GetCountriesMap).Methods("GET")
    router.HandleFunc("/api/codes", controllers.GetCodes).Methods("GET")

    handler := cors.Default().Handler(router)
    http.ListenAndServe(":8080", handler)
}

