package main

import (
    "fmt"
    "net/http"
    "database/sql"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "strings"
    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

type Entry struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Country string `json:"country"`
    Countries int `json:"countries"`
    Time int `json:"time"`
}

func GetEntry(writer http.ResponseWriter, request *http.Request) {
    id, err := strconv.Atoi(mux.Vars(request)["id"])
    if(err != nil) {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    database, err := sql.Open("sqlite3", connectionString)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, _ := database.Prepare("SELECT * FROM leaderboard WHERE id = ?")
    defer statement.Close()

    var name string
    var country string
    var countries int
    var time int
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

    entry := Entry { id, name, country, countries, time }
    json.NewEncoder(writer).Encode(entry)
}

func GetEntries(writer http.ResponseWriter, request *http.Request) {
    database, err := sql.Open("sqlite3", connectionString)
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

    var entries = []Entry{}
    for rows.Next() {
        var id int
        var name string
        var country string
        var countries int
        var time int
        err = rows.Scan(&id, &name, &country, &countries, &time)
        if err != nil {
            fmt.Fprintf(writer, err.Error())
        }

        var entry = Entry { id, name, country, countries, time };
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

    var newEntry Entry
    err = json.Unmarshal(requestBody, &newEntry)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    database, err := sql.Open("sqlite3", connectionString)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, _ := database.Prepare("INSERT into leaderboard (name, country, countries, time) values (?, ?, ?, ?)")
    defer statement.Close()

    statement.Exec(newEntry.Name, newEntry.Country, newEntry.Countries, newEntry.Time)


    statement, _ = database.Prepare("SELECT id FROM leaderboard WHERE name = ? AND country = ? AND countries = ? AND time = ?")
    defer statement.Close()

    var id int
    err = statement.QueryRow(newEntry.Name, newEntry.Country, newEntry.Countries, newEntry.Time).Scan(&id)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
    }

    writer.WriteHeader(http.StatusCreated)
    newEntry.Id = id
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

    var updatedEntry Entry
    err = json.Unmarshal(requestBody, &updatedEntry)
    if err != nil {
        writer.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(writer, err.Error())
        return
    }

    database, err := sql.Open("sqlite3", connectionString)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, _ := database.Prepare("SELECT id FROM leaderboard WHERE id = ?")
    defer statement.Close()

    var temp int
    err = statement.QueryRow(id).Scan(&temp)
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

    statement, _ = database.Prepare("UPDATE leaderboard set name = ?, country = ?, countries = ?, time = ? where id = ?")
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

    database, err := sql.Open("sqlite3", connectionString)
    if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer database.Close()

    statement, _ := database.Prepare("SELECT * FROM leaderboard WHERE id = ?")
    defer statement.Close()

    var name string
    var country string
    var countries int
    var time int
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

    statement, _ = database.Prepare("DELETE FROM leaderboard WHERE id = ?")
    defer statement.Close()

    statement.Exec(id)

    entry := Entry { id, name, country, countries, time }
    json.NewEncoder(writer).Encode(entry)
}