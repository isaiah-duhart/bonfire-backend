package utils

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Secrets struct {
	DB_URL     string `json:"DB_URL"`
	JWT_SECRET string `json:"JWT_SECRET"`
}

func IsRunningOnEC2() bool {
	client := http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get("http://169.254.169.254/latest/meta-data/")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
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

	return secrets
}
