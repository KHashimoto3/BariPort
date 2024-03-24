package main

import (
	"bariport"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(bariport.HandlerGetMessages)
}
