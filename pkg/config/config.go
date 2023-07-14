package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Host string
		Port int
		Name string
		Url  string
	}
	Debug struct {
		Debug bool
	}
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     string
		Database string
	}
	Smtp struct {
		Addr   string
		Passwd string
		Port   int
		From   string
	}
	Realname struct {
		SecretID  string
		SecretKey string
	}
}

var conf Config

func LoadConfig(file string) error {
	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return err
	}
	return err
}

func GetConfig() *Config {
	return &conf
}
