package etl

import (
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func Transform(filename string) []*dynamodb.WriteRequest {
	data := &swData{}
	fromFile(filename, data)

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

func fromFile(filename string, data *swData) {
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonBytes, data)
}

func createPutRequest(in interface{}) *dynamodb.PutRequest {
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		panic(err)
	}
	return &dynamodb.PutRequest{Item: item}
}
