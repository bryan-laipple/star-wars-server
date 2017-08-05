package etl

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bryan-laipple/star-wars-server/storage"
	"github.com/leejarvis/swapi"
)

type swapiData struct {
	Characters []storage.Character `json:"characters"`
	Starships  []storage.Starship  `json:"starships"`
	Planets    []storage.Planet    `json:"planets"`
}

const TABLE_NAME = "StarWars"

func createTable(svc *dynamodb.DynamoDB) (success bool) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("type"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("type"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(TABLE_NAME),
	}
	fmt.Printf("Creating DynamoDB table '%s'...\n", TABLE_NAME)
	result, err := svc.CreateTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				fmt.Println(dynamodb.ErrCodeResourceInUseException, aerr.Error())
			case dynamodb.ErrCodeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return false
	}

	fmt.Println(result)
	return true
}

func tableExists(svc *dynamodb.DynamoDB) (exists bool) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(TABLE_NAME),
	}
	if _, err := svc.DescribeTable(input); err != nil {
		// Expecting ResourceNotFoundException.
		// If another error encountered, lets record it.
		expected := false
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				expected = true
			}
		}
		if !expected {
			fmt.Println(err.Error())
		}
		return false
	}
	return true
}

func urlToId(url string) string {
	offset := 1
	if strings.HasSuffix(url, "/") {
		offset = 2
	}
	split := strings.Split(url, "/")
	return split[len(split)-offset]
}

func extract(data *swapiData) {
	fmt.Println("Extracting data from swapi.co...")
	urlToName := make(map[string]string)
	characters, err := storage.GetCharacters()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, character := range characters {
		urlToName[character.URL] = character.Name
	}

	films, err := storage.GetFilms()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, film := range films {
		urlToName[film.URL] = film.Title
	}

	starships, err := storage.GetStarships()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, starship := range starships {
		urlToName[starship.URL] = starship.Name
	}

	species, err := storage.GetSpecies()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, sp := range species {
		urlToName[sp.URL] = sp.Name
	}

	planets, err := storage.GetPlanets()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, planet := range planets {
		urlToName[planet.URL] = planet.Name
	}

	for i, _ := range characters {
		person := &characters[i]
		person.Id = urlToId(person.URL)
		person.Type = "character"
		person.Homeworld = urlToName[person.Homeworld]
		for i, film := range person.Films {
			person.Films[i] = urlToName[film]
		}
		for i, sp := range person.Species {
			person.Species[i] = urlToName[sp]
		}
		for i, starship := range person.Starships {
			person.Starships[i] = urlToName[starship]
		}
	}

	for i, _ := range starships {
		starship := &starships[i]
		starship.Id = urlToId(starship.URL)
		starship.Type = "starship"
		for i, film := range starship.Films {
			starship.Films[i] = urlToName[film]
		}
		for i, pilot := range starship.Pilots {
			starship.Pilots[i] = urlToName[pilot]
		}
	}

	for i, _ := range planets {
		planet := &planets[i]
		planet.Id = urlToId(planet.URL)
		planet.Type = "planet"
		for i, film := range planet.Films {
			planet.Films[i] = urlToName[film]
		}
		for i, resident := range planet.Residents {
			planet.Residents[i] = urlToName[resident]
		}
	}

	data.Characters = characters
	data.Starships = starships
	data.Planets = planets
}

func createPutRequest(in interface{}) *dynamodb.PutRequest {
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		panic(err)
	}
	return &dynamodb.PutRequest{Item: item}
}

func transform(data *swapiData) []*dynamodb.WriteRequest {
	itemCount := len(data.Characters) + len(data.Starships) + len(data.Planets)
	ret := make([]*dynamodb.WriteRequest, 0, itemCount)

	for i, _ := range data.Characters {
		ret = append(ret, &dynamodb.WriteRequest{PutRequest: createPutRequest(&data.Characters[i])})
	}

	for i, _ := range data.Starships {
		ret = append(ret, &dynamodb.WriteRequest{PutRequest: createPutRequest(&data.Starships[i])})
	}

	for i, _ := range data.Planets {
		ret = append(ret, &dynamodb.WriteRequest{PutRequest: createPutRequest(&data.Planets[i])})
	}

	return ret
}

func writeBatch(svc *dynamodb.DynamoDB, batch []*dynamodb.WriteRequest) {
	var items = make(map[string][]*dynamodb.WriteRequest)
	items[TABLE_NAME] = batch
	input := &dynamodb.BatchWriteItemInput{RequestItems: items}

	result, err := svc.BatchWriteItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func load(svc *dynamodb.DynamoDB, all []*dynamodb.WriteRequest) {
	fmt.Printf("Loading %d items into DynamoDB table '%s'...\n", len(all), TABLE_NAME)

	// limit of 25 requests per batch
	batchSize := 25
	numOfBatches := int(len(all)/batchSize) + 1
	fmt.Printf("Number of batches = %d\n", numOfBatches)
	l, c := 0, batchSize
	for i := 0; i < numOfBatches; i++ {
		fmt.Printf("Writing %d - %d\n", l+1, c)
		writeBatch(svc, all[l:c])
		l += batchSize
		c += batchSize
		if c > len(all) {
			c = len(all)
		}
	}
}

func BuildStarWarsDB() {
	// Create a session to share configuration, and load external configuration.
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	svc := dynamodb.New(sess)

	// Create table if it doesn't exist
	if !tableExists(svc) && !createTable(svc) {
		fmt.Printf("Table '%s' could not be created", TABLE_NAME)
		return
	}

	swapiData := &swapiData{}
	extract(swapiData)
	awsData := transform(swapiData)
	load(svc, awsData)
}

//
// Some experimenting below with generic results structure
//
type pagedResponse struct {
	Count    int                      `json:"count"`
	Next     string                   `json:"next"`
	Previous string                   `json:"previous"`
	Results  []map[string]interface{} `json:"results"`
}

func getList(url string) (list []map[string]interface{}, err error) {
	for url != "" {
		var res pagedResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		list = append(list, res.Results...)
	}
	return
}

func extractToGenericMap() {
	urlToName := make(map[string]string)
	characters, err := getList("https://swapi.co/api/people/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, character := range characters {
		urlToName[character["url"].(string)] = character["name"].(string)
	}

	films, err := getList("https://swapi.co/api/films/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, film := range films {
		urlToName[film["url"].(string)] = film["title"].(string)
	}

	starships, err := getList("https://swapi.co/api/starships/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, starship := range starships {
		urlToName[starship["url"].(string)] = starship["name"].(string)
	}

	species, err := getList("https://swapi.co/api/species/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, sp := range species {
		urlToName[sp["url"].(string)] = sp["name"].(string)
	}

	planets, err := getList("https://swapi.co/api/planets/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, planet := range planets {
		urlToName[planet["url"].(string)] = planet["name"].(string)
	}

	fmt.Printf("%+v\n", urlToName)
}
