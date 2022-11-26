package envloader

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bavatech/envloader/internal/secretmanager"
)

type awsConfigOption struct {
	AwsConfig aws.Config
}

func (g awsConfigOption) apply(opts *options) {
	opts.awsConfig = g.AwsConfig
}

func WithAwsConfig(awsConfig aws.Config) Option {
	return awsConfigOption{
		AwsConfig: awsConfig,
	}
}

func loadSecretManager(vars interface{}, awsConfig aws.Config) (map[string]interface{}, error) {
	var secretsMap map[string]interface{}

	secretName := os.Getenv("SECRET_NAME")
	region := os.Getenv("REGION")

	if secretName != "" {
		sm, err := secretmanager.NewInstance(region, awsConfig)
		if err != nil {
			return nil, err
		}

		secret, err := sm.GetSecretString(secretName)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(*secret), &secretsMap); err != nil {
			return nil, err
		}
	}

	return secretsMap, nil
}
