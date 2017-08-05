package etl

import (
	"github.com/leejarvis/swapi"
)

type Person struct {
	DynamoDBKey
	swapi.Person
}

type PeopleResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Person `json:"results"`
}

func GetPeople() (people []Person, err error) {
	url := "https://swapi.co/api/people/"
	for url != "" {
		var res PeopleResponse
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
