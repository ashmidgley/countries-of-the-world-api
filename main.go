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

    statement := "CREATE TABLE IF NOT EXISTS leaderboard (id INTEGER PRIMARY KEY, name TEXT, country TEXT, countries INTEGER, time TEXT)"
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
    Time string `json:"time"`
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
        var time string
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

    statement, err := database.Prepare("DELETE FROM leaderboard WHERE id = ?")
	if err != nil {
        writer.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(writer, err.Error())
        return
	}
    defer statement.Close()

    statement.Exec(id)
}

var countries = [...]string{
    "andorra",
    "united arab emirates",
    "afghanistan",
    "antigua and barbuda",
    "anguilla",
    "albania",
    "armenia",
    "angola",
    "argentina",
    "american samoa",
    "austria",
    "australia",
    "aruba",
    "aland islands",
    "azerbaijan",
    "bosnia and herzegovina",
    "barbados",
    "bangladesh",
    "belgium",
    "burkina faso",
    "bulgaria",
    "bahrain",
    "burundi",
    "benin",
    "saint barthelemy",
    "brunei darussalam",
    "bolivia",
    "bermuda",
    "bonaire, saint eustachius and saba",
    "brazil",
    "bahamas",
    "bhutan",
    "bouvet island",
    "botswana",
    "belarus",
    "belize",
    "canada",
    "cocos (keeling) islands",
    "democratic republic of congo",
    "central african republic",
    "republic of congo",
    "switzerland",
    "côte d'ivoire",
    "cook islands",
    "chile",
    "cameroon",
    "china",
    "colombia",
    "costa rica",
    "cuba",
    "cape verde",
    "curaçao",
    "christmas island",
    "cyprus",
    "czech republic",
    "germany",
    "djibouti",
    "denmark",
    "dominica",
    "dominican republic",
    "algeria",
    "ecuador",
    "egypt",
    "estonia",
    "western sahara",
    "eritrea",
    "spain",
    "ethiopia",
    "finland",
    "fiji",
    "falkland islands",
    "federated states of micronesia",
    "faroe islands",
    "france",
    "gabon",
    "united kingdom",
    "georgia",
    "grenada",
    "french guiana",
    "guernsey",
    "ghana",
    "gibraltar",
    "greenland",
    "gambia",
    "guinea",
    "glorioso islands",
    "guadeloupe",
    "equatorial guinea",
    "greece",
    "south georgia and south sandwich islands",
    "guatemala",
    "guam",
    "guinea-bissau",
    "guyana",
    "hong kong",
    "heard island and mcdonald islands",
    "honduras",
    "croatia",
    "haiti",
    "hungary",
    "indonesia",
    "ireland",
    "israel",
    "isle of man",
    "india",
    "british indian ocean territory",
    "iraq",
    "iran",
    "iceland",
    "italy",
    "jersey",
    "jamaica",
    "jordan",
    "japan",
    "juan de nova island",
    "kenya",
    "kyrgyzstan",
    "cambodia",
    "kiribati",
    "comoros",
    "saint kitts and nevis",
    "north korea",
    "south korea",
    "kosovo",
    "kuwait",
    "cayman islands",
    "kazakhstan",
    "lao people's democratic republic",
    "lebanon",
    "saint lucia",
    "liechtenstein",
    "sri lanka",
    "liberia",
    "lesotho",
    "lithuania",
    "luxembourg",
    "latvia",
    "libya",
    "morocco",
    "monaco",
    "moldova",
    "madagascar",
    "montenegro",
    "saint martin",
    "marshall islands",
    "macedonia",
    "mali",
    "macau",
    "myanmar",
    "mongolia",
    "northern mariana islands",
    "martinique",
    "mauritania",
    "montserrat",
    "malta",
    "mauritius",
    "maldives",
    "malawi",
    "mexico",
    "malaysia",
    "mozambique",
    "namibia",
    "new caledonia",
    "niger",
    "norfolk island",
    "nigeria",
    "nicaragua",
    "netherlands",
    "norway",
    "nepal",
    "nauru",
    "niue",
    "new zealand",
    "oman",
    "panama",
    "peru",
    "french polynesia",
    "papua new guinea",
    "philippines",
    "pakistan",
    "poland",
    "saint pierre and miquelon",
    "pitcairn islands",
    "puerto rico",
    "palestinian territories",
    "portugal",
    "palau",
    "paraguay",
    "qatar",
    "reunion",
    "romania",
    "serbia",
    "russia",
    "rwanda",
    "saudi arabia",
    "solomon islands",
    "seychelles",
    "sudan",
    "sweden",
    "singapore",
    "saint helena",
    "slovenia",
    "svalbard and jan mayen",
    "slovakia",
    "sierra leone",
    "san marino",
    "senegal",
    "somalia",
    "suriname",
    "south sudan",
    "sao tome and principe",
    "el salvador",
    "saint martin",
    "syria",
    "swaziland",
    "turks and caicos islands",
    "chad",
    "french southern and antarctic lands",
    "togo",
    "thailand",
    "tajikistan",
    "tokelau",
    "timor-leste",
    "turkmenistan",
    "tunisia",
    "tonga",
    "turkey",
    "trinidad and tobago",
    "tuvalu",
    "taiwan",
    "tanzania",
    "ukraine",
    "uganda",
    "jarvis island",
    "baker island",
    "howland island",
    "johnston atoll",
    "midway islands",
    "wake island",
    "united states",
    "uruguay",
    "uzbekistan",
    "vatican city",
    "saint vincent and the grenadines",
    "venezuela",
    "british virgin islands",
    "us virgin islands",
    "vietnam",
    "vanuatu",
    "wallis and futuna",
    "samoa",
    "yemen",
    "mayotte",
    "south africa",
    "zambia",
    "zimbabwe",
}

