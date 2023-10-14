package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	DocId primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
}

type User struct {
	UserId string `bson:"userId" json:"userId"`
	Name   string `bson:"name" json:"name"`
	Phone  string `bson:"phone" json:"phone"`
}
