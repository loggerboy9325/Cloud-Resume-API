package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	body := fmt.Sprintf("Your method is %s", req.RequestContext.HTTP.Method)

	headers := map[string]string{
		"Content-Type": "text/plain",
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       body,
		StatusCode: 200,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(handler)
}
