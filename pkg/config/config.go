package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type (
	Config struct {
		Database DatabaseConfig `yaml:"database"`
	}

	DatabaseConfig struct {
		Driver          string `yaml:"driver"`
		User            string `yaml:"user"`
		Password        string `yaml:"password"`
		DBName          string `yaml:"dbname"`
		Host            string `yaml:"host"`
		Port            int    `yaml:"port"`
		MaxOpenConns    int    `yaml:"max_open_conns"`
		MaxIdleConns    int    `yaml:"max_idle_conns"`
		ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	}
)

const (
	configPath    = "files/etc/product-catalogue/"
	appConfigName = "product-catalogue"
)

var (
	cfg Config
)

func Init() error {
	configFilePath := fmt.Sprintf("%s%s.yaml", configPath, appConfigName)
	err := read(&cfg, configFilePath)
	if err != nil {
		return err
	}

	return nil
}

func Get() Config {
	return cfg
}

func read(config interface{}, path string) (err error) {
	// check if the path is exist
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return errors.New("no config file found")
	}

	// load config
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, config)
}
