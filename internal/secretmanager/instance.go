package secretmanager

import (
	"context"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type secretManager struct {
	client *secretsmanager.Client
}

type SecretManager interface {
	GetSecretString(secretName string) (*string, error)
}

func NewInstance(region string, awsConfig aws.Config) (sm SecretManager, err error) {
	if reflect.DeepEqual(awsConfig, aws.Config{}) {
		if awsConfig, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
		); err != nil {
			return nil, err
		}
	}

	return secretManager{
		client: secretsmanager.NewFromConfig(awsConfig),
	}, err
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
