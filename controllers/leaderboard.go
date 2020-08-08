package controllers

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "../models"
)

var entries = []models.Entry{}

func GetEntry(writer http.ResponseWriter, request *http.Request) {
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

func GetEntries(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(entries)
}

func CreateEntry(writer http.ResponseWriter, request *http.Request) {
    requestBody, err := ioutil.ReadAll(request.Body)
    if err != nil {
        fmt.Fprintf(writer, "Request body invalid.")
    }

    var newEntry models.Entry
    json.Unmarshal(requestBody, &newEntry)
    entries = append(entries, newEntry)
    writer.WriteHeader(http.StatusCreated)
    json.NewEncoder(writer).Encode(newEntry)
}

func UpdateEntry(writer http.ResponseWriter, request *http.Request) {
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

    var updatedEntry models.Entry
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

func DeleteEntry(writer http.ResponseWriter, request *http.Request) {
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

