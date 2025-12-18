package utils

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secrets struct {
	DB_URL     string `json:"DB_URL"`
	JWT_SECRET string `json:"JWT_SECRET"`
}

func GetAWSSecrets() Secrets {
	secretName := "bonfire-secrets"
	region := "us-east-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	secrets := Secrets{}
	if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		log.Fatalf("Error decoding params: %v", err)
	}

	log.Println("Retrieved secrets!")
	return secrets
}
