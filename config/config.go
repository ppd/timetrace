// Package config provides functions for reading the per-user timetrace config.
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Root       string `json:"root"`
	Use12Hours bool   `json:"use_12_hours"`
	Editor     string `json:"editor"`
}

// FromFile reads a configuration file called config.yml and returns it as a
// Config instance. If no configuration file is found, nil and no error will be
// returned. The configuration must live in one of the following directories:
//
//	- /etc/timetrace
//	- $HOME/.timetrace
//	- .
//
// In case multiple configuration files are found, the one in the most specific
// or "closest" directory will be preferred.
func FromFile() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/timetrace/")
	viper.AddConfigPath("$HOME/.timetrace")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}