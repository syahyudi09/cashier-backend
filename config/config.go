package config

import (
	"cashier/utils/common"
	"errors"
	"os"
)

type DbConfig struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	Driver       string
	AccessToken  string
	RefreshToken string
}
type Config struct {
	DbConfig
}

func (c *Config) readConfigFile() error {
	err := common.LoadFileEnv(".env")
	if err != nil {
		return err
	}
	c.DbConfig = DbConfig{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Name:         os.Getenv("DB_NAME"),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Driver:       os.Getenv("DB_DRIVER"),
		AccessToken:  os.Getenv("ACCESS_TOKEN"),
		RefreshToken: os.Getenv("REFRESH_TOKEN"),
	}

	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Name == "" ||
		c.DbConfig.User == "" || c.DbConfig.Password == "" || c.DbConfig.Driver == "" {
		return errors.New("missing required environment variables")
	}
	return nil
}
func NewConfig() (Config, error) {
	cfg := Config{}
	err := cfg.readConfigFile()
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
