package config

import (
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	WebServer struct {
		Domain string
		Port int
	}
}

var GlobalConfig Configuration

func init() {
	_, err := toml.DecodeFile("config.toml", &GlobalConfig)
	if err != nil {
		panic("Failed to read configuration file, " + err.Error())
	}
}
