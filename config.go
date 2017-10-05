package jlog

import (
	"encoding/json"
	"io"
	"os"

	"github.com/jerejones/jlog/target"
	"github.com/pkg/errors"
)

type RouterConfig struct {
	Name     string
	WriteTo  string
	MinLevel string
	MaxLevel string
}

// Config describes how loggers will be setup
type Config struct {
	AutoReload bool
	Routes     []RouterConfig
	Targets    []target.Config

	sourceFileName string
}

// DefaultConfig returns a new copy of the default config
// The default config is equivalent to the following json:
//   {
//     "autoreload": true,
//     "routes": [
//       {
//         "name": "*",
//         "writeto": "console",
//         "minlevel": "info"
//       }
//     ],
//     "targets": [
//       {
//         "name": "console",
//         "type": "console",
//         "layout": "${datetime} [${level}] ${message}"
//       }
//     ]
//   }
func DefaultConfig() *Config {
	return &Config{
		AutoReload: true,
		Routes: []RouterConfig{
			{Name: "*", WriteTo: "console", MinLevel: "info", MaxLevel: "fatal"},
		},
		Targets: []target.Config{
			{Name: "console", Layout: "${datetime} [${level}] ${message}", Type: "console"},
		},
	}
}

func LoadConfig(filename string) *Config {
	if filename == "" {
		filename = "jlog.config.json"
	}

	var cfg *Config

	configFile, err := os.OpenFile(filename, os.O_RDONLY, 0000)
	if err != nil {
		cfg = DefaultConfig()
	} else {
		cfg, err = UnmarshalConfig(configFile)
		configFile.Close()
		if err != nil {
			cfg = DefaultConfig()
		}
	}
	cfg.sourceFileName = filename

	return cfg
}

// UnmarshalConfig will unmarshal a config from an io.Reader in json format
func UnmarshalConfig(rdr io.Reader) (*Config, error) {
	var cfg Config
	decoder := json.NewDecoder(rdr)
	err := decoder.Decode(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode config")
	}
	return &cfg, nil
}
