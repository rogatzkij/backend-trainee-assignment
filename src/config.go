package main

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port int    `default:"9000"`
	Host string `default:"localhost"`
}

func initConfig() (*Config, error) {
	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
