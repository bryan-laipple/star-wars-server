package etl

import (
	"github.com/leejarvis/swapi"
)

type Starship struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Key  string `json:"key"`
	swapi.Starship
}

type StarshipsResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Starship `json:"results"`
}

func GetStarships() (starships []Starship, err error) {
	url := "https://swapi.co/api/starships/"
	for url != "" {
		var res StarshipsResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		for _, starship := range res.Results {
			starships = append(starships, starship)
		}
	}
	return
}
