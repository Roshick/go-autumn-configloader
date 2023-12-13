package configloader

import (
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
