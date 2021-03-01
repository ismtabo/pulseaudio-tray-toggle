package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config represent application configuration
type Config struct {
	Pulseaudio struct {
		Source string
		Sink   string
	}
}

// Load from file path
func Load(configPath *string) (*Config, error) {
	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		return nil, err
	}

	config := Config{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
