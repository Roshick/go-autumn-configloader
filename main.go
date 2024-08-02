package main

import (
	"fmt"

	"github.com/Roshick/go-autumn-configloader/pkg/configloader"
)

type Config struct {
	Example string
}

func (c *Config) ConfigItems() []*configloader.DefaultConfigItem {
	return []*configloader.DefaultConfigItem{
		{
			Key: "EXAMPLE",
		},
	}
}

func (c *Config) ObtainValues(getter func(string) string) error {
	c.Example = getter("EXAMPLE")
	return nil
}

func main() {
	config := Config{}

	yamlProvider := configloader.CreateYAMLConfigFileProvider("config.yaml")
	envProvider := configloader.CreateEnvironmentVariablesProvider()

	configLoader := configloader.New[*configloader.DefaultConfigItem]()
	if err := configLoader.LoadConfig(&config, yamlProvider, envProvider); err != nil {
		panic("failed to load config values: " + err.Error())
	}

	fmt.Println("successfully loaded config value: ", config.Example)
}

func ptr[E any](e E) *E {
	return &e
}
