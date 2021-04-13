package main

import (
	"animals/internal/animal"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var s *animal.Service

func init() {
	region, ok := os.LookupEnv("DYNAMODB_AWS_REGION")
	if !ok {
		panic("DYNAMODB_AWS_REGION not set")
	}

	tableName, ok := os.LookupEnv("DYNAMODB_TABLE")
	if !ok {
		panic("DYNAMODB_TABLE not set")
	}

	s = animal.NewService(tableName, region)
}

func main() {
	lambda.Start(func() (animal.Animals, error) {
		fmt.Println("ListAnimals invoked!")
		return s.List()
	})
}
