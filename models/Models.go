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

type Post struct {
	DocID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Content     string             `bson:"content" json:"content"`
	PostedBy    primitive.ObjectID `bson:"postedBy" json:"postedBy"`
	Likes       int32              `bson:"likes" json:"likes"`
	Comments    int32              `bson:"comments" json:"comments"`
	IsPublished bool               `bson:"isPublished" json:"isPublished"` //is public?
	CreatedAt   int64              `bson:"created_at" json:"created_at"`
	UpdatedAt   int64              `bson:"updated_at" json:"updated_at"`
}

type Like struct {
	DocID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
}

type Comment struct {
	DocID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}
