package config

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/yimincai/gopunch/pkg/logger"
)

type Config struct {
	Prefix       string `json:"prefix"`
	DiscordToken string `json:"discord_token"`
	Endpoint     string `json:"endpoint"`
	PunchApiPath string `json:"punch_api_path"`
	LoginApiPath string `json:"login_api_path"`
}

var cfg *Config
var cfgOnce sync.Once

// New init env, this function will load .env file at first if exist it will load environment variable
// APP_ENV is required, if not set, it will panic
// .env file is for local development, environment variable is for production
func New() *Config {
	cfgOnce.Do(func() {
		err := json.Unmarshal(EnvFile, &cfg)
		if err != nil {
			panic(fmt.Sprintf("Error reading config file: %s\n", err))
		}

		p := &Config{
			Prefix:       cfg.Prefix,
			DiscordToken: "",
			Endpoint:     cfg.Endpoint,
			PunchApiPath: cfg.PunchApiPath,
			LoginApiPath: cfg.LoginApiPath,
		}

		if cfg.DiscordToken == "" {
			panic("Discord token is required")
		} else {
			p.DiscordToken = cfg.DiscordToken
		}

		logger.Infof("Config: \n%s", prettyPrint(p))
	})

	return cfg
}

func GetEnv() *Config {
	return cfg
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
