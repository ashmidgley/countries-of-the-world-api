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
    "strings"
    "../models"
)

func GetEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, err := database.Prepare("SELECT * FROM leaderboard WHERE id = ?")
	if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer statement.Close()

    var name string
    var country string
    var countries int
    var time string
    err = statement.QueryRow(id).Scan(&id, &name, &country, &countries, &time)
    if err != nil {
        message := err.Error()
        if strings.Contains(message, "no rows in result set") {
            writer.WriteHeader(http.StatusNotFound)
        } else {
            writer.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(writer, message)
        }
        return
    }

    entry := models.Entry { id, name, country, countries, time }
    json.NewEncoder(writer).Encode(entry)
}

func GetEntries(writer http.ResponseWriter, request *http.Request) {
    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    rows, err := database.Query("SELECT * FROM leaderboard")
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
    }
    defer rows.Close()

    var entries = []models.Entry{}
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

        var entry = models.Entry { id, name, country, countries, time };
        entries = append(entries, entry)
    }

    json.NewEncoder(writer).Encode(entries)
}

func CreateEntry(writer http.ResponseWriter, request *http.Request) {
    requestBody, err := ioutil.ReadAll(request.Body)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    var newEntry models.Entry
    err = json.Unmarshal(requestBody, &newEntry)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, err := database.Prepare("INSERT into leaderboard (name, country, countries, time) values (?, ?, ?, ?)")
	if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer statement.Close()

    statement.Exec(newEntry.Name, newEntry.Country, newEntry.Countries, newEntry.Time)

    writer.WriteHeader(http.StatusCreated)
    json.NewEncoder(writer).Encode(newEntry)
}

func UpdateEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    requestBody, err := ioutil.ReadAll(request.Body)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    var updatedEntry models.Entry
    err = json.Unmarshal(requestBody, &updatedEntry)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, err := database.Prepare("UPDATE leaderboard set name = ?, country = ?, countries = ?, time = ? where id = ?")
	if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer statement.Close()

    statement.Exec(updatedEntry.Name, updatedEntry.Country, updatedEntry.Countries, updatedEntry.Time, id)

    updatedEntry.Id = id
    json.NewEncoder(writer).Encode(updatedEntry)
}

func DeleteEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, err := database.Prepare("DELETE FROM leaderboard WHERE id = ?")
	if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer statement.Close()

    statement.Exec(id)
}

