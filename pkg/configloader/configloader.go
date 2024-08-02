package configloader

type ConfigLoader[ConfigItem any] struct {
	providers []Provider[ConfigItem]
	values    map[string]string
}

type Config[ConfigItem any] interface {
	ConfigItems() []ConfigItem

	ObtainValues(getter func(string) string) error
}

type Provider[ConfigItem any] func([]ConfigItem) (map[string]string, error)

func New[ConfigItem any]() *ConfigLoader[ConfigItem] {
	return &ConfigLoader[ConfigItem]{
		values: make(map[string]string),
	}
}

func (l *ConfigLoader[ConfigItem]) LoadConfig(config Config[ConfigItem], providers ...Provider[ConfigItem]) error {
	if err := l.LoadValues(config.ConfigItems(), providers...); err != nil {
		return err
	}
	return config.ObtainValues(l.Get)
}

func (l *ConfigLoader[ConfigItem]) Get(key string) string {
	return l.values[key]
}

func (l *ConfigLoader[ConfigItem]) LoadValues(
	configItems []ConfigItem,
	providers ...Provider[ConfigItem],
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

func loadValues[ConfigItem any](
	configItems []ConfigItem,
	providers ...Provider[ConfigItem],
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
