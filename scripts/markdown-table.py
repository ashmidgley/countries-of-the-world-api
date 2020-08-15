#!/usr/bin/env python3

def codes_to_dictionary(codes):
    dict = {};
    for x in codes:
        line = x.strip().split('"')
        key = line[3]
        value = line[1]
        dict[key] = value
    return dict;

if __name__ == "__main__":
    filename = "codes.json"
    with open(filename) as f:
        codes = f.readlines()
    codes = codes[1:-1]
    dict = codes_to_dictionary(codes);

    filename = "countries.txt"
    with open(filename) as f:
        countries = f.readlines()
    list = [x.strip().lower() for x in countries] 
    list.sort()

    i = 0
    for (key, value) in sorted(dict.items()):
        print(value + " | " + key + " | "+ list[i]) 
        i+= 1
