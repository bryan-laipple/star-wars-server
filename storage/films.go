package storage

import (
	"github.com/leejarvis/swapi"
)

type filmsResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Film `json:"results"`
}

func GetFilms() (films []Film, err error) {
	url := "https://swapi.co/api/films/"
	for url != "" {
		var res filmsResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		for _, film := range res.Results {
			films = append(films, film)
		}
	}
	return
}
