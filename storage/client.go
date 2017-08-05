package storage

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type StarWarsDBClient interface {
	GetCharacters() []Character
	GetStarships() []Starship
	GetPlanets() []Planet
	GetCharacter(id string) (Character, bool)
	GetStarship(id string) (Starship, bool)
	GetPlanet(id string) (Planet, bool)
}

type dynamoDBClient struct {
	client         *dynamodb.DynamoDB
	characterCache []Character
	starshipCache  []Starship
	planetCache    []Planet
}

func NewStarWarsDBClient() StarWarsDBClient {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	db := &dynamoDBClient{}
	db.client = svc
	return db
}

func (db *dynamoDBClient) GetCharacters() []Character {
	if db.characterCache == nil {
		db.characterCache = scanCharacters(db.client)
	}
	return db.characterCache
}

func (db *dynamoDBClient) GetStarships() []Starship {
	if db.starshipCache == nil {
		db.starshipCache = scanStarships(db.client)
	}
	return db.starshipCache
}

func (db *dynamoDBClient) GetPlanets() []Planet {
	if db.planetCache == nil {
		db.planetCache = scanPlanets(db.client)
	}
	return db.planetCache
}

func (db *dynamoDBClient) GetCharacter(id string) (Character, bool) {
	characters := db.GetCharacters()
	for i := 0; i < len(characters); i++ {
		if characters[i].Id == id {
			return characters[i], true
		}
	}
	return Character{}, false
}

func (db *dynamoDBClient) GetStarship(id string) (Starship, bool) {
	starships := db.GetStarships()
	for i := 0; i < len(starships); i++ {
		if starships[i].Id == id {
			return starships[i], true
		}
	}
	return Starship{}, false
}

func (db *dynamoDBClient) GetPlanet(id string) (Planet, bool) {
	planets := db.GetPlanets()
	for i := 0; i < len(planets); i++ {
		if planets[i].Id == id {
			return planets[i], true
		}
	}
	return Planet{}, false
}

func getScanInput(recordType string) *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		ExpressionAttributeNames: map[string]*string{
			"#T": aws.String("type"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(recordType),
			},
		},
		FilterExpression: aws.String("#T = :t"),
		TableName:        aws.String("StarWars"),
	}
}

func scanCharacters(svc *dynamodb.DynamoDB) (records []Character) {
	// Use the ScanPages method to perform the scan with pagination. Use
	// just Scan method to make the API call without pagination.
	svc.ScanPages(getScanInput("character"), func(page *dynamodb.ScanOutput, last bool) bool {
		recs := []Character{}

		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &recs)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
		}

		records = append(records, recs...)

		return true // keep paging
	})
	return
}

func scanStarships(svc *dynamodb.DynamoDB) (records []Starship) {
	// Use the ScanPages method to perform the scan with pagination. Use
	// just Scan method to make the API call without pagination.
	svc.ScanPages(getScanInput("starship"), func(page *dynamodb.ScanOutput, last bool) bool {
		recs := []Starship{}

		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &recs)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
		}

		records = append(records, recs...)

		return true // keep paging
	})
	return
}

func scanPlanets(svc *dynamodb.DynamoDB) (records []Planet) {
	// Use the ScanPages method to perform the scan with pagination. Use
	// just Scan method to make the API call without pagination.
	svc.ScanPages(getScanInput("planet"), func(page *dynamodb.ScanOutput, last bool) bool {
		recs := []Planet{}

		err := dynamodbattribute.UnmarshalListOfMaps(page.Items, &recs)
		if err != nil {
			panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
		}

		records = append(records, recs...)

		return true // keep paging
	})
	return
}
