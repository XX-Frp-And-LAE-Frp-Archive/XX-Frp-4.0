package config

import (
	"github.com/BurntSushi/toml"
	"github.com/ahmr-bot/ME-Frp/pkg/define"
)

var conf define.Config

func LoadConfig(file string) error {
	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return err
	}
	return err
}

func GetConfig() *define.Config {
	return &conf
}
