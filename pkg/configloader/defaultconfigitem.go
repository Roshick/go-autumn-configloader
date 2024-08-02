package configloader

type DefaultConfigItem struct {
	Key            string
	Description    string
	DefaultValue   *string
	ConfigFileKey  *string
	EnvironmentKey *string
}

func (c *DefaultConfigItem) GetKey() string {
	return c.Key
}

func (c *DefaultConfigItem) GetDefaultValue() *string {
	return c.DefaultValue
}

func (c *DefaultConfigItem) GetConfigFileKey() *string {
	return c.ConfigFileKey
}

func (c *DefaultConfigItem) GetEnvironmentKey() *string {
	return c.EnvironmentKey
}
