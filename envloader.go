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
				} else if strings.Index(key, DefaultKey+"=") == 0 {
					opts := strings.Split(key, "=")

					if value.CanSet() {
						defaultValue = opts[1]
					}
				}
			}
		}

		if field.CanSet() {
			envValue := os.Getenv(fieldKey)

			if secretsMap[fieldKey] != nil {
				envValue = secretsMap[fieldKey].(string)
			}

			if envValue == "" && defaultValue != "" {
				envValue = defaultValue
			}

			switch field.Kind() {
			case reflect.Slice:
				splitEnvValue := strings.Split(envValue, " ")

				switch field.Interface().(type) {
				case []int:
					if err := appendIntValues[int](splitEnvValue, value); err != nil {
						return err
					}
				case []int16:
					if err := appendIntValues[int16](splitEnvValue, value); err != nil {
						return err
					}
				case []int32:
					if err := appendIntValues[int32](splitEnvValue, value); err != nil {
						return err
					}
				case []int64:
					if err := appendIntValues[int64](splitEnvValue, value); err != nil {
						return err
					}
				case []float32:
					if err := appendFloatValues[float32](splitEnvValue, value); err != nil {
						return err
					}
				case []float64:
					if err := appendFloatValues[float64](splitEnvValue, value); err != nil {
						return err
					}
				case []string:
					value.Set(reflect.ValueOf(splitEnvValue))
				}
			case reflect.String:
				if field.String() == "" {
					field.SetString(envValue)
				}
			case reflect.Int:
				if err := appendIntValue[int](envValue, value); err != nil {
					return err
				}
			case reflect.Int16:
				if err := appendIntValue[int16](envValue, value); err != nil {
					return err
				}
			case reflect.Int32:
				if err := appendIntValue[int32](envValue, value); err != nil {
					return err
				}
			case reflect.Int64:
				if err := appendIntValue[int64](envValue, value); err != nil {
					return err
				}
			case reflect.Float32:
				if err := appendFloatValue[float32](envValue, value); err != nil {
					return err
				}
			case reflect.Float64:
				if err := appendFloatValue[float64](envValue, value); err != nil {
					return err
				}
			}

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

func scanNumericsValues[T int | int16 | int32 | int64](values []string) ([]T, error) {
	var ret []T

	for _, strValue := range values {
		intValue, err := strconv.Atoi(strValue)

		if err != nil {
			return ret, err
		}

		ret = append(ret, T(intValue))
	}

	return ret, nil
}

func appendIntValues[T int | int16 | int32 | int64](splitEnvValue []string, reflectValue reflect.Value) error {
	scannedValues, err := scanNumericsValues[T](splitEnvValue)

	if err != nil {
		return err
	}

	reflectValue.Set(reflect.ValueOf(scannedValues))

	return nil
}

func scanFloatValues[T float32 | float64](values []string) ([]T, error) {
	var floatSize int
	var ret []T

	switch any(ret).(type) {
	case []float32:
		floatSize = 32
	case []float64:
		floatSize = 64
	}

	for _, strValue := range values {
		strConv, err := strconv.ParseFloat(strValue, floatSize)

		if err != nil {
			return ret, err
		}

		ret = append(ret, T(strConv))
	}

	return ret, nil
}

func appendFloatValues[T float32 | float64](splitEnvValue []string, reflectValue reflect.Value) error {
	scannedValues, err := scanFloatValues[T](splitEnvValue)

	if err != nil {
		return err
	}

	reflectValue.Set(reflect.ValueOf(scannedValues))

	return nil
}

func appendIntValue[T int | int16 | int32 | int64](value string, reflectValue reflect.Value) error {
	strConv, err := strconv.Atoi(value)

	if err != nil {
		return err
	}

	reflectValue.Set(reflect.ValueOf(T(strConv)))

	return nil
}

func appendFloatValue[T float32 | float64](value string, reflectValue reflect.Value) (err error) {
	var floatSize int
	var ret T

	switch any(ret).(type) {
	case float32:
		floatSize = 32
	case float64:
		floatSize = 64
	}

	strConv, err := strconv.ParseFloat(value, floatSize)

	if err != nil {
		return err
	}

	reflectValue.Set(reflect.ValueOf(T(strConv)))

	return nil
}
