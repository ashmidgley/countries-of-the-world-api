package models

type Entry struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Country string `json:"country"`
    Countries int `json:"countries"`
    Time string `json:"time"`
}
