package main

import (
	"github.com/BurntSushi/toml"
)

type ConfigToml struct {
	Env            string `toml:"Env"`
	HTTPServerAddr string `toml:"HttpServerAddr"`
	UseRepo        string `toml:"UseRepo"`
	DBAddress      string `toml:"DBAddress"`
}

var (
	Config ConfigToml
)

func InitConfig(fileName string) error {
	if _, err := toml.DecodeFile(fileName, &Config); err != nil {
		return err
	}
	return nil
}
