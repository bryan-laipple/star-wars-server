package storage

import (
	"github.com/leejarvis/swapi"
)

type dynamoDBKey struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Character struct {
	dynamoDBKey
	swapi.Person
}

type Film struct {
	dynamoDBKey
	swapi.Film
}

type Planet struct {
	dynamoDBKey
	swapi.Planet
}

type Species struct {
	dynamoDBKey
	swapi.Species
}

type Starship struct {
	dynamoDBKey
	swapi.Starship
}
