package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "./controllers"
)

func main() {
    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/api/leaderboard", controllers.GetEntries).Methods("GET")
    router.HandleFunc("/api/leaderboard/{id}", controllers.GetEntry).Methods("GET")
    router.HandleFunc("/api/leaderboard", controllers.CreateEntry).Methods("POST")
    router.HandleFunc("/api/leaderboard/{id}", controllers.UpdateEntry).Methods("PATCH")
	router.HandleFunc("/api/leaderboard/{id}", controllers.DeleteEntry).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8080", router))
}

