package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Define the structure of the JSON data
type Resume struct {
	Basics struct {
		Name     string `json:"name"`
		Label    string `json:"label"`
		Image    string `json:"image"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		URL      string `json:"url"`
		Summary  string `json:"summary"`
		Location struct {
			Address     string `json:"address"`
			PostalCode  string `json:"postalCode"`
			City        string `json:"city"`
			CountryCode string `json:"countryCode"`
			Region      string `json:"region"`
		} `json:"location"`
		Profiles []struct {
			Network  string `json:"network"`
			Username string `json:"username"`
			URL      string `json:"url"`
		} `json:"profiles"`
	} `json:"basics"`
	Work []struct {
		Name       string   `json:"name"`
		Position   string   `json:"position"`
		URL        string   `json:"url"`
		StartDate  string   `json:"startDate"`
		EndDate    string   `json:"endDate"`
		Summary    string   `json:"summary"`
		Highlights []string `json:"highlights"`
	} `json:"work"`
	Volunteer []struct {
		Organization string   `json:"organization"`
		Position     string   `json:"position"`
		URL          string   `json:"url"`
		StartDate    string   `json:"startDate"`
		EndDate      string   `json:"endDate"`
		Summary      string   `json:"summary"`
		Highlights   []string `json:"highlights"`
	} `json:"volunteer"`
	Education []struct {
		Institution string   `json:"institution"`
		URL         string   `json:"url"`
		Area        string   `json:"area"`
		StudyType   string   `json:"studyType"`
		StartDate   string   `json:"startDate"`
		EndDate     string   `json:"endDate"`
		Score       string   `json:"score"`
		Courses     []string `json:"courses"`
	} `json:"education"`
	Awards []struct {
		Title   string `json:"title"`
		Date    string `json:"date"`
		Awarder string `json:"awarder"`
		Summary string `json:"summary"`
	} `json:"awards"`
	Certificates []struct {
		Name   string `json:"name"`
		Date   string `json:"date"`
		Issuer string `json:"issuer"`
		URL    string `json:"url"`
	} `json:"certificates"`
	Publications []struct {
		Name        string `json:"name"`
		Publisher   string `json:"publisher"`
		ReleaseDate string `json:"releaseDate"`
		URL         string `json:"url"`
		Summary     string `json:"summary"`
	} `json:"publications"`
	Skills []struct {
		Name     string   `json:"name"`
		Level    string   `json:"level"`
		Keywords []string `json:"keywords"`
	} `json:"skills"`
	Languages []struct {
		Language string `json:"language"`
		Fluency  string `json:"fluency"`
	} `json:"languages"`
	Interests []struct {
		Name     string   `json:"name"`
		Keywords []string `json:"keywords"`
	} `json:"interests"`
	References []struct {
		Name      string `json:"name"`
		Reference string `json:"reference"`
	} `json:"references"`
	Projects []struct {
		Name        string   `json:"name"`
		StartDate   string   `json:"startDate"`
		EndDate     string   `json:"endDate"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Highlights  []string `json:"highlights"`
	} `json:"projects"`
}

var (
	tableName  = "Resume-api"
	primaryKey = "id"
	resumeID   = "1"
)

func fetchResumeFromDynamoDB(ctx context.Context, svc *dynamodb.Client) (*Resume, error) {
	result, err := svc.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			primaryKey: &types.AttributeValueMemberS{Value: resumeID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("no item found in DynamoDB with id: %s", resumeID)
	}

	var resume Resume
	err = attributevalue.UnmarshalMap(result.Item, &resume)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DynamoDB item: %w", err)
	}

	return &resume, nil
}

func handler(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	resume, err := fetchResumeFromDynamoDB(ctx, svc)
	if err != nil {
		log.Printf("failed to fetch resume from DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	responseBody, err := json.Marshal(resume)
	if err != nil {
		log.Printf("failed to marshal response body: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}
