#!/usr/bin/env python3

def countries_to_list(countries):
    return [x.strip() for x in countries] 
               
def countries_to_dictionary(countries):
    dict = {};
    for x in countries:
        key = x.lower()
        dict[key] = x
    return dict

if __name__ == "__main__":
    filename = "countries.txt"
    with open(filename) as f:
        countries = f.readlines()
    country_list = countries_to_list(countries)
    country_list.sort()
    for x in country_list:
        print('"' + x.lower() + '",')
    country_dict = countries_to_dictionary(country_list)
    for key in country_dict:
        print('"' + key + '": "' + country_dict[key] + '",')
