package envs

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	Log log
	Api api
)

type config struct {
	log log `yaml:"log"`
	api api `yaml:"api"`
}

type log struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type api struct {
	Port           int    `yaml:"port"`
	AllowedOrigins string `yaml:"allowed_origins"`
}

func init() {
	cfg := config{}

	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		panic(err)
	}

	if cfg.log.Level != "debug" {
		cfg.log.Level = "debug"
	}

	if cfg.api.Port == 0 {
		cfg.api.Port = 8080
	}

	Log = cfg.log
	Api = cfg.api
}
