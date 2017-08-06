package etl

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const tableName = "StarWars"

func Load(all []*dynamodb.WriteRequest) {
	// Create a session to share configuration, and load external configuration.
	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	svc := dynamodb.New(sess)

	fmt.Printf("Loading %d items into DynamoDB table '%s'...\n", len(all), tableName)

	// Create table if it doesn't exist
	if !tableExists(svc) && !createTable(svc) {
		fmt.Printf("Table '%s' could not be created", tableName)
		return
	}

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
		TableName: aws.String(tableName),
	}
	fmt.Printf("Creating DynamoDB table '%s'...\n", tableName)
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
		TableName: aws.String(tableName),
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

func writeBatch(svc *dynamodb.DynamoDB, batch []*dynamodb.WriteRequest) {
	var items = make(map[string][]*dynamodb.WriteRequest)
	items[tableName] = batch
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
