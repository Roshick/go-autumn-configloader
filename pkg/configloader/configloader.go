package configloader

type ConfigItem interface {
	GetKey() string
}

type ConfigLoader[C ConfigItem] struct {
	providers []Provider[C]
	values    map[string]string
}

type Config[C ConfigItem] interface {
	ConfigItems() []C

	ObtainValues(getter func(string) string) error
}

type Provider[C ConfigItem] func([]C) (map[string]string, error)

func New[C ConfigItem]() *ConfigLoader[C] {
	return &ConfigLoader[C]{
		values: make(map[string]string),
	}
}

func (l *ConfigLoader[C]) LoadConfig(config Config[C], providers ...Provider[C]) error {
	if err := l.LoadValues(config.ConfigItems(), providers...); err != nil {
		return err
	}
	return config.ObtainValues(l.Get)
}

func (l *ConfigLoader[C]) Get(key string) string {
	return l.values[key]
}

func (l *ConfigLoader[C]) LoadValues(
	configItems []C,
	providers ...Provider[C],
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

func loadValues[C ConfigItem](
	configItems []C,
	providers ...Provider[C],
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
