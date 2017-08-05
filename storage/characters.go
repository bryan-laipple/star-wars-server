package storage

import (
	"github.com/leejarvis/swapi"
)

type charactersResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []Character `json:"results"`
}

func GetCharacters() (people []Character, err error) {
	url := "https://swapi.co/api/people/"
	for url != "" {
		var res charactersResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		for _, person := range res.Results {
			people = append(people, person)
		}
	}
	return
}
