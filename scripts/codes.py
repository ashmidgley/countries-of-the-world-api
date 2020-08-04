#!/usr/bin/env python3

def codes_to_dictionary(codes):
    dict = {};
    for x in codes:
        line = x.strip().split('"')
        key = line[3]
        value = line[1].lower()
        dict[key] = value
    return dict;

def print_codes(codes):
    for x in codes:
        line = x.strip().split('"')
        key = line[3]
        value = line[1].lower()
        print('"' + key + '": "' + value + '",') 

if __name__ == "__main__":
    filename = "codes.json"
    with open(filename) as f:
        codes = f.readlines()
    codes = codes[1:-1]
    dict = codes_to_dictionary(codes);
    #print_codes(codes)
