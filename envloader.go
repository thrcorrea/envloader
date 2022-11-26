package envloader

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/joho/godotenv"
)

const (
	OptionalKey = "optional"
	DefaultKey  = "default"
)

type options struct {
	godotenv  []string
	awsConfig aws.Config
}

type Option interface {
	apply(*options)
}

func Load(vars interface{}, configOptions ...Option) error {
	options := options{
		godotenv:  []string{},
		awsConfig: aws.Config{},
	}

	for _, opt := range configOptions {
		opt.apply(&options)
	}

	godotenv.Load(options.godotenv...)

	secretsMap, err := loadSecretManager(vars, options.awsConfig)
	if err != nil {
		return err
	}

	err = loadEnvVars(vars, secretsMap)
	if err != nil {
		return err
	}

	return nil
}
