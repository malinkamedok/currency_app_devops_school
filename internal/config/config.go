package config

import (
	"go.uber.org/zap"
	logger2 "p.solovev/pkg/logger"

	"github.com/go-playground/validator/v10"

	"github.com/caarlos0/env/v7"
)

type Config struct {
	Version string `env:"VERSION" envDefault:"0.3.0" validate:"required"`
	Port    int    `env:"PORT" envDefault:"8000" validate:"required,number,gte=0,lte=65535"`
	Hostname string `env:"HOSTNAME" envDefault:"host"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		logger2.Fatal("Error in parsing config: ", zap.Error(err))
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(cfg)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger2.Fatal("Invalid validation error: ", zap.Error(err))
		}
		for _, err := range err.(validator.ValidationErrors) {
			logger2.Error("Validation error: ", zap.Error(err))
		}
		return nil, err
	}
	return cfg, nil
}
