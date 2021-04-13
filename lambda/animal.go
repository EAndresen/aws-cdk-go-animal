package animal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Service struct {
	table dynamo.Table
}

func NewService(tableName, region string) *Service {
	sess, _ := session.NewSession()
	db := dynamo.New(sess, &aws.Config{Region: aws.String(region)})
	table := db.Table(tableName)

	return &Service{table: table}
}

type Animal struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

type Animals []Animal

func (s Service) Create(animal Animal) (Animal, error) {
	return animal, s.table.Put(animal).Run()
}

func (s Service) List() (Animals, error) {
	var animals Animals
	err := s.table.Scan().All(&animals)
	return animals, err
}
