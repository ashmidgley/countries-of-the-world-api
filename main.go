package main

import (
    "fmt"
    "log"
    "net/http"
    "database/sql"
    "io/ioutil"
    "encoding/json"
    "strconv"
    "strings"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    database, err := sql.Open("sqlite3", "./leaderboard.db")
    if err != nil {
		log.Fatal(err)
	}
    defer database.Close()

    statement := "CREATE TABLE IF NOT EXISTS leaderboard (id INTEGER PRIMARY KEY, name TEXT, country TEXT, countries INTEGER, time INTEGER)"
    _, err = database.Exec(statement)
    if err != nil {
		log.Printf("%q: %s\n", err, statement)
		return
	}

    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/api/leaderboard", getEntries).Methods("GET")
    router.HandleFunc("/api/leaderboard/{id}", getEntry).Methods("GET")
    router.HandleFunc("/api/leaderboard", createEntry).Methods("POST")
    router.HandleFunc("/api/leaderboard/{id}", updateEntry).Methods("PUT")
    router.HandleFunc("/api/leaderboard/{id}", deleteEntry).Methods("DELETE")
    router.HandleFunc("/api/countries", getCountries).Methods("GET")
    router.HandleFunc("/api/countries/alternatives", getAlternativeNamings).Methods("GET")
    router.HandleFunc("/api/countries/map", getCountriesMap).Methods("GET")
    router.HandleFunc("/api/codes", getCodes).Methods("GET")

    handler := cors.Default().Handler(router)
    http.ListenAndServe(":8080", handler)
}

type Entry struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Country string `json:"country"`
    Countries int `json:"countries"`
    Time int `json:"time"`
}

