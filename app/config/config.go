package config

import (
	"github.com/ivanovaleksey/rusprofile/app/services/rusprofile"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const prefix = "rusprofile"

type Config struct {
	HTTPAddr string            `default:":7001" split_words:"true"`
	GRPCAddr string            `default:":7002" split_words:"true"`
	Client   rusprofile.Config `envconfig:"client"`
}

func New() (Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return Config{}, errors.Wrap(err, "can't process config")
	}
	return cfg, nil
}
