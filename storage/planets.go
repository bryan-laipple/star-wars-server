package storage

import (
	"github.com/leejarvis/swapi"
)

type planetsResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Planet `json:"results"`
}

func GetPlanets() (planets []Planet, err error) {
	url := "https://swapi.co/api/planets/"
	for url != "" {
		var res planetsResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		for _, planet := range res.Results {
			planets = append(planets, planet)
		}
	}
	return
}
