package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Entry is the database object for a leaderboard entry.
type Entry struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Countries int    `json:"countries"`
	Time      int    `json:"time"`
}

// EntriesDto is used to display a paged result of leaderboard entries.
type EntriesDto struct {
	Entries []Entry `json:"entries"`
	HasMore bool    `json:"hasMore"`
}

// GetEntry gets a leaderboard entry by id.
func GetEntry(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
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

	entry := Entry{id, name, country, countries, time}
	json.NewEncoder(writer).Encode(entry)
}

// GetEntries gets the leaderboard entries for a given page.
func GetEntries(writer http.ResponseWriter, request *http.Request) {
	pageParam := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	database, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT * FROM leaderboard ORDER BY countries DESC, time limit 10 offset " + strconv.Itoa(page*10))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
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
			fmt.Fprintf(writer, "%v\n", err)
		}

		var entry = Entry{id, name, country, countries, time}
		entries = append(entries, entry)
	}

	statement, _ := database.Prepare("SELECT id FROM leaderboard ORDER BY countries DESC, time limit 1 offset ?")
	defer statement.Close()

	var id string
	hasMore := true
	err = statement.QueryRow(strconv.Itoa((page + 1) * 10)).Scan(&id)
	if err != nil {
		hasMore = false
	}

	entriesDto := EntriesDto{entries, hasMore}
	json.NewEncoder(writer).Encode(entriesDto)
}

// CreateEntry creates a new leaderboard entry.
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
	newEntry.ID = id
	json.NewEncoder(writer).Encode(newEntry)
}

// UpdateEntry updates an existing leaderboard entry.
func UpdateEntry(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
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

	updatedEntry.ID = id
	json.NewEncoder(writer).Encode(updatedEntry)
}

// DeleteEntry deletes an existing leaderboard entry.
func DeleteEntry(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
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

	entry := Entry{id, name, country, countries, time}
	json.NewEncoder(writer).Encode(entry)
}
