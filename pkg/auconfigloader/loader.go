package auconfigloader

import (
	"errors"
	"fmt"

	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

type ConfigurationLoader struct {
	providers []Provider
}

func New(providers ...Provider) *ConfigurationLoader {
	return &ConfigurationLoader{
		providers: providers,
	}
}

func (l *ConfigurationLoader) ObtainValues(
	configItems []auconfigapi.ConfigItem,
) (map[string]string, error) {
	values, err := loadValues(configItems, l.providers...)
	if err != nil {
		return nil, err
	}
	if err = validateValues(configItems, values); err != nil {
		return nil, err
	}
	return values, nil
}

func loadValues(
	configItems []auconfigapi.ConfigItem,
	providers ...Provider,
) (map[string]string, error) {
	rawValues := make(map[string]string)
	for _, provider := range providers {
		if provider == nil {
			continue
		}
		loaderRawValues, err := provider(configItems)
		if err != nil {
			return nil, err
		}
		for key, value := range loaderRawValues {
			rawValues[key] = value
		}
	}
	return rawValues, nil
}

func validateValues(
	configItems []auconfigapi.ConfigItem,
	values map[string]string,
) error {
	var errs = make([]error, 0)
	for _, it := range configItems {
		if it.Validate != nil {
			err := it.Validate(values[it.Key])
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to validate configuration value of %s: %s", it.Key, err.Error()))
			}
		}
	}
	return errors.Join(errs...)
}
