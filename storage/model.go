package storage

import (
	"github.com/leejarvis/swapi"
)

type dynamoDBKey struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Character struct {
	dynamoDBKey
	swapi.Person
	Avatar string   `json:"avatar"`
	Image  string   `json:"image"`
	Links  []Link   `json:"links"`
	Images []string `json:"images"`
}

type Film struct {
	dynamoDBKey
	swapi.Film
}

type Planet struct {
	dynamoDBKey
	swapi.Planet
	Avatar string   `json:"avatar"`
	Image  string   `json:"image"`
	Links  []Link   `json:"links"`
	Images []string `json:"images"`
}

type Species struct {
	dynamoDBKey
	swapi.Species
}

type Starship struct {
	dynamoDBKey
	swapi.Starship
	Avatar string   `json:"avatar"`
	Image  string   `json:"image"`
	Links  []Link   `json:"links"`
	Images []string `json:"images"`
}
