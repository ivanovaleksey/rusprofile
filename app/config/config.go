package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const prefix = "rusprofile"

type Config struct {
	HTTPAddr string `default:":7001" split_words:"true"`
	GRPCAddr string `default:":7002" split_words:"true"`
}

func New() (Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "can't process config")
	}
	return cfg, nil
}
