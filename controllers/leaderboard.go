package controllers

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "../models"
)

func GetEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        fmt.Fprintf(writer, "Id value is not an int.")
    }

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer database.Close()

    statement, err := database.Prepare("SELECT * FROM leaderboard WHERE id = ?")
	if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer statement.Close()

    var entry models.Entry
    err = statement.QueryRow(id).Scan(&entry)
    if err != nil {
        fmt.Fprintf(writer, err.Error())
    }

    json.NewEncoder(writer).Encode(entry)
}

func GetEntries(writer http.ResponseWriter, request *http.Request) {
    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer database.Close()

    rows, err := database.Query("SELECT * FROM leaderboard")
    if err != nil {
        fmt.Fprintf(writer, err.Error())
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var name string
        var country string
        var countries int
        var time string
        err = rows.Scan(&id, &name, &country, &countries, &time)
        if err != nil {
            fmt.Fprintf(writer, err.Error())
        }
    }

    json.NewEncoder(writer).Encode(rows)
}

func CreateEntry(writer http.ResponseWriter, request *http.Request) {
    requestBody, err := ioutil.ReadAll(request.Body)
    if err != nil {
        fmt.Fprintf(writer, "Request body invalid.")
    }

    var newEntry models.Entry
    json.Unmarshal(requestBody, &newEntry)

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer database.Close()

    statement, err := database.Prepare("INSERT into leaderboard (name, country, countries, time) values (?, ?, ?, ?)")
	if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer statement.Close()

    statement.Exec(newEntry.Name, newEntry.Country, newEntry.Countries, newEntry.Time)

    writer.WriteHeader(http.StatusCreated)
    json.NewEncoder(writer).Encode(newEntry)
}

func UpdateEntry(writer http.ResponseWriter, request *http.Request) {
    id, parseErr := strconv.Atoi(mux.Vars(request)["id"])
    if(parseErr != nil) {
        fmt.Fprintf(writer, "Id value is not an int.")
    }

    requestBody, readErr := ioutil.ReadAll(request.Body)
    if readErr != nil {
        fmt.Fprintf(writer, "Request body invalid.")
    }

    var updatedEntry models.Entry
    json.Unmarshal(requestBody, &updatedEntry)

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer database.Close()

    statement, err := database.Prepare("UPDATE leaderboard set name = ?, country = ?, countries = ?, time = ? where id = ?")
	if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer statement.Close()

    statement.Exec(updatedEntry.Name, updatedEntry.Country, updatedEntry.Countries, updatedEntry.Time, id)

    updatedEntry.Id = id
    json.NewEncoder(writer).Encode(updatedEntry)
}

func DeleteEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        fmt.Fprintf(writer, "Id value is not an int.")
    }

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer database.Close()

    statement, err := database.Prepare("DELETE FROM leaderboard WHERE id = ?")
	if err != nil {
        fmt.Fprintf(writer, err.Error())
	}
    defer statement.Close()

    statement.Exec(id)
}

