package configloader

import (
	"fmt"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

type ConfigLoader struct {
	providers []Provider
	values    map[string]string
}

type Config interface {
	ConfigItems() []auconfigapi.ConfigItem

	ObtainValues(getter func(string) string) error
}

type Provider func([]auconfigapi.ConfigItem) (map[string]string, error)

func New() *ConfigLoader {
	return &ConfigLoader{
		values: make(map[string]string),
	}
}

func (l *ConfigLoader) LoadConfig(config Config, providers ...Provider) error {
	if err := l.LoadValues(config.ConfigItems(), providers...); err != nil {
		return err
	}
	return config.ObtainValues(l.Get)
}

func (l *ConfigLoader) Get(key string) string {
	return l.values[key]
}

func (l *ConfigLoader) LoadValues(
	configItems []auconfigapi.ConfigItem,
	providers ...Provider,
) error {
	values, err := loadValues(configItems, providers...)
	if err != nil {
		return err
	}
	for key, value := range values {
		l.values[key] = value
	}
	return nil
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

// Validate checks all values using the validation function
// configured in configItems, if any.
//
// The first return value is a map of individual validation errors
// by configuration item key, or nil if no errors occurred.
//
// The second return value indicates overall validation success (nil)
// or failure with a single summary error.
//
// Validate should be called after respective LoadConfig() and/or LoadValues()
// calls have been made for the keys being validated.
func (l *ConfigLoader) Validate(
	configItems []auconfigapi.ConfigItem,
) (map[string]error, error) {
	validationErrors := make(map[string]error)
	for _, item := range configItems {
		if item.Validate != nil {
			err := item.Validate(item.Key)
			if err != nil {
				validationErrors[item.Key] = err
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors, fmt.Errorf("configuration validation failed for %d configuration keys", len(validationErrors))
	} else {
		return nil, nil
	}
}
