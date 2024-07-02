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

type locationDetail struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
	GameIndices []struct {
		GameIndex  int `json:"game_index"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"game_indices"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}

type Pokemon struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

/*
make http request to the PokeAPI to get location data
*/
func GetLocationData(url string) (Location, []byte, error) {
	location := Location{}
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return location, []byte{}, err
	}
	rawData, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatal("Some sort of HTTP error")
		return location, []byte{}, err
	}
	if err != nil {
		log.Fatal(err)
		return location, []byte{}, err
	}
	err = json.Unmarshal(rawData, &location)
	if err != nil {
		log.Fatal(err)
		return location, []byte{}, err
	}
	return location, rawData, nil
}

func getLocationUrl(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	rawData, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatal("HTTP error")
		return "", err
	}
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var locationDetail locationDetail
	err = json.Unmarshal(rawData, &locationDetail)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return locationDetail.Areas[0].URL, err

}

func GetPokemonInLocation(url string) (Pokemon, []byte, error) {
	var pokemon Pokemon

	pokemonUrl, err := getLocationUrl(url)
	if err != nil {
		log.Fatal(err)
		return pokemon, []byte{}, err
	}

	res, err := http.Get(pokemonUrl)

	if res.StatusCode > 299 {
		log.Fatal("HTTP error")
		return pokemon, []byte{}, nil
	}
	if err != nil {
		log.Fatal(err)
		return pokemon, []byte{}, err
	}

	rawData, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
		return Pokemon{}, []byte{}, err
	}

	err = json.Unmarshal(rawData, &pokemon)
	if err != nil {
		log.Fatal(err)
		return pokemon, []byte{}, err
	}
	return pokemon, rawData, nil

}
