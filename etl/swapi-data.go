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
)

//func GetList(url string, structType interface{}) interface{} {
//	//resultValue := reflect.ValueOf(resultType)
//	//listType := reflect.SliceOf(resultValue.Type())
//	resultType := reflect.TypeOf(structType)
//	listType := reflect.SliceOf(resultType)
//	fmt.Printf("listType = %s\n", listType.String())
//	list := reflect.MakeSlice(listType, 0, 0)
//
//	var err error
//	for url != "" {
//		var res struct {
//			Count    int      `json:"count"`
//			Next     string   `json:"next"`
//			Previous string   `json:"previous"`
//			Results  []map[string]interface{} `json:"results"`
//		}
//		if err = swapi.Get(url, &res); err != nil {
//			fmt.Printf(url)
//		}
//		url = res.Next
//		for _, person := range res.Results {
//			fmt.Printf("%+v\n", person)
//			// TODO problem here because person is a map and I'm trying to add map to slice of etl.Person
//			val := reflect.ValueOf(person)
//			reflect.AppendSlice(list, val)
//		}
//	}
//	return list.Interface()
//}

func GetList(url string) (list []map[string]interface{}, err error) {
	for url != "" {
		var res struct {
			Count    int                      `json:"count"`
			Next     string                   `json:"next"`
			Previous string                   `json:"previous"`
			Results  []map[string]interface{} `json:"results"`
		}
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		for _, one := range res.Results {
			list = append(list, one)
		}
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
	for i, person := range people {
		people[i].Homeworld = urlToName[person.Homeworld] // TODO why? we must get shallow copy perhaps
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

	for _, starship := range starships {
		for i, film := range starship.Films {
			starship.Films[i] = urlToName[film]
		}
		for i, pilot := range starship.Pilots {
			starship.Pilots[i] = urlToName[pilot]
		}
	}

	for _, planet := range planets {
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
}

func WriteBatchToDynamoDB() {
	url := "https://swapi.co/api/people/"
	var list []map[string]interface{}
	var err error
	if list, err = GetList(url); err != nil {
		fmt.Printf("some error occurred")
		return
	}
	//fmt.Printf("%+v\n", list)
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