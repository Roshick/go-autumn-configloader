# go-autumn-configloader

Inspired by the go-autumn framework, go-autumn-configloader allows users to load configuration values from various sources like configuration files and environment variables using a system of modular providers.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
    - [YAML Config Files](#yaml-config-files)
    - [Environment Variables](#environment-variables)
- [Extensibility](#extensibility)
    - [Creating Custom Providers](#creating-custom-providers)
- [Examples](#examples)
    - [Loading from YAML and Environment Variables](#loading-from-yaml-and-environment-variables)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Default Providers**: Ships with two default providers for YAML configuration files and environment variables.
- **Extensible**: Easily extend functionality via the provider interface.

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
      Default:     "",
      Description: "",
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
      Default:     "",
      Description: "",
    },
  }
}

func (c *Config) ObtainValues(getter func(string) string) error {
  c.Example = getter("EXAMPLE")
  return nil
}

func main() {
  config := Config{}

  envProvider := configloader.CreateDefaultValuesProvider()

  configLoader := configloader.New()
  if err := configLoader.LoadConfig(&config, envProvider); err != nil {
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

type Config struct {
  Example string
}

func (c *Config) ConfigItems() []auconfigapi.ConfigItem {
  return []auconfigapi.ConfigItem{
    {
      Key:         "EXAMPLE",
      Default:     "",
      Description: "",
    },
  }
}

func (c *Config) ObtainValues(getter func(string) string) error {
  c.Example = getter("EXAMPLE")
  return nil
}

func CreateNoopProvider() configloader.Provider {
  return func(configItems []auconfigapi.ConfigItem) (map[string]string, error) {
    return make(map[string]string), nil
  }
}

func main() {
  config := Config{}

  noopProvider := configloader.CreateNoopProvider()

  configLoader := configloader.New()
  if err := configLoader.LoadConfig(&config, noopProvider); err != nil {
    panic("failed to load config values: " + err.Error())
  }

  fmt.Println("successfully loaded config value: ", config.Example)
}

```

## Examples

### Loading from YAML and Environment Variables

Here's a comprehensive example demonstrating how to load configuration values from both a YAML file and environment variables, the order of registering providers is important, in this case values provided by the environment would override the ones provided via the yaml:

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
      Default:     "",
      Description: "",
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

  configLoader := configloader.New()
  if err := configLoader.LoadConfig(&config, yamlProvider, envProvider); err != nil {
    panic("failed to load config values: " + err.Error())
  }

  fmt.Println("successfully loaded config value: ", config.Example)
}

```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request, or open an issue to report bugs or request new features.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
