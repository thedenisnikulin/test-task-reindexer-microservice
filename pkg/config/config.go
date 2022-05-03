package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)


// TODO use viper later

type Config struct {
	ServerConfig ServerConfig
	DbConfig     DbConfig
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type DbConfig struct {
	DbUser string `yaml:"user"`
	DbPass string `yaml:"pass"`
	DbAddr string `yaml:"addr"`
	DbName string `yaml:"name"`
}

type rootConfig struct {
	serverConfig        ServerConfig `yaml:"server"`
	databaseConnections struct {
		ReindexerDefault DbConfig `yaml:"reindexer_default"`
	} `yaml:"database_connections"`
}

func NewConfigYaml(yamlConfigPath string) (*Config, error) {
	rootConfig := &rootConfig{}

	file, err := os.Open(yamlConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(rootConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerConfig: rootConfig.serverConfig,
		DbConfig:     rootConfig.databaseConnections.ReindexerDefault}, nil

}

func NewConfigDotEnv(envConfigPath string) (*Config, error) {
	if err := godotenv.Load(envConfigPath); err != nil {
		return nil, err
	}

	return NewConfigEnv()

}

func NewConfigEnv() (*Config, error) {
	for _, k := range []string{"srv_addr", "db_user", "db_pass", "db_addr", "db_name"} {
		if _, set := os.LookupEnv(k); !set {
			return nil, errors.New("One or more environment variables are not set")
		}
	}

	return &Config{
		ServerConfig: ServerConfig{
			Addr: os.Getenv("srv_addr"),
		},
		DbConfig: DbConfig{
			DbUser: os.Getenv("db_user"),
			DbPass: os.Getenv("db_pass"),
			DbAddr: os.Getenv("db_pass"),
			DbName: os.Getenv("db_name"),
		},
	}, nil
}
