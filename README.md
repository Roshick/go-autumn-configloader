# go-autumn-configloader

Inspired by the go-autumn framework, go-autumn-configloader enables users to load configuration values from various sources such as configuration files and environment variables, using a system of modular providers.

## Features

- **Extensible**: Easily extend functionality via the provider interface.
- **Default Providers**: Ships with three default providers for YAML configuration files, environment variables, and default values.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
    - [YAML Config Files](#yaml-config-files)
    - [Environment Variables](#environment-variables)
    - [Default Values](#default-values)
- [Extensibility](#extensibility)
    - [Creating Custom Providers](#creating-custom-providers)
- [Examples](#examples)
    - [Loading from YAML, Environment Variables, and Default Values](#loading-from-yaml-environment-variables-and-default-values)
- [Contributing](#contributing)
- [License](#license)

## Installation

Install go-autumn-configloader with the following command:

```sh
go get github.com/Roshick/go-autumn-configloader
```

## Usage

### YAML Config Files

To load configuration from a YAML file, use the YAML provider:

```go
package main

import (
	"fmt"
	
	"github.com/Roshick/go-autumn-configloader/pkg/configloader"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

type Config struct {
	Example string
}

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:         "EXAMPLE",
			Description: "An example configuration item",
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

	configLoader := configloader.New()
	if err := configLoader.LoadConfig(&config, yamlProvider); err != nil {
		panic("failed to load config values: " + err.Error())
	}

	fmt.Println("successfully loaded config value: ", config.Example)
}
```

### Environment Variables

To load configuration values from environment variables, use the environment provider:

```go
package main

import (
	"fmt"
	
	"github.com/Roshick/go-autumn-configloader/pkg/configloader"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

type Config struct {
	Example string
}

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:         "EXAMPLE",
			Description: "An example configuration item",
		},
	}
}

func (c *Config) ObtainValues(getter func(string) string) error {
	c.Example = getter("EXAMPLE")
	return nil
}

func main() {
	config := Config{}

	envProvider := configloader.CreateEnvironmentVariablesProvider()

	configLoader := configloader.New()
	if err := configLoader.LoadConfig(&config, envProvider); err != nil {
		panic("failed to load config values: " + err.Error())
	}

	fmt.Println("successfully loaded config value: ", config.Example)
}
```

### Default Values

To load configuration values from default values, use the default values provider:

```go
package main

import (
	"fmt"
	
	"github.com/Roshick/go-autumn-configloader/pkg/configloader"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

type Config struct {
	Example string
}

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:         "EXAMPLE",
			Default:     "default-value",
			Description: "An example configuration item",
		},
	}
}

func (c *Config) ObtainValues(getter func(string) string) error {
	c.Example = getter("EXAMPLE")
	return nil
}

func main() {
	config := Config{}

	defaultProvider := configloader.CreateDefaultValuesProvider()

	configLoader := configloader.New()
	if err := configLoader.LoadConfig(&config, defaultProvider); err != nil {
		panic("failed to load config values: " + err.Error())
	}

	fmt.Println("successfully loaded config value: ", config.Example)
}
```

## Extensibility

go-autumn-configloader is designed to be extensible. You can create custom providers to load configuration from other sources.

### Creating Custom Providers

Implement the `Provider` interface to create a custom provider:

```go
package main

import (
	"fmt"
	
	"github.com/Roshick/go-autumn-configloader/pkg/configloader"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

type CustomProvider struct {
	// Custom provider fields
}

func (p *CustomProvider) Load(configItems []auconfigapi.ConfigItem) (map[string]string, error) {
	// Implement your custom loading logic here
	return map[string]string{
		"EXAMPLE": "custom-value",
	}, nil
}

type Config struct {
	Example string
}

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:         "EXAMPLE",
			Default:     "",
			Description: "An example configuration item",
		},
	}
}

func (c *Config) ObtainValues(getter func(string) string) error {
	c.Example = getter("EXAMPLE")
	return nil
}

func main() {
	config := Config{}

	customProvider := &CustomProvider{}

	configLoader := configloader.New()
	if err := configLoader.LoadConfig(&config, customProvider); err != nil {
		panic("failed to load config values: " + err.Error())
	}

	fmt.Println("successfully loaded config value: ", config.Example)
}
```

## Examples

### Loading from YAML, Environment Variables, and Default Values

Here's an example demonstrating how to load configuration values from a YAML file, environment variables, and default values. The order of registering providers is important: values from later providers override those from earlier providers.

```go
package main

import (
	"fmt"
	
	"github.com/Roshick/go-autumn-configloader/pkg/configloader"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

type Config struct {
	Example string
}

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
	return []auconfigapi.ConfigItem{
		{
			Key:         "EXAMPLE",
			Default:     "default-value",
			Description: "An example configuration item",
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
	defaultProvider := configloader.CreateDefaultValuesProvider()

	configLoader := configloader.New()
	if err := configLoader.LoadConfig(&config, defaultProvider, yamlProvider, envProvider); err != nil {
		panic("failed to load config values: " + err.Error())
	}

	fmt.Println("successfully loaded config value: ", config.Example)
}
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request, or open an issue to report bugs or request new features.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
