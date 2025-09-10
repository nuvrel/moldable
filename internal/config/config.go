package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func LoadYaml[T any](filepath string) (T, error) {
	var zero T

	k := koanf.New(".")

	if err := k.Load(file.Provider(filepath), yaml.Parser()); err != nil {
		return zero, fmt.Errorf("loading from disk: %w", err)
	}

	var destination T

	if err := k.Unmarshal("", &destination); err != nil {
		return zero, fmt.Errorf("unmarshalling: %w", err)
	}

	return destination, nil
}
