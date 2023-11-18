package auconfigloader

import (
	"fmt"
	"regexp"
	"strconv"

	auconfigapi "github.com/StephanHCB/go-autumn-config-api"
)

func ParseBoolean(value string) (bool, error) {
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("value %s is not a valid boolean", value)
	}
	return boolValue, nil
}

func ValidateIsBoolean(value string) error {
	if _, err := ParseBoolean(value); err != nil {
		return err
	}
	return nil
}

func ParseBooleanPtr(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}
	boolValue, err := ParseBoolean(value)
	if err != nil {
		return nil, fmt.Errorf("value %s is not a valid boolean pointer", value)
	}
	return &boolValue, nil
}

func ValidateIsBooleanPtr(value string) error {
	if _, err := ParseBooleanPtr(value); err != nil {
		return err
	}
	return nil
}

func ParseUint(value string) (uint, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value %s is not a valid integer: %s", value, err.Error())
	}
	if intValue < 0 {
		return 0, fmt.Errorf("value %s is not a valid unsigned integer", value)
	}
	return uint(intValue), nil
}

func ValidateIsUint(value string) error {
	if _, err := ParseUint(value); err != nil {
		return err
	}
	return nil
}

func CreateUintRangeValidator(min uint, max uint) auconfigapi.ConfigValidationFunc {
	return func(value string) error {
		uintValue, err := ParseUint(value)
		if err != nil {
			return err
		}

		if uintValue < min || uintValue > max {
			return fmt.Errorf("value %s is out of range [%d,...,%d]", value, min, max)
		}
		return nil
	}
}

func CreatePatternValidator(pattern string) auconfigapi.ConfigValidationFunc {
	return func(value string) error {
		if matched, err := regexp.MatchString(pattern, value); err != nil {
			return err
		} else if matched {
			return nil
		}
		return fmt.Errorf("does not match pattern %s", pattern)
	}
}