var countriesMap = map[string]string{
    "andorra": "Andorra",
    "united arab emirates": "United Arab Emirates",
    "afghanistan": "Afghanistan",
    "antigua and barbuda": "Antigua and Barbuda",
    "anguilla": "Anguilla",
    "albania": "Albania",
    "armenia": "Armenia",
    "angola": "Angola",
    "argentina": "Argentina",
    "american samoa": "American Samoa",
    "austria": "Austria",
    "australia": "Australia",
    "aruba": "Aruba",
    "aland islands": "Aland Islands",
    "azerbaijan": "Azerbaijan",
    "bosnia and herzegovina": "Bosnia and Herzegovina",
    "barbados": "Barbados",
    "bangladesh": "Bangladesh",
    "belgium": "Belgium",
    "burkina faso": "Burkina Faso",
    "bulgaria": "Bulgaria",
    "bahrain": "Bahrain",
    "burundi": "Burundi",
    "benin": "Benin",
    "saint barthelemy": "Saint Barthelemy",
    "brunei darussalam": "Brunei Darussalam",
    "bolivia": "Bolivia",
    "bermuda": "Bermuda",
    "bonaire, saint eustachius and saba": "Bonaire, Saint Eustachius and Saba",
    "brazil": "Brazil",
    "bahamas": "Bahamas",
    "bhutan": "Bhutan",
    "bouvet island": "Bouvet Island",
    "botswana": "Botswana",
    "belarus": "Belarus",
    "belize": "Belize",
    "canada": "Canada",
    "cocos (keeling) islands": "Cocos (Keeling) Islands",
    "democratic republic of congo": "Democratic Republic of Congo",
    "central african republic": "Central African Republic",
    "republic of congo": "Republic of Congo",
    "switzerland": "Switzerland",
    "côte d'ivoire": "Côte d'Ivoire",
    "cook islands": "Cook Islands",
    "chile": "Chile",
    "cameroon": "Cameroon",
    "china": "China",
    "colombia": "Colombia",
    "costa rica": "Costa Rica",
    "cuba": "Cuba",
    "cape verde": "Cape Verde",
    "curaçao": "Curaçao",
    "christmas island": "Christmas Island",
    "cyprus": "Cyprus",
    "czech republic": "Czech Republic",
    "germany": "Germany",
    "djibouti": "Djibouti",
    "denmark": "Denmark",
    "dominica": "Dominica",
    "dominican republic": "Dominican Republic",
    "algeria": "Algeria",
    "ecuador": "Ecuador",
    "egypt": "Egypt",
    "estonia": "Estonia",
    "western sahara": "Western Sahara",
    "eritrea": "Eritrea",
    "spain": "Spain",
    "ethiopia": "Ethiopia",
    "finland": "Finland",
    "fiji": "Fiji",
    "falkland islands": "Falkland Islands",
    "federated states of micronesia": "Federated States of Micronesia",
    "faroe islands": "Faroe Islands",
    "france": "France",
    "gabon": "Gabon",
    "united kingdom": "United Kingdom",
    "georgia": "Georgia",
    "grenada": "Grenada",
    "french guiana": "French Guiana",
    "guernsey": "Guernsey",
    "ghana": "Ghana",
    "gibraltar": "Gibraltar",
    "greenland": "Greenland",
    "gambia": "Gambia",
    "guinea": "Guinea",
    "glorioso islands": "Glorioso Islands",
    "guadeloupe": "Guadeloupe",
    "equatorial guinea": "Equatorial Guinea",
    "greece": "Greece",
    "south georgia and south sandwich islands": "South Georgia and South Sandwich Islands",
    "guatemala": "Guatemala",
    "guam": "Guam",
    "guinea-bissau": "Guinea-Bissau",
    "guyana": "Guyana",
    "hong kong": "Hong Kong",
    "heard island and mcdonald islands": "Heard Island and McDonald Islands",
    "honduras": "Honduras",
    "croatia": "Croatia",
    "haiti": "Haiti",
    "hungary": "Hungary",
    "indonesia": "Indonesia",
    "ireland": "Ireland",
    "israel": "Israel",
    "isle of man": "Isle of Man",
    "india": "India",
    "british indian ocean territory": "British Indian Ocean Territory",
    "iraq": "Iraq",
    "iran": "Iran",
    "iceland": "Iceland",
    "italy": "Italy",
    "jersey": "Jersey",
    "jamaica": "Jamaica",
    "jordan": "Jordan",
    "japan": "Japan",
    "juan de nova island": "Juan De Nova Island",
    "kenya": "Kenya",
    "kyrgyzstan": "Kyrgyzstan",
    "cambodia": "Cambodia",
    "kiribati": "Kiribati",
    "comoros": "Comoros",
    "saint kitts and nevis": "Saint Kitts and Nevis",
    "north korea": "North Korea",
    "south korea": "South Korea",
    "kosovo": "Kosovo",
    "kuwait": "Kuwait",
    "cayman islands": "Cayman Islands",
    "kazakhstan": "Kazakhstan",
    "lao people's democratic republic": "Lao People's Democratic Republic",
    "lebanon": "Lebanon",
    "saint lucia": "Saint Lucia",
    "liechtenstein": "Liechtenstein",
    "sri lanka": "Sri Lanka",
    "liberia": "Liberia",
    "lesotho": "Lesotho",
    "lithuania": "Lithuania",
    "luxembourg": "Luxembourg",
    "latvia": "Latvia",
    "libya": "Libya",
    "morocco": "Morocco",
    "monaco": "Monaco",
    "moldova": "Moldova",
    "madagascar": "Madagascar",
    "montenegro": "Montenegro",
    "saint martin": "Saint Martin",
    "marshall islands": "Marshall Islands",
    "macedonia": "Macedonia",
    "mali": "Mali",
    "macau": "Macau",
    "myanmar": "Myanmar",
    "mongolia": "Mongolia",
    "northern mariana islands": "Northern Mariana Islands",
    "martinique": "Martinique",
    "mauritania": "Mauritania",
    "montserrat": "Montserrat",
    "malta": "Malta",
    "mauritius": "Mauritius",
    "maldives": "Maldives",
    "malawi": "Malawi",
    "mexico": "Mexico",
    "malaysia": "Malaysia",
    "mozambique": "Mozambique",
    "namibia": "Namibia",
    "new caledonia": "New Caledonia",
    "niger": "Niger",
    "norfolk island": "Norfolk Island",
    "nigeria": "Nigeria",
    "nicaragua": "Nicaragua",
    "netherlands": "Netherlands",
    "norway": "Norway",
    "nepal": "Nepal",
    "nauru": "Nauru",
    "niue": "Niue",
    "new zealand": "New Zealand",
    "oman": "Oman",
    "panama": "Panama",
    "peru": "Peru",
    "french polynesia": "French Polynesia",
    "papua new guinea": "Papua New Guinea",
    "philippines": "Philippines",
    "pakistan": "Pakistan",
    "poland": "Poland",
    "saint pierre and miquelon": "Saint Pierre and Miquelon",
    "pitcairn islands": "Pitcairn Islands",
    "puerto rico": "Puerto Rico",
    "palestinian territories": "Palestinian Territories",
    "portugal": "Portugal",
    "palau": "Palau",
    "paraguay": "Paraguay",
    "qatar": "Qatar",
    "reunion": "Reunion",
    "romania": "Romania",
    "serbia": "Serbia",
    "russia": "Russia",
    "rwanda": "Rwanda",
    "saudi arabia": "Saudi Arabia",
    "solomon islands": "Solomon Islands",
    "seychelles": "Seychelles",
    "sudan": "Sudan",
    "sweden": "Sweden",
    "singapore": "Singapore",
    "saint helena": "Saint Helena",
    "slovenia": "Slovenia",
    "svalbard and jan mayen": "Svalbard and Jan Mayen",
    "slovakia": "Slovakia",
    "sierra leone": "Sierra Leone",
    "san marino": "San Marino",
    "senegal": "Senegal",
    "somalia": "Somalia",
    "suriname": "Suriname",
    "south sudan": "South Sudan",
    "sao tome and principe": "Sao Tome and Principe",
    "el salvador": "El Salvador",
    "syria": "Syria",
    "swaziland": "Swaziland",
    "turks and caicos islands": "Turks and Caicos Islands",
    "chad": "Chad",
    "french southern and antarctic lands": "French Southern and Antarctic Lands",
    "togo": "Togo",
    "thailand": "Thailand",
    "tajikistan": "Tajikistan",
    "tokelau": "Tokelau",
    "timor-leste": "Timor-Leste",
    "turkmenistan": "Turkmenistan",
    "tunisia": "Tunisia",
    "tonga": "Tonga",
    "turkey": "Turkey",
    "trinidad and tobago": "Trinidad and Tobago",
    "tuvalu": "Tuvalu",
    "taiwan": "Taiwan",
    "tanzania": "Tanzania",
    "ukraine": "Ukraine",
    "uganda": "Uganda",
    "jarvis island": "Jarvis Island",
    "baker island": "Baker Island",
    "howland island": "Howland Island",
    "johnston atoll": "Johnston Atoll",
    "midway islands": "Midway Islands",
    "wake island": "Wake Island",
    "united states": "United States",
    "uruguay": "Uruguay",
    "uzbekistan": "Uzbekistan",
    "vatican city": "Vatican City",
    "saint vincent and the grenadines": "Saint Vincent and the Grenadines",
    "venezuela": "Venezuela",
    "british virgin islands": "British Virgin Islands",
    "us virgin islands": "US Virgin Islands",
    "vietnam": "Vietnam",
    "vanuatu": "Vanuatu",
    "wallis and futuna": "Wallis and Futuna",
    "samoa": "Samoa",
    "yemen": "Yemen",
    "mayotte": "Mayotte",
    "south africa": "South Africa",
    "zambia": "Zambia",
    "zimbabwe": "Zimbabwe",
}

