package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Document struct {
	DocId primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
}

type User struct {
	DocId    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
}

type UserDetail struct {
	DocId     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId    primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	Following int32              `bson:"following" json:"following"`
	Followers int32              `bson:"followers" json:"followers"`
}
