package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "io/ioutil"
    "encoding/json"
    "strconv"
)

type entry struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Country string `json:"country"`
    Countries int `json:"countries"`
    Time string `json:"time"`
}

var entries = []entry{}

func getEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        fmt.Fprintf(writer, "Id value is not an int.")
        return
    }

    for _, value := range entries {
        if value.Id == id {
            json.NewEncoder(writer).Encode(value)
            return
        }
    }

    writer.WriteHeader(http.StatusNotFound)
}

func getEntries(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(entries)
}

func createEntry(writer http.ResponseWriter, request *http.Request) {
    requestBody, err := ioutil.ReadAll(request.Body)
    if err != nil {
        fmt.Fprintf(writer, "Request body invalid.")
    }

    var newEntry entry
    json.Unmarshal(requestBody, &newEntry)
    entries = append(entries, newEntry)
    writer.WriteHeader(http.StatusCreated)
    json.NewEncoder(writer).Encode(newEntry)
}

func updateEntry(writer http.ResponseWriter, request *http.Request) {
    id, parseErr := strconv.Atoi(mux.Vars(request)["id"])
    if(parseErr != nil) {
        fmt.Fprintf(writer, "Id value is not an int.")
        return
    }

    requestBody, readErr := ioutil.ReadAll(request.Body)
    if readErr != nil {
        fmt.Fprintf(writer, "Request body invalid.")
        return
    }

    var updatedEntry entry
    json.Unmarshal(requestBody, &updatedEntry)

    for i, value := range entries {
        if value.Id == id {
            value.Name = updatedEntry.Name
            value.Country = updatedEntry.Country
            value.Countries = updatedEntry.Countries
            value.Time = updatedEntry.Time
            entries = append(entries[:i], value)
            json.NewEncoder(writer).Encode(value)
            return
        }
    }

    writer.WriteHeader(http.StatusNotFound)
}

func deleteEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        fmt.Fprintf(writer, "Id value is not an int.")
        return
    }

    for i, value := range entries {
        if value.Id == id {
            entries = append(entries[:i], entries[i+1:]...)
            json.NewEncoder(writer).Encode(value)
            return
        }
    }

    writer.WriteHeader(http.StatusNotFound)
}

func main() {
    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/entries", getEntries).Methods("GET")
    router.HandleFunc("/entries/{id}", getEntry).Methods("GET")
    router.HandleFunc("/entries", createEntry).Methods("POST")
    router.HandleFunc("/entries/{id}", updateEntry).Methods("PATCH")
	router.HandleFunc("/entries/{id}", deleteEntry).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8080", router))
}

