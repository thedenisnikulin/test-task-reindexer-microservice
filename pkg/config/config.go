package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DbUser string `yaml:"user"`
	DbPass string `yaml:"pass"`
	DbAddr string `yaml:"addr"`
	DbName string `yaml:"name"`
}

func NewConfigYaml(yamlConfigPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open("configPath")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil

}

func NewConfigDotEnv(envConfigPath string) (*Config, error) {
	if err := godotenv.Load(envConfigPath); err != nil {
		return nil, err
	}

	return NewConfigEnv()


}

func NewConfigEnv() (*Config, error) {
	for _, k := range []string { "db_user", "db_pass", "db_addr", "db_name" } {
		if _, set := os.LookupEnv(k); !set {
			return nil, errors.New("One or more environment variables are not set")
		}
	}

	return &Config{
		DbUser: os.Getenv("db_user"),
		DbPass: os.Getenv("db_pass"),
		DbAddr: os.Getenv("db_pass"),
		DbName: os.Getenv("db_name"),
	}, nil
}
