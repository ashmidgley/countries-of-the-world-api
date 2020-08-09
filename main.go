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
    "curaçao",
    "cyprus",
    "czech republic",
    "côte d'ivoire",
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
    "lao people's democratic republic",
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
    "palestinian territories",
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
    "south georgia and south sandwich islands",
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
    "cyprus": "Cyprus",
    "czech republic": "Czech Republic",
    "côte d'ivoire": "Côte d'Ivoire",
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

