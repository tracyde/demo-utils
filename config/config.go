package config

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
	"github.com/iamolegga/enviper"
	"github.com/spf13/viper"
)

// Required environment variables
// SERVER_PORT
// SERVER_INPUT_LOG_PATH
type Conf struct {
	Server ConfServer
}

type ConfServer struct {
	Port         int    `validate:"required,min=1,max=65535"`
	Host         string `validate:"required"`
	Endpoint     string `validate:"required"`
	InputLogPath string `validate:"required"`
}

// TODO: Figure out how to selectively validate the config
func New() *Conf {
	var config Conf

	e := enviper.New(viper.New())

	e.AutomaticEnv()
	e.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := e.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &config
}

func NewWithValidation() *Conf {
	config := New()

	// validate Config
	validate := validator.New()
	err := validate.Struct(&config)
	if err != nil {
		log.Fatalf("Failed to validate config: %v", err)
	}

	return config
}
