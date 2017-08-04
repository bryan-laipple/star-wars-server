package etl

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/leejarvis/swapi"
	//"github.com/aws/aws-sdk-go/aws"
	"encoding/json"
	"strconv"
	"strings"
)

type PagedResponse struct {
	Count    int                      `json:"count"`
	Next     string                   `json:"next"`
	Previous string                   `json:"previous"`
	Results  []map[string]interface{} `json:"results"`
}

func GetList(url string) (list []map[string]interface{}, err error) {
	for url != "" {
		var res PagedResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		list = append(list, res.Results...)
	}
	return
}

func BuildDynamoDBTable() {
	// Create a session to share configuration, and load external configuration.
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	svc := dynamodb.New(sess)

	input := &dynamodb.ListTablesInput{}

	result, err := svc.ListTables(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
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

func DataStuffWithGenericMap() {
	urlToName := make(map[string]string)
	people, err := GetList("https://swapi.co/api/people/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, person := range people {
		urlToName[person["url"].(string)] = person["name"].(string)
	}

	films, err := GetList("https://swapi.co/api/films/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, film := range films {
		urlToName[film["url"].(string)] = film["title"].(string)
	}

	starships, err := GetList("https://swapi.co/api/starships/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, starship := range starships {
		urlToName[starship["url"].(string)] = starship["name"].(string)
	}

	species, err := GetList("https://swapi.co/api/species/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, sp := range species {
		urlToName[sp["url"].(string)] = sp["name"].(string)
	}

	planets, err := GetList("https://swapi.co/api/planets/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, planet := range planets {
		urlToName[planet["url"].(string)] = planet["name"].(string)
	}


	fmt.Printf("%+v\n", urlToName)
}

func toId(url string) (int, error) {
	offset := 1
	if strings.HasSuffix(url, "/") {
		offset = 2
	}
	split := strings.Split(url, "/")
	return strconv.Atoi(split[len(split)-offset])
}

func toKey(typ string, id int) string {
	return strings.Join([]string{typ, ":", strconv.Itoa(id)}, "")
}

func DataStuff() {
	urlToName := make(map[string]string)
	people, err := GetPeople()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, person := range people {
		urlToName[person.URL] = person.Name
	}

	films, err := GetFilms()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, film := range films {
		urlToName[film.URL] = film.Title
	}

	starships, err := GetStarships()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, starship := range starships {
		urlToName[starship.URL] = starship.Name
	}

	species, err := GetSpecies()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, sp := range species {
		urlToName[sp.URL] = sp.Name
	}

	planets, err := GetPlanets()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, planet := range planets {
		urlToName[planet.URL] = planet.Name
	}

	// TODO URL to Id...somehow
	for i, _ := range people {
		person := &people[i]
		person.Id, err = toId(person.URL) // TODO error
		person.Type = "people"
		person.Key = toKey(person.Type, person.Id)
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
		for i, film := range starship.Films {
			starship.Films[i] = urlToName[film]
		}
		for i, pilot := range starship.Pilots {
			starship.Pilots[i] = urlToName[pilot]
		}
	}

	for i, _ := range planets {
		planet := &planets[i]
		for i, film := range planet.Films {
			planet.Films[i] = urlToName[film]
		}
		for i, resident := range planet.Residents {
			planet.Residents[i] = urlToName[resident]
		}
	}

	data := struct {
		People    []Person   `json:"people"`
		Starships []Starship `json:"starships"`
		Planets   []Planet   `json:"planets"`
	}{
		people,
		starships,
		planets,
	}

	//fmt.Printf("%+v\n", data)
	jsonData, _ := json.Marshal(&data)
	fmt.Println(string(jsonData))

	//WriteBatchToDynamoDB(people)// TODO https://golang.org/doc/faq#convert_slice_of_interface
}

func WriteBatchToDynamoDB(list []interface{}) {
	itemList := make([]map[string]*dynamodb.AttributeValue, len(list), len(list))
	for i, one := range list {
		var item map[string]*dynamodb.AttributeValue
		var err error
		if item, err = dynamodbattribute.MarshalMap(one); err != nil {
			//error
		}
		itemList[i] = item
	}
	fmt.Printf("%+v\n", itemList)
	//sess := session.Must(session.NewSession())
	//svc := dynamodb.New(sess)
	//input := &dynamodb.BatchWriteItemInput{
	//	RequestItems: map[string][]*dynamodb.WriteRequest{
	//		"Music": {
	//			{
	//				PutRequest: &dynamodb.PutRequest{
	//					Item: map[string]*dynamodb.AttributeValue{
	//						"AlbumTitle": {
	//							S: aws.String("Somewhat Famous"),
	//						},
	//						"Artist": {
	//							S: aws.String("No One You Know"),
	//						},
	//						"SongTitle": {
	//							S: aws.String("Call Me Today"),
	//						},
	//					},
	//				},
	//			},
	//			{
	//				PutRequest: &dynamodb.PutRequest{
	//					Item: map[string]*dynamodb.AttributeValue{
	//						"AlbumTitle": {
	//							S: aws.String("Songs About Life"),
	//						},
	//						"Artist": {
	//							S: aws.String("Acme Band"),
	//						},
	//						"SongTitle": {
	//							S: aws.String("Happy Day"),
	//						},
	//					},
	//				},
	//			},
	//			{
	//				PutRequest: &dynamodb.PutRequest{
	//					Item: map[string]*dynamodb.AttributeValue{
	//						"AlbumTitle": {
	//							S: aws.String("Blue Sky Blues"),
	//						},
	//						"Artist": {
	//							S: aws.String("No One You Know"),
	//						},
	//						"SongTitle": {
	//							S: aws.String("Scared of My Shadow"),
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	//
	//result, err := svc.BatchWriteItem(input)
	//if err != nil {
	//	if aerr, ok := err.(awserr.Error); ok {
	//		switch aerr.Code() {
	//		case dynamodb.ErrCodeProvisionedThroughputExceededException:
	//			fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
	//		case dynamodb.ErrCodeResourceNotFoundException:
	//			fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
	//		case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
	//			fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
	//		case dynamodb.ErrCodeInternalServerError:
	//			fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
	//		default:
	//			fmt.Println(aerr.Error())
	//		}
	//	} else {
	//		// Print the error, cast err to awserr.Error to get the Code and
	//		// Message from an error.
	//		fmt.Println(err.Error())
	//	}
	//	return
	//}
	//
	//fmt.Println(result)
}
