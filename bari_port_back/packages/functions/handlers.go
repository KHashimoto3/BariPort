package bariport

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func HandlerHello(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	body, _ := json.Marshal(HelloResponse{
		Message: "Hello, World!",
	})
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}