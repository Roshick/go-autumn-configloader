package configloader

import (
	"errors"
	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
	"github.com/stretchr/testify/require"
	"testing"
)

type tstConfig struct {
	values      map[string]string
	obtainError error
}

func (c *tstConfig) ConfigItems() []auconfigapi.ConfigItem {
	// these are also the test cases
	return []auconfigapi.ConfigItem{
		{
			Key:     "KEY1",
			Default: "value1default",
			Validate: func(key string) error {
				if key != "KEY1" {
					return errors.New("called with wrong key")
				}
				val, ok := c.values[key]
				if !ok {
					return errors.New("not found in values")
				}
				if val != "value1default" {
					return errors.New("value must be value1default")
				}
				return nil
			},
		},
		{
			Key:     "KEY2",
			Default: "value2default",
		},
	}
}

func (c *tstConfig) ObtainValues(getter func(string) string) error {
	c.values = make(map[string]string)
	for _, item := range c.ConfigItems() {
		c.values[item.Key] = getter(item.Key)
	}
	return c.obtainError
}

func tstCreateFailingProvider() Provider {
	return func(configItems []auconfigapi.ConfigItem) (map[string]string, error) {
		return nil, errors.New("some load error")
	}
}

// tests for LoadConfig fully cover Get, LoadValues, loadValues

func TestLoadConfig_Success(t *testing.T) {
	configuration := &tstConfig{}
	cut := New()
	err := cut.LoadConfig(configuration, CreateDefaultValuesProvider())
	require.Nil(t, err)
	expected := map[string]string{
		"KEY1": "value1default",
		"KEY2": "value2default",
	}
	require.EqualValues(t, expected, configuration.values)
}

func TestLoadConfig_ProviderNilSuccess(t *testing.T) {
	configuration := &tstConfig{}
	cut := New()
	err := cut.LoadConfig(configuration, nil)
	require.Nil(t, err)
	expected := map[string]string{
		"KEY1": "",
		"KEY2": "",
	}
	require.EqualValues(t, expected, configuration.values)
}

func TestLoadConfig_LoadFailure(t *testing.T) {
	configuration := &tstConfig{}
	cut := New()
	err := cut.LoadConfig(configuration, tstCreateFailingProvider())
	require.EqualError(t, err, "some load error")
}

func TestLoadConfig_ObtainFailure(t *testing.T) {
	configuration := &tstConfig{
		obtainError: errors.New("some obtain error"),
	}
	cut := New()
	err := cut.LoadConfig(configuration, CreateDefaultValuesProvider())
	require.EqualError(t, err, "some obtain error")
}

func TestValidate_Success(t *testing.T) {
	configuration := &tstConfig{}
	cut := New()
	err := cut.LoadConfig(configuration, CreateDefaultValuesProvider())
	require.Nil(t, err)
	actual, err := cut.Validate(configuration.ConfigItems())
	require.Nil(t, err)
	require.Nil(t, actual)
}

func TestValidate_Failure(t *testing.T) {
	configuration := &tstConfig{}
	cut := New()
	err := cut.LoadConfig(configuration, nil)
	require.Nil(t, err)
	actual, err := cut.Validate(configuration.ConfigItems())
	require.EqualError(t, err, "configuration validation failed for 1 configuration keys")
	require.NotNil(t, actual)
	actualErr, ok := actual["KEY1"]
	require.True(t, ok)
	require.EqualError(t, actualErr, "value must be value1default")
}
