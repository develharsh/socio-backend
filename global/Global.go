package globals

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var ENV string

type listofsecrets struct {
	SELF_URL       string
	MONGO_URL      string
	MONGO_DATABASE string
	REDIS_URL      string
	REDIS_PASSW    string
	JWT_SECRET     string
}

var SECRETS = listofsecrets{}

func LoadGlobals(envPath string) {

	err := godotenv.Load(envPath)
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	ENV = os.Getenv("ENV")
	if ENV == "" {
		logrus.Fatal("ENV in .env file is missing")
	}

	result := GetSecretValue(ENV + "_SECRETS")

	json.Unmarshal([]byte(result), &SECRETS)

}

func GetSecretValue(secretName string) (secretString string) {
	region := "ap-south-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		logrus.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	secretString = *result.SecretString

	return secretString
}
