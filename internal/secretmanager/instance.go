package secretmanager

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type secretManager struct {
	client *secretsmanager.Client
}

type SecretManager interface {
	GetSecretString(secretName string) (*string, error)
}

func NewInstance(region string) (SecretManager, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	return secretManager{
		client: secretsmanager.NewFromConfig(cfg),
	}, nil
}

func (s secretManager) GetSecretString(secretName string) (*string, error) {
	secret, err := s.client.GetSecretValue(
		context.TODO(),
		&secretsmanager.GetSecretValueInput{SecretId: &secretName},
	)

	if err != nil {
		return nil, err
	}

	return secret.SecretString, nil
}
