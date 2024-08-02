package provider

import (
	"fmt"
	"os"

	"github.com/Roshick/go-autumn-configloader/pkg/configloader"
	"gopkg.in/yaml.v3"
)

type DefaultValuesProviderConfigItem interface {
	configloader.ConfigItem
	GetDefaultValue() *string
}

func CreateDefaultValuesProvider() configloader.Provider[DefaultValuesProviderConfigItem] {
	return func(configItems []DefaultValuesProviderConfigItem) (map[string]string, error) {
		rawValues := make(map[string]string)
		for _, it := range configItems {
			if it.GetDefaultValue() == nil {
				continue
			}
			rawValues[it.GetKey()] = *it.GetDefaultValue()
		}
		return rawValues, nil
	}
}

type YAMLConfigFileProviderConfigItem interface {
	configloader.ConfigItem
	GetConfigFileKey() *string
}

func CreateYAMLConfigFileProvider(filename string) configloader.Provider[YAMLConfigFileProviderConfigItem] {
	return func(configItems []YAMLConfigFileProviderConfigItem) (map[string]string, error) {
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
			if it.GetConfigFileKey() == nil {
				continue
			}
			if value, ok := allValues[*it.GetConfigFileKey()]; ok {
				values[it.GetKey()] = value
			}
		}
		return values, nil
	}
}

type EnvironmentVariablesProviderConfigItem interface {
	configloader.ConfigItem
	GetEnvironmentKey() *string
}

func CreateEnvironmentVariablesProvider() configloader.Provider[EnvironmentVariablesProviderConfigItem] {
	return func(configItems []EnvironmentVariablesProviderConfigItem) (map[string]string, error) {
		values := make(map[string]string)
		for _, it := range configItems {
			if it.GetEnvironmentKey() == nil {
				continue
			}
			value, ok := os.LookupEnv(*it.GetEnvironmentKey())
			if ok {
				values[it.GetKey()] = value
			}
		}
		return values, nil
	}
}
