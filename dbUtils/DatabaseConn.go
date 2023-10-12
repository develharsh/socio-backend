package dbUtils

import (
	"context"
	"fmt"

	globals "github.com/harshvsinghme/socio-backend.git/global"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var RedisClient *redis.Client

func InitRedisConn() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     globals.SECRETS.REDIS_URL,
		Password: globals.SECRETS.REDIS_PASSW,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		logrus.Fatal(fmt.Sprintf("Failed to connect to Redis: %s", err))
	}

	logrus.Info("Connected to Redis")
	RedisClient = redisClient
}

func InitMongoConn() {
	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI(globals.SECRETS.MONGO_URL)

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	// Ping MongoDB to verify the connection
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Connected to MongoDB")

	MongoClient = mongoClient

}
