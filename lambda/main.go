package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// Define the DynamoDB client interface
type DynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

// Define the response structure
type Response struct {
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	StatusCode int               `json:"statusCode"`
}

// Handler function
func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (Response, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	tableName := "Resume-api"
	key := event.PathParameters["id"]

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to get item from DynamoDB: %v", err)
		return Response{StatusCode: 500, Body: "Internal Server Error", Headers: map[string]string{"Content-Type": "application/json"}}, nil
	}

	if result.Item == nil {
		return Response{StatusCode: 404, Body: "Not Found", Headers: map[string]string{"Content-Type": "application/json"}}, nil
	}

	jsonResponse, err := json.Marshal(result.Item)
	if err != nil {
		log.Fatalf("Failed to marshal result to JSON: %v", err)
		return Response{StatusCode: 500, Body: "Internal Server Error", Headers: map[string]string{"Content-Type": "application/json"}}, nil
	}

	return Response{
		StatusCode: 200,
		Body:       string(jsonResponse),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
