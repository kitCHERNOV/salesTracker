package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string   `env:"ENV" env-default:"local"`
	Database   Database `env-prefix:"DB_"`
	ServerPort string   `env:"SERVER_PORT" env-default:"8080"`
}

type Database struct {
	Driver   string `env:"DRIVER" env-required:"true"`
	Host     string `env:"HOST" env-required:"true"`
	Port     string `env:"PORT" env-required:"true"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Name     string `env:"NAME" env-required:"true"`
}

func (d Database) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Name)
}

func MustLoad() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
