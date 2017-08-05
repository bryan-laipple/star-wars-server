package storage

import (
	"github.com/leejarvis/swapi"
)

type speciesResponse struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous string    `json:"previous"`
	Results  []Species `json:"results"`
}

func GetSpecies() (species []Species, err error) {
	url := "https://swapi.co/api/species/"
	for url != "" {
		var res speciesResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		for _, sp := range res.Results {
			species = append(species, sp)
		}
	}
	return
}
