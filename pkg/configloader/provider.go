package configloader

import (
	"fmt"
	"os"

	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	"gopkg.in/yaml.v3"
)

type Provider func([]auconfigapi.ConfigItem) (map[string]string, error)

func CreateDefaultValuesProvider() Provider {
	return func(configItems []auconfigapi.ConfigItem) (map[string]string, error) {
		rawValues := make(map[string]string)
		for _, it := range configItems {
			defaultValue, ok := it.Default.(string)
			if !ok {
				return nil, fmt.Errorf("failed to load default value of key %s: value is not a string", it.Key)
			}
			rawValues[it.Key] = defaultValue
		}
		return rawValues, nil
	}
}

func CreateYAMLConfigFileProvider(filename string) Provider {
	return func(configItems []auconfigapi.ConfigItem) (map[string]string, error) {
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			// this is NOT an error
			return nil, nil
		}

		reader, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open configuration file %s: %s", filename, err.Error())
		}
		defer reader.Close()
		decoder := yaml.NewDecoder(reader)
		decoder.KnownFields(true)

		var allValues = make(map[string]string)
		if err = decoder.Decode(&allValues); err != nil {
			return nil, fmt.Errorf("failed to parse configuration file %s: %s", filename, err.Error())
		}

		values := make(map[string]string)
		for _, it := range configItems {
			if value, ok := allValues[it.Key]; ok {
				values[it.Key] = value
			}
		}
		return values, nil
	}
}

func CreateEnvironmentVariablesProvider() Provider {
	return func(configItems []auconfigapi.ConfigItem) (map[string]string, error) {
		values := make(map[string]string)
		for _, it := range configItems {
			envValue, ok := os.LookupEnv(it.EnvName)
			if ok {
				values[it.Key] = envValue
			}
		}
		return values, nil
	}
}
