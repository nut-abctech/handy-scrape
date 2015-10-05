package parser

import (
    "encoding/json"
    "io/ioutil"
)

type attributes struct {
    Link      string `json:"link"`
    Name      string `json:"name"`
    Detail    string `json:"detail"`
    ContactNo string `json:"contactNo"`
    Location  string `json:"location"`
}
type routes struct {
    Entry string     `json:"entry"`
    Info  attributes `json:"attributes"`
}

// Role structure
type Role struct {
    Title  string `json:"title"`
    URL    string `json:"url"`
    Get    string `json:"getURL"`
    Routes routes `json:"routes"`
}

// Parse the json file and return &Role
func Parse(filename string) (*Role, error) {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    var roletype Role
    if err := json.Unmarshal(bytes, &roletype); err != nil {
        return nil, err
    }
    return &roletype, nil
}
