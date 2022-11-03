package envloader

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/bavatech/envloader/internal/secretmanager"
	"github.com/joho/godotenv"
)

const (
	OptionalKey = "optional"
	DefaultKey  = "default"
)

func loadSecretManager(vars interface{}) (map[string]interface{}, error) {
	var secretsMap map[string]interface{}

	secretName := os.Getenv("SECRET_NAME")
	region := os.Getenv("REGION")

	if secretName != "" {
		sm, err := secretmanager.NewInstance(region)
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

func loadEnvVars(vars interface{}, secretsMap map[string]interface{}) error {
	pointr := reflect.ValueOf(vars)
	values := pointr.Elem()
	typeOfValues := values.Type()

	for i := 0; i < values.NumField(); i++ {
		value := values.Field(i)
		field := pointr.Elem().Field(i)
		fieldName := typeOfValues.Field(i).Name

		// value := values.Field(i).String()
		// field := pointr.Elem().Field(i)
		// fieldName := typeOfValues.Field(i).Name

		fieldKey := fieldName
		optional := false
		defaultValue := ""

		tag := typeOfValues.Field(i).Tag.Get("env")
		if tag != "" {
			tagOpts := strings.Split(tag, ",")
			fieldKey = tagOpts[0]
			keys := tagOpts[1:]
			for _, key := range keys {
				if key == OptionalKey {
					optional = true
				} else if strings.Index(key, DefaultKey+"=") == 0 && value.String() == "" {
					opts := strings.Split(key, "=")
					defaultValue = opts[1]
				}
			}
		}

		fmt.Println(value, field, fieldName, os.Getenv(fieldKey), optional, defaultValue)

		if field.CanSet() {
			envValue := os.Getenv(fieldKey)

			switch field.Kind() {
			case reflect.Slice:
				splitEnvValue := strings.Split(envValue, " ")

				switch field.Interface().(type) {
				case []int:
					for _, strValue := range splitEnvValue {
						intValue, err := strconv.Atoi(strValue)

						if err != nil {
							return err
						}

						value.Set(reflect.Append(value, reflect.ValueOf(intValue)))
					}
				case []int32:
					for _, strValue := range splitEnvValue {
						intValue, err := strconv.Atoi(strValue)

						if err != nil {
							return err
						}

						value.Set(reflect.Append(value, reflect.ValueOf(int32(intValue))))
					}
				case []int64:
					for _, strValue := range splitEnvValue {
						intValue, err := strconv.Atoi(strValue)

						if err != nil {
							return err
						}

						value.Set(reflect.Append(value, reflect.ValueOf(int64(intValue))))
					}
				case []string:
					for _, strValue := range splitEnvValue {
						value.Set(reflect.Append(value, reflect.ValueOf(strValue)))
					}
				}
			}
			// field.SetString(os.Getenv(fieldKey))

			// fmt.Println(field.String())

			// if secretsMap[fieldKey] != nil {
			// 	field.SetString(secretsMap[fieldKey].(string))
			// }

			// if field.String() == "" {
			// 	field.SetString(defaultValue)
			// }
		}

		if !optional && field.String() == "" {
			return fmt.Errorf(`env "%s", fieldname "%s" must be defined`, fieldKey, fieldName)
		}
	}

	return nil
}

func Load(vars interface{}, filenames ...string) error {
	godotenv.Load(filenames...)

	secretsMap, err := loadSecretManager(vars)
	if err != nil {
		return err
	}

	err = loadEnvVars(vars, secretsMap)
	if err != nil {
		return err
	}

	return nil
}
