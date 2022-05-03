package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type DatabaseConnectionsConfig struct {
	DatabaseConnections struct {
		ReindexerDefault DbConfig `yaml:"reindexer_default"`} `yaml:"database_connections"`
	
}

type DbConfig struct {
	DbUser string `yaml:"user"`
	DbPass string `yaml:"pass"`
	DbAddr string `yaml:"addr"`
	DbName string `yaml:"name"`
}

func NewConfigYaml(yamlConfigPath string) (*DbConfig, error) {
	dbconnConfig := &DatabaseConnectionsConfig{}

	file, err := os.Open(yamlConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(dbconnConfig)
	if err != nil {
		return nil, err
	}

	return &dbconnConfig.DatabaseConnections.ReindexerDefault, nil

}

func NewConfigDotEnv(envConfigPath string) (*DbConfig, error) {
	if err := godotenv.Load(envConfigPath); err != nil {
		return nil, err
	}

	return NewConfigEnv()

}

func NewConfigEnv() (*DbConfig, error) {
	for _, k := range []string{"db_user", "db_pass", "db_addr", "db_name"} {
		if _, set := os.LookupEnv(k); !set {
			return nil, errors.New("One or more environment variables are not set")
		}
	}

	return &DbConfig{
		DbUser: os.Getenv("db_user"),
		DbPass: os.Getenv("db_pass"),
		DbAddr: os.Getenv("db_pass"),
		DbName: os.Getenv("db_name"),
	}, nil
}
