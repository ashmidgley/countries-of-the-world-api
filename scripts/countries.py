#!/usr/bin/env python3

def countries_to_list(countries):
    return [x.strip().lower() for x in countries] 
               
def countries_to_dictionary(countries):
    dict = {};
    for x in countries:
        key = x.strip().lower()
        value = x.strip()
        dict[key] = value
    return dict

def print_countries(countries):
   for x in countries:
       print('"' + x.strip().lower() + '",')

def print_countries_map():
    for x in countries:
        key = x.strip().lower()
        value = x.strip()
        print('"' + key + '": "' + value + '",')

if __name__ == "__main__":
    filename = "countries.txt"
    with open(filename) as f:
        countries = f.readlines()
    country_list = countries_to_list(countries)
    country_dict = countries_to_dictionary(countries)
    #print_countries(countries)
    #print_countries_map()

