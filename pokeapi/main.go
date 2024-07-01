package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

/*
make http request to the PokeAPI to get location data
*/
func GetLocationData(url string) (Location, []byte, error) {
	location := Location{}
	rawData := []byte{}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return location, rawData, err
	}
	rawData, err = io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatal("Some sort of HTTP error")
		return location, rawData, err
	}
	if err != nil {
		log.Fatal(err)
		return location, rawData, err
	}
	err = json.Unmarshal(rawData, &location)
	if err != nil {
		log.Fatal(err)
		return location, rawData, err
	}
	return location, rawData, nil
}
