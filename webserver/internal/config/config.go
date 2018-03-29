package config

import (
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	QuoteServer struct {
		Domain string
		Port int
	}

	WebServer struct {
		Domain string
		Port int
	}

	Redis struct {
		Domain string
		Port int
	}

	Database struct {
		Username string
		Password string
		Database string
		Domain string
		Port int
		SSLMode string
    }
	
	LoadBalancer struct {
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
