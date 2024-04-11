package config

import (
	_ "embed"
)

//go:embed config.json
var EnvFile []byte