func getEntry(writer http.ResponseWriter, request *http.Request) {
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

func getEntries(writer http.ResponseWriter, request *http.Request) {
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

func createEntry(writer http.ResponseWriter, request *http.Request) {
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

func updateEntry(writer http.ResponseWriter, request *http.Request) {
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

    database, err := sql.Open("sqlite3", "./leaderboard.db")
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

    statement, err = database.Prepare("UPDATE leaderboard set name = ?, country = ?, countries = ?, time = ? where id = ?")
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

func deleteEntry(writer http.ResponseWriter, request *http.Request) {
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

    statement, err = database.Prepare("DELETE FROM leaderboard WHERE id = ?")
	if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer statement.Close()

    statement.Exec(id)

    entry := Entry { id, name, country, countries, time }
    json.NewEncoder(writer).Encode(entry)
}

var countries = [...]string{
    "afghanistan",
    "aland islands",
    "albania",
    "algeria",
    "american samoa",
    "andorra",
    "angola",
    "anguilla",
    "antigua and barbuda",
    "argentina",
    "armenia",
    "aruba",
    "australia",
    "austria",
    "azerbaijan",
    "bahamas",
    "bahrain",
    "baker island",
    "bangladesh",
    "barbados",
    "belarus",
    "belgium",
    "belize",
    "benin",
    "bermuda",
    "bhutan",
    "bolivia",
    "bonaire, saint eustachius and saba",
    "bosnia and herzegovina",
    "botswana",
    "bouvet island",
    "brazil",
    "british indian ocean territory",
    "british virgin islands",
    "brunei darussalam",
    "bulgaria",
    "burkina faso",
    "burundi",
    "cambodia",
    "cameroon",
    "canada",
    "cape verde",
    "cayman islands",
    "central african republic",
    "chad",
    "chile",
    "china",
    "christmas island",
    "cocos (keeling) islands",
    "colombia",
    "comoros",
    "cook islands",
    "costa rica",
    "croatia",
    "cuba",
    "curacao",
    "cyprus",
    "czech republic",
    "democratic republic of congo",
    "denmark",
    "djibouti",
    "dominica",
    "dominican republic",
    "ecuador",
    "egypt",
    "el salvador",
    "equatorial guinea",
    "eritrea",
    "estonia",
    "ethiopia",
    "falkland islands",
    "faroe islands",
    "federated states of micronesia",
    "fiji",
    "finland",
    "france",
    "french guiana",
    "french polynesia",
    "french southern and antarctic lands",
    "gabon",
    "gambia",
    "georgia",
    "germany",
    "ghana",
    "gibraltar",
    "glorioso islands",
    "greece",
    "greenland",
    "grenada",
    "guadeloupe",
    "guam",
    "guatemala",
    "guernsey",
    "guinea",
    "guinea-bissau",
    "guyana",
    "haiti",
    "heard island and mcdonald islands",
    "honduras",
    "hong kong",
    "howland island",
    "hungary",
    "iceland",
    "india",
    "indonesia",
    "iran",
    "iraq",
    "ireland",
    "isle of man",
    "israel",
    "italy",
    "ivory coast",
    "jamaica",
    "japan",
    "jarvis island",
    "jersey",
    "johnston atoll",
    "jordan",
    "juan de nova island",
    "kazakhstan",
    "kenya",
    "kiribati",
    "kosovo",
    "kuwait",
    "kyrgyzstan",
    "laos",
    "latvia",
    "lebanon",
    "lesotho",
    "liberia",
    "libya",
    "liechtenstein",
    "lithuania",
    "luxembourg",
    "macau",
    "macedonia",
    "madagascar",
    "malawi",
    "malaysia",
    "maldives",
    "mali",
    "malta",
    "marshall islands",
    "martinique",
    "mauritania",
    "mauritius",
    "mayotte",
    "mexico",
    "midway islands",
    "moldova",
    "monaco",
    "mongolia",
    "montenegro",
    "montserrat",
    "morocco",
    "mozambique",
    "myanmar",
    "namibia",
    "nauru",
    "nepal",
    "netherlands",
    "new caledonia",
    "new zealand",
    "nicaragua",
    "niger",
    "nigeria",
    "niue",
    "norfolk island",
    "north korea",
    "northern mariana islands",
    "norway",
    "oman",
    "pakistan",
    "palau",
    "palestine",
    "panama",
    "papua new guinea",
    "paraguay",
    "peru",
    "philippines",
    "pitcairn islands",
    "poland",
    "portugal",
    "puerto rico",
    "qatar",
    "republic of congo",
    "reunion",
    "romania",
    "russia",
    "rwanda",
    "saint barthelemy",
    "saint helena",
    "saint kitts and nevis",
    "saint lucia",
    "saint martin",
    "saint pierre and miquelon",
    "saint vincent and the grenadines",
    "samoa",
    "san marino",
    "sao tome and principe",
    "saudi arabia",
    "senegal",
    "serbia",
    "seychelles",
    "sierra leone",
    "singapore",
    "slovakia",
    "slovenia",
    "solomon islands",
    "somalia",
    "south africa",
    "south georgia and the south sandwich islands",
    "south korea",
    "south sudan",
    "spain",
    "sri lanka",
    "sudan",
    "suriname",
    "svalbard and jan mayen",
    "swaziland",
    "sweden",
    "switzerland",
    "syria",
    "taiwan",
    "tajikistan",
    "tanzania",
    "thailand",
    "timor-leste",
    "togo",
    "tokelau",
    "tonga",
    "trinidad and tobago",
    "tunisia",
    "turkey",
    "turkmenistan",
    "turks and caicos islands",
    "tuvalu",
    "us virgin islands",
    "uganda",
    "ukraine",
    "united arab emirates",
    "united kingdom",
    "united states",
    "uruguay",
    "uzbekistan",
    "vanuatu",
    "vatican city",
    "venezuela",
    "vietnam",
    "wake island",
    "wallis and futuna",
    "western sahara",
    "yemen",
    "zambia",
    "zimbabwe",
}

var alternativeNamings = [...]string{
    "curaçao",
    "côte d'ivoire",
    "cote d'ivoire",
    "lao people's democratic republic",
    "palestinian territories",
    "south georgia and south sandwich islands",
    "macao",
}

var countriesMap = map[string]string{
    "afghanistan": "Afghanistan",
    "aland islands": "Aland Islands",
    "albania": "Albania",
    "algeria": "Algeria",
    "american samoa": "American Samoa",
    "andorra": "Andorra",
    "angola": "Angola",
    "anguilla": "Anguilla",
    "antigua and barbuda": "Antigua and Barbuda",
    "argentina": "Argentina",
    "armenia": "Armenia",
    "aruba": "Aruba",
    "australia": "Australia",
    "austria": "Austria",
    "azerbaijan": "Azerbaijan",
    "bahamas": "Bahamas",
    "bahrain": "Bahrain",
    "baker island": "Baker Island",
    "bangladesh": "Bangladesh",
    "barbados": "Barbados",
    "belarus": "Belarus",
    "belgium": "Belgium",
    "belize": "Belize",
    "benin": "Benin",
    "bermuda": "Bermuda",
    "bhutan": "Bhutan",
    "bolivia": "Bolivia",
    "bonaire, saint eustachius and saba": "Bonaire, Saint Eustachius and Saba",
    "bosnia and herzegovina": "Bosnia and Herzegovina",
    "botswana": "Botswana",
    "bouvet island": "Bouvet Island",
    "brazil": "Brazil",
    "british indian ocean territory": "British Indian Ocean Territory",
    "british virgin islands": "British Virgin Islands",
    "brunei darussalam": "Brunei Darussalam",
    "bulgaria": "Bulgaria",
    "burkina faso": "Burkina Faso",
    "burundi": "Burundi",
    "cambodia": "Cambodia",
    "cameroon": "Cameroon",
    "canada": "Canada",
    "cape verde": "Cape Verde",
    "cayman islands": "Cayman Islands",
    "central african republic": "Central African Republic",
    "chad": "Chad",
    "chile": "Chile",
    "china": "China",
    "christmas island": "Christmas Island",
    "cocos (keeling) islands": "Cocos (Keeling) Islands",
    "colombia": "Colombia",
    "comoros": "Comoros",
    "cook islands": "Cook Islands",
    "costa rica": "Costa Rica",
    "croatia": "Croatia",
    "cuba": "Cuba",
    "curaçao": "Curaçao",
    "curacao": "Curaçao",
    "cyprus": "Cyprus",
    "czech republic": "Czech Republic",
    "côte d'ivoire": "Côte d'Ivoire",
    "cote d'ivoire": "Côte d'Ivoire",
    "ivory coast": "Côte d'Ivoire",
    "democratic republic of congo": "Democratic Republic of Congo",
    "denmark": "Denmark",
    "djibouti": "Djibouti",
    "dominica": "Dominica",
    "dominican republic": "Dominican Republic",
    "ecuador": "Ecuador",
    "egypt": "Egypt",
    "el salvador": "El Salvador",
    "equatorial guinea": "Equatorial Guinea",
    "eritrea": "Eritrea",
    "estonia": "Estonia",
    "ethiopia": "Ethiopia",
    "falkland islands": "Falkland Islands",
    "faroe islands": "Faroe Islands",
    "federated states of micronesia": "Federated States of Micronesia",
    "fiji": "Fiji",
    "finland": "Finland",
    "france": "France",
    "french guiana": "French Guiana",
    "french polynesia": "French Polynesia",
    "french southern and antarctic lands": "French Southern and Antarctic Lands",
    "gabon": "Gabon",
    "gambia": "Gambia",
    "georgia": "Georgia",
    "germany": "Germany",
    "ghana": "Ghana",
    "gibraltar": "Gibraltar",
    "glorioso islands": "Glorioso Islands",
    "greece": "Greece",
    "greenland": "Greenland",
    "grenada": "Grenada",
    "guadeloupe": "Guadeloupe",
    "guam": "Guam",
    "guatemala": "Guatemala",
    "guernsey": "Guernsey",
    "guinea": "Guinea",
    "guinea-bissau": "Guinea-Bissau",
    "guyana": "Guyana",
    "haiti": "Haiti",
    "heard island and mcdonald islands": "Heard Island and McDonald Islands",
    "honduras": "Honduras",
    "hong kong": "Hong Kong",
    "howland island": "Howland Island",
    "hungary": "Hungary",
    "iceland": "Iceland",
    "india": "India",
    "indonesia": "Indonesia",
    "iran": "Iran",
    "iraq": "Iraq",
    "ireland": "Ireland",
    "isle of man": "Isle of Man",
    "israel": "Israel",
    "italy": "Italy",
    "jamaica": "Jamaica",
    "japan": "Japan",
    "jarvis island": "Jarvis Island",
    "jersey": "Jersey",
    "johnston atoll": "Johnston Atoll",
    "jordan": "Jordan",
    "juan de nova island": "Juan De Nova Island",
    "kazakhstan": "Kazakhstan",
    "kenya": "Kenya",
    "kiribati": "Kiribati",
    "kosovo": "Kosovo",
    "kuwait": "Kuwait",
    "kyrgyzstan": "Kyrgyzstan",
    "laos": "Lao People's Democratic Republic",
    "lao people's democratic republic": "Lao People's Democratic Republic",
    "latvia": "Latvia",
    "lebanon": "Lebanon",
    "lesotho": "Lesotho",
    "liberia": "Liberia",
    "libya": "Libya",
    "liechtenstein": "Liechtenstein",
    "lithuania": "Lithuania",
    "luxembourg": "Luxembourg",
    "macau": "Macau",
    "macao": "Macau",
    "macedonia": "Macedonia",
    "madagascar": "Madagascar",
    "malawi": "Malawi",
    "malaysia": "Malaysia",
    "maldives": "Maldives",
    "mali": "Mali",
    "malta": "Malta",
    "marshall islands": "Marshall Islands",
    "martinique": "Martinique",
    "mauritania": "Mauritania",
    "mauritius": "Mauritius",
    "mayotte": "Mayotte",
    "mexico": "Mexico",
    "midway islands": "Midway Islands",
    "moldova": "Moldova",
    "monaco": "Monaco",
    "mongolia": "Mongolia",
    "montenegro": "Montenegro",
    "montserrat": "Montserrat",
    "morocco": "Morocco",
    "mozambique": "Mozambique",
    "myanmar": "Myanmar",
    "namibia": "Namibia",
    "nauru": "Nauru",
    "nepal": "Nepal",
    "netherlands": "Netherlands",
    "new caledonia": "New Caledonia",
    "new zealand": "New Zealand",
    "nicaragua": "Nicaragua",
    "niger": "Niger",
    "nigeria": "Nigeria",
    "niue": "Niue",
    "norfolk island": "Norfolk Island",
    "north korea": "North Korea",
    "northern mariana islands": "Northern Mariana Islands",
    "norway": "Norway",
    "oman": "Oman",
    "pakistan": "Pakistan",
    "palau": "Palau",
    "palestine": "Palestinian Territories",
    "palestinian territories": "Palestinian Territories",
    "panama": "Panama",
    "papua new guinea": "Papua New Guinea",
    "paraguay": "Paraguay",
    "peru": "Peru",
    "philippines": "Philippines",
    "pitcairn islands": "Pitcairn Islands",
    "poland": "Poland",
    "portugal": "Portugal",
    "puerto rico": "Puerto Rico",
    "qatar": "Qatar",
    "republic of congo": "Republic of Congo",
    "reunion": "Reunion",
    "romania": "Romania",
    "russia": "Russia",
    "rwanda": "Rwanda",
    "saint barthelemy": "Saint Barthelemy",
    "saint helena": "Saint Helena",
    "saint kitts and nevis": "Saint Kitts and Nevis",
    "saint lucia": "Saint Lucia",
    "saint martin": "Saint Martin",
    "saint pierre and miquelon": "Saint Pierre and Miquelon",
    "saint vincent and the grenadines": "Saint Vincent and the Grenadines",
    "samoa": "Samoa",
    "san marino": "San Marino",
    "sao tome and principe": "Sao Tome and Principe",
    "saudi arabia": "Saudi Arabia",
    "senegal": "Senegal",
    "serbia": "Serbia",
    "seychelles": "Seychelles",
    "sierra leone": "Sierra Leone",
    "singapore": "Singapore",
    "slovakia": "Slovakia",
    "slovenia": "Slovenia",
    "solomon islands": "Solomon Islands",
    "somalia": "Somalia",
    "south africa": "South Africa",
    "south georgia and the south sandwich islands": "South Georgia and South Sandwich Islands",
    "south georgia and south sandwich islands": "South Georgia and South Sandwich Islands",
    "south korea": "South Korea",
    "south sudan": "South Sudan",
    "spain": "Spain",
    "sri lanka": "Sri Lanka",
    "sudan": "Sudan",
    "suriname": "Suriname",
    "svalbard and jan mayen": "Svalbard and Jan Mayen",
    "swaziland": "Swaziland",
    "sweden": "Sweden",
    "switzerland": "Switzerland",
    "syria": "Syria",
    "taiwan": "Taiwan",
    "tajikistan": "Tajikistan",
    "tanzania": "Tanzania",
    "thailand": "Thailand",
    "timor-leste": "Timor-Leste",
    "togo": "Togo",
    "tokelau": "Tokelau",
    "tonga": "Tonga",
    "trinidad and tobago": "Trinidad and Tobago",
    "tunisia": "Tunisia",
    "turkey": "Turkey",
    "turkmenistan": "Turkmenistan",
    "turks and caicos islands": "Turks and Caicos Islands",
    "tuvalu": "Tuvalu",
    "us virgin islands": "US Virgin Islands",
    "uganda": "Uganda",
    "ukraine": "Ukraine",
    "united arab emirates": "United Arab Emirates",
    "united kingdom": "United Kingdom",
    "united states": "United States",
    "uruguay": "Uruguay",
    "uzbekistan": "Uzbekistan",
    "vanuatu": "Vanuatu",
    "vatican city": "Vatican City",
    "venezuela": "Venezuela",
    "vietnam": "Vietnam",
    "wake island": "Wake Island",
    "wallis and futuna": "Wallis and Futuna",
    "western sahara": "Western Sahara",
    "yemen": "Yemen",
    "zambia": "Zambia",
    "zimbabwe": "Zimbabwe",
}

func getCountries(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(countries)
}

func getAlternativeNamings(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(alternativeNamings)
}

func getCountriesMap(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(countriesMap)
}

var codes = map[string]string{
    "Afghanistan": "af",
    "Aland Islands": "ax",
    "Albania": "al",
    "Algeria": "dz",
    "American Samoa": "as",
    "Andorra": "ad",
    "Angola": "ao",
    "Anguilla": "ai",
    "Antigua and Barbuda": "ag",
    "Argentina": "ar",
    "Armenia": "am",
    "Aruba": "aw",
    "Australia": "au",
    "Austria": "at",
    "Azerbaijan": "az",
    "Bahamas": "bs",
    "Bahrain": "bh",
    "Baker Island": "us",
    "Bangladesh": "bd",
    "Barbados": "bb",
    "Belarus": "by",
    "Belgium": "be",
    "Belize": "bz",
    "Benin": "bj",
    "Bermuda": "bm",
    "Bhutan": "bt",
    "Bolivia": "bo",
    "Bonaire, Saint Eustachius and Saba": "bq",
    "Bosnia and Herzegovina": "ba",
    "Botswana": "bw",
    "Bouvet Island": "bv",
    "Brazil": "br",
    "British Indian Ocean Territory": "io",
    "British Virgin Islands": "vg",
    "Brunei Darussalam": "bn",
    "Bulgaria": "bg",
    "Burkina Faso": "bf",
    "Burundi": "bi",
    "Cambodia": "kh",
    "Cameroon": "cm",
    "Canada": "ca",
    "Cape Verde": "cv",
    "Cayman Islands": "ky",
    "Central African Republic": "cf",
    "Chad": "td",
    "Chile": "cl",
    "China": "cn",
    "Christmas Island": "cx",
    "Cocos (Keeling) Islands": "cc",
    "Colombia": "co",
    "Comoros": "km",
    "Cook Islands": "ck",
    "Costa Rica": "cr",
    "Croatia": "hr",
    "Cuba": "cu",
    "Curaçao": "cw",
    "Cyprus": "cy",
    "Czech Republic": "cz",
    "Côte d'Ivoire": "ci",
    "Democratic Republic of Congo": "cd",
    "Denmark": "dk",
    "Djibouti": "dj",
    "Dominica": "dm",
    "Dominican Republic": "do",
    "Ecuador": "ec",
    "Egypt": "eg",
    "El Salvador": "sv",
    "Equatorial Guinea": "gq",
    "Eritrea": "er",
    "Estonia": "ee",
    "Ethiopia": "et",
    "Falkland Islands": "fk",
    "Faroe Islands": "fo",
    "Federated States of Micronesia": "fm",
    "Fiji": "fj",
    "Finland": "fi",
    "France": "fr",
    "French Guiana": "gf",
    "French Polynesia": "pf",
    "French Southern and Antarctic Lands": "tf",
    "Gabon": "ga",
    "Gambia": "gm",
    "Georgia": "ge",
    "Germany": "de",
    "Ghana": "gh",
    "Gibraltar": "gi",
    "Glorioso Islands": "tf",
    "Greece": "gr",
    "Greenland": "gl",
    "Grenada": "gd",
    "Guadeloupe": "gp",
    "Guam": "gu",
    "Guatemala": "gt",
    "Guernsey": "gg",
    "Guinea": "gn",
    "Guinea-Bissau": "gw",
    "Guyana": "gy",
    "Haiti": "ht",
    "Heard Island and McDonald Islands": "hm",
    "Honduras": "hn",
    "Hong Kong": "hk",
    "Howland Island": "us",
    "Hungary": "hu",
    "Iceland": "is",
    "India": "in",
    "Indonesia": "id",
    "Iran": "ir",
    "Iraq": "iq",
    "Ireland": "ie",
    "Isle of Man": "im",
    "Israel": "il",
    "Italy": "it",
    "Jamaica": "jm",
    "Japan": "jp",
    "Jarvis Island": "us",
    "Jersey": "je",
    "Johnston Atoll": "us",
    "Jordan": "jo",
    "Juan De Nova Island": "tf",
    "Kazakhstan": "kz",
    "Kenya": "ke",
    "Kiribati": "ki",
    "Kosovo": "xk",
    "Kuwait": "kw",
    "Kyrgyzstan": "kg",
    "Lao People's Democratic Republic": "la",
    "Latvia": "lv",
    "Lebanon": "lb",
    "Lesotho": "ls",
    "Liberia": "lr",
    "Libya": "ly",
    "Liechtenstein": "li",
    "Lithuania": "lt",
    "Luxembourg": "lu",
    "Macau": "mo",
    "Macedonia": "mk",
    "Madagascar": "mg",
    "Malawi": "mw",
    "Malaysia": "my",
    "Maldives": "mv",
    "Mali": "ml",
    "Malta": "mt",
    "Marshall Islands": "mh",
    "Martinique": "mq",
    "Mauritania": "mr",
    "Mauritius": "mu",
    "Mayotte": "yt",
    "Mexico": "mx",
    "Midway Islands": "us",
    "Moldova": "md",
    "Monaco": "mc",
    "Mongolia": "mn",
    "Montenegro": "me",
    "Montserrat": "ms",
    "Morocco": "ma",
    "Mozambique": "mz",
    "Myanmar": "mm",
    "Namibia": "na",
    "Nauru": "nr",
    "Nepal": "np",
    "Netherlands": "nl",
    "New Caledonia": "nc",
    "New Zealand": "nz",
    "Nicaragua": "ni",
    "Niger": "ne",
    "Nigeria": "ng",
    "Niue": "nu",
    "Norfolk Island": "nf",
    "North Korea": "kp",
    "Northern Mariana Islands": "mp",
    "Norway": "no",
    "Oman": "om",
    "Pakistan": "pk",
    "Palau": "pw",
    "Palestinian Territories": "ps",
    "Panama": "pa",
    "Papua New Guinea": "pg",
    "Paraguay": "py",
    "Peru": "pe",
    "Philippines": "ph",
    "Pitcairn Islands": "pn",
    "Poland": "pl",
    "Portugal": "pt",
    "Puerto Rico": "pr",
    "Qatar": "qa",
    "Republic of Congo": "cg",
    "Reunion": "re",
    "Romania": "ro",
    "Russia": "ru",
    "Rwanda": "rw",
    "Saint Barthelemy": "bl",
    "Saint Helena": "sh",
    "Saint Kitts and Nevis": "kn",
    "Saint Lucia": "lc",
    "Saint Martin": "mf",
    "Saint Pierre and Miquelon": "pm",
    "Saint Vincent and the Grenadines": "vc",
    "Samoa": "ws",
    "San Marino": "sm",
    "Sao Tome and Principe": "st",
    "Saudi Arabia": "sa",
    "Senegal": "sn",
    "Serbia": "rs",
    "Seychelles": "sc",
    "Sierra Leone": "sl",
    "Singapore": "sg",
    "Slovakia": "sk",
    "Slovenia": "si",
    "Solomon Islands": "sb",
    "Somalia": "so",
    "South Africa": "za",
    "South Georgia and South Sandwich Islands": "gs",
    "South Korea": "kr",
    "South Sudan": "ss",
    "Spain": "es",
    "Sri Lanka": "lk",
    "Sudan": "sd",
    "Suriname": "sr",
    "Svalbard and Jan Mayen": "sj",
    "Swaziland": "sz",
    "Sweden": "se",
    "Switzerland": "ch",
    "Syria": "sy",
    "Taiwan": "tw",
    "Tajikistan": "tj",
    "Tanzania": "tz",
    "Thailand": "th",
    "Timor-Leste": "tl",
    "Togo": "tg",
    "Tokelau": "tk",
    "Tonga": "to",
    "Trinidad and Tobago": "tt",
    "Tunisia": "tn",
    "Turkey": "tr",
    "Turkmenistan": "tm",
    "Turks and Caicos Islands": "tc",
    "Tuvalu": "tv",
    "US Virgin Islands": "vi",
    "Uganda": "ug",
    "Ukraine": "ua",
    "United Arab Emirates": "ae",
    "United Kingdom": "gb",
    "United States": "us",
    "Uruguay": "uy",
    "Uzbekistan": "uz",
    "Vanuatu": "vu",
    "Vatican City": "va",
    "Venezuela": "ve",
    "Vietnam": "vn",
    "Wake Island": "us",
    "Wallis and Futuna": "wf",
    "Western Sahara": "eh",
    "Yemen": "ye",
    "Zambia": "zm",
    "Zimbabwe": "zw",
}

func getCodes(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(codes)
}

