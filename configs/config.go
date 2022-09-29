package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func ReadConfig(configPath string) *Config {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config: %v\n", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("failed to parse config: %v\n", err)
	}
	return &cfg
}
