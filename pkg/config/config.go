package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server struct {
		Host string
		Port int
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
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	V2 struct {
		Dir       string
		OSSPrefix string
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