func getCountries(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(countries)
}

func getCountriesMap(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(countriesMap)
}

var codes = map[string]string{
    "Andorra": "ad",
    "United Arab Emirates": "ae",
    "Afghanistan": "af",
    "Antigua and Barbuda": "ag",
    "Anguilla": "ai",
    "Albania": "al",
    "Armenia": "am",
    "Netherlands Antilles": "an",
    "Angola": "ao",
    "Antarctica": "aq",
    "Argentina": "ar",
    "American Samoa": "as",
    "Austria": "at",
    "Australia": "au",
    "Aruba": "aw",
    "\u00c5land Islands": "ax",
    "Azerbaijan": "az",
    "Bosnia and Herzegovina": "ba",
    "Barbados": "bb",
    "Bangladesh": "bd",
    "Belgium": "be",
    "Burkina Faso": "bf",
    "Bulgaria": "bg",
    "Bahrain": "bh",
    "Burundi": "bi",
    "Benin": "bj",
    "Saint Barthélemy": "bl",
    "Bermuda": "bm",
    "Brunei Darussalam": "bn",
    "Bolivia, Plurinational State of": "bo",
    "Caribbean Netherlands": "bq",
    "Brazil": "br",
    "Bahamas": "bs",
    "Bhutan": "bt",
    "Bouvet Island": "bv",
    "Botswana": "bw",
    "Belarus": "by",
    "Belize": "bz",
    "Canada": "ca",
    "Cocos (Keeling) Islands": "cc",
    "Congo, the Democratic Republic of the": "cd",
    "Central African Republic": "cf",
    "Congo": "cg",
    "Switzerland": "ch",
    "C\u00f4te d'Ivoire": "ci",
    "Cook Islands": "ck",
    "Chile": "cl",
    "Cameroon": "cm",
    "China": "cn",
    "Colombia": "co",
    "Costa Rica": "cr",
    "Cuba": "cu",
    "Cape Verde": "cv",
    "Cura\u00e7ao": "cw",
    "Christmas Island": "cx",
    "Cyprus": "cy",
    "Czech Republic": "cz",
    "Germany": "de",
    "Djibouti": "dj",
    "Denmark": "dk",
    "Dominica": "dm",
    "Dominican Republic": "do",
    "Algeria": "dz",
    "Ecuador": "ec",
    "Estonia": "ee",
    "Egypt": "eg",
    "Western Sahara": "eh",
    "Eritrea": "er",
    "Spain": "es",
    "Ethiopia": "et",
    "Europe": "eu",
    "Finland": "fi",
    "Fiji": "fj",
    "Falkland Islands (Malvinas)": "fk",
    "Micronesia, Federated States of": "fm",
    "Faroe Islands": "fo",
    "France": "fr",
    "Gabon": "ga",
    "England": "gb-eng",
    "Northern Ireland": "gb-nir",
    "Scotland": "gb-sct",
    "Wales": "gb-wls",
    "United Kingdom": "gb",
    "Grenada": "gd",
    "Georgia": "ge",
    "French Guiana": "gf",
    "Guernsey": "gg",
    "Ghana": "gh",
    "Gibraltar": "gi",
    "Greenland": "gl",
    "Gambia": "gm",
    "Guinea": "gn",
    "Guadeloupe": "gp",
    "Equatorial Guinea": "gq",
    "Greece": "gr",
    "South Georgia and the South Sandwich Islands": "gs",
    "Guatemala": "gt",
    "Guam": "gu",
    "Guinea-Bissau": "gw",
    "Guyana": "gy",
    "Hong Kong": "hk",
    "Heard Island and McDonald Islands": "hm",
    "Honduras": "hn",
    "Croatia": "hr",
    "Haiti": "ht",
    "Hungary": "hu",
    "Indonesia": "id",
    "Ireland": "ie",
    "Israel": "il",
    "Isle of Man": "im",
    "India": "in",
    "British Indian Ocean Territory": "io",
    "Iraq": "iq",
    "Iran, Islamic Republic of": "ir",
    "Iceland": "is",
    "Italy": "it",
    "Jersey": "je",
    "Jamaica": "jm",
    "Jordan": "jo",
    "Japan": "jp",
    "Kenya": "ke",
    "Kyrgyzstan": "kg",
    "Cambodia": "kh",
    "Kiribati": "ki",
    "Comoros": "km",
    "Saint Kitts and Nevis": "kn",
    "Korea, Democratic People's Republic of": "kp",
    "Korea, Republic of": "kr",
    "Kuwait": "kw",
    "Cayman Islands": "ky",
    "Kazakhstan": "kz",
    "Lao People's Democratic Republic": "la",
    "Lebanon": "lb",
    "Saint Lucia": "lc",
    "Liechtenstein": "li",
    "Sri Lanka": "lk",
    "Liberia": "lr",
    "Lesotho": "ls",
    "Lithuania": "lt",
    "Luxembourg": "lu",
    "Latvia": "lv",
    "Libya": "ly",
    "Morocco": "ma",
    "Monaco": "mc",
    "Moldova, Republic of": "md",
    "Montenegro": "me",
    "Saint Martin": "mf",
    "Madagascar": "mg",
    "Marshall Islands": "mh",
    "Macedonia, the former Yugoslav Republic of": "mk",
    "Mali": "ml",
    "Myanmar": "mm",
    "Mongolia": "mn",
    "Macao": "mo",
    "Northern Mariana Islands": "mp",
    "Martinique": "mq",
    "Mauritania": "mr",
    "Montserrat": "ms",
    "Malta": "mt",
    "Mauritius": "mu",
    "Maldives": "mv",
    "Malawi": "mw",
    "Mexico": "mx",
    "Malaysia": "my",
    "Mozambique": "mz",
    "Namibia": "na",
    "New Caledonia": "nc",
    "Niger": "ne",
    "Norfolk Island": "nf",
    "Nigeria": "ng",
    "Nicaragua": "ni",
    "Netherlands": "nl",
    "Norway": "no",
    "Nepal": "np",
    "Nauru": "nr",
    "Niue": "nu",
    "New Zealand": "nz",
    "Oman": "om",
    "Panama": "pa",
    "Peru": "pe",
    "French Polynesia": "pf",
    "Papua New Guinea": "pg",
    "Philippines": "ph",
    "Pakistan": "pk",
    "Poland": "pl",
    "Saint Pierre and Miquelon": "pm",
    "Pitcairn": "pn",
    "Puerto Rico": "pr",
    "Palestine": "ps",
    "Portugal": "pt",
    "Palau": "pw",
    "Paraguay": "py",
    "Qatar": "qa",
    "Réunion": "re",
    "Romania": "ro",
    "Serbia": "rs",
    "Russian Federation": "ru",
    "Rwanda": "rw",
    "Saudi Arabia": "sa",
    "Solomon Islands": "sb",
    "Seychelles": "sc",
    "Sudan": "sd",
    "Sweden": "se",
    "Singapore": "sg",
    "Saint Helena, Ascension and Tristan da Cunha": "sh",
    "Slovenia": "si",
    "Svalbard and Jan Mayen Islands": "sj",
    "Slovakia": "sk",
    "Sierra Leone": "sl",
    "San Marino": "sm",
    "Senegal": "sn",
    "Somalia": "so",
    "Suriname": "sr",
    "South Sudan": "ss",
    "Sao Tome and Principe": "st",
    "El Salvador": "sv",
    "Sint Maarten (Dutch part)": "sx",
    "Syrian Arab Republic": "sy",
    "Swaziland": "sz",
    "Turks and Caicos Islands": "tc",
    "Chad": "td",
    "French Southern Territories": "tf",
    "Togo": "tg",
    "Thailand": "th",
    "Tajikistan": "tj",
    "Tokelau": "tk",
    "Timor-Leste": "tl",
    "Turkmenistan": "tm",
    "Tunisia": "tn",
    "Tonga": "to",
    "Turkey": "tr",
    "Trinidad and Tobago": "tt",
    "Tuvalu": "tv",
    "Taiwan": "tw",
    "Tanzania, United Republic of": "tz",
    "Ukraine": "ua",
    "Uganda": "ug",
    "US Minor Outlying Islands": "um",
    "United States": "us",
    "Uruguay": "uy",
    "Uzbekistan": "uz",
    "Holy See (Vatican City State)": "va",
    "Saint Vincent and the Grenadines": "vc",
    "Venezuela, Bolivarian Republic of": "ve",
    "Virgin Islands, British": "vg",
    "Virgin Islands, U.S.": "vi",
    "Viet Nam": "vn",
    "Vanuatu": "vu",
    "Wallis and Futuna Islands": "wf",
    "Kosovo": "xk",
    "Samoa": "ws",
    "Yemen": "ye",
    "Mayotte": "yt",
    "South Africa": "za",
    "Zambia": "zm",
    "Zimbabwe": "zw",
}

func getCodes(writer http.ResponseWriter, request *http.Request) {
    json.NewEncoder(writer).Encode(codes)
}

