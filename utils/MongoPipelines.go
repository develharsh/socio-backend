package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PipelineUserAccount(userId primitive.ObjectID) primitive.A {
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userId}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "userdetails"},
					{Key: "localField", Value: "_id"},
					{Key: "foreignField", Value: "userId"},
					{Key: "as", Value: "userdetails"},
				},
			},
		},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$userdetails"}}}},
		bson.D{
			{Key: "$project",
				Value: bson.D{
					{Key: "_id", Value: false},
					{Key: "name", Value: "$name"},
					{Key: "userId", Value: "$_id"},
					{Key: "email", Value: "$email"},
					{Key: "followers", Value: "$userdetails.followers"},
					{Key: "following", Value: "$userdetails.following"},
				},
			},
		},
	}
	return pipeline
}
