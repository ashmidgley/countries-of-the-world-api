package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ashmidgley/countries-of-the-world-api/database"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	statement := "SELECT * FROM leaderboard WHERE id = $1;"
	row := database.DBConnection.QueryRow(statement, id)

	var entry Entry
	switch err = row.Scan(&entry.ID, &entry.Name, &entry.Country, &entry.Countries, &entry.Time); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		json.NewEncoder(writer).Encode(entry)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
}

// GetEntries gets the leaderboard entries for a given page.
func GetEntries(writer http.ResponseWriter, request *http.Request) {
	pageParam := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	rows, err := database.DBConnection.Query("SELECT * FROM leaderboard ORDER BY countries DESC, time LIMIT $1 OFFSET $2;", 10, strconv.Itoa(page*10))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}
	defer rows.Close()

	var entries = []Entry{}
	for rows.Next() {
		var entry Entry
		err = rows.Scan(&entry.ID, &entry.Name, &entry.Country, &entry.Countries, &entry.Time)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		entries = append(entries, entry)
	}

	err = rows.Err()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	statement := "SELECT id FROM leaderboard ORDER BY countries DESC, time LIMIT $1 OFFSET $2;"
	row := database.DBConnection.QueryRow(statement, 1, strconv.Itoa((page+1)*10))

	var id int
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		entriesDto := EntriesDto{entries, false}
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := EntriesDto{entries, true}
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
}

// CreateEntry creates a new leaderboard entry.
func CreateEntry(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	var newEntry Entry
	err = json.Unmarshal(requestBody, &newEntry)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	statement := "INSERT INTO leaderboard (name, country, countries, time) VALUES ($1, $2, $3, $4) RETURNING id;"

	var id int
	err = database.DBConnection.QueryRow(statement, newEntry.Name, newEntry.Country, newEntry.Countries, newEntry.Time).Scan(&id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
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
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	var updatedEntry Entry
	err = json.Unmarshal(requestBody, &updatedEntry)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	statement := "UPDATE leaderboard set name = $2, country = $3, countries = $4, time = $5 where id = $1 RETURNING *;"

	row := database.DBConnection.QueryRow(statement, id, updatedEntry.Name, updatedEntry.Country, updatedEntry.Countries, updatedEntry.Time)

	var entry Entry
	switch err = row.Scan(&entry.ID, &entry.Name, &entry.Country, &entry.Countries, &entry.Time); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		json.NewEncoder(writer).Encode(entry)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
}

// DeleteEntry deletes an existing leaderboard entry.
func DeleteEntry(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	statement := "DELETE FROM leaderboard WHERE id = $1 RETURNING *;"
	row := database.DBConnection.QueryRow(statement, id)

	var entry Entry
	switch err = row.Scan(&entry.ID, &entry.Name, &entry.Country, &entry.Countries, &entry.Time); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		json.NewEncoder(writer).Encode(entry)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
}
