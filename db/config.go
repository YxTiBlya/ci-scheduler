package db

import "fmt"

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Schema   string `yaml:"schema"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSL      bool   `yaml:"ssl" default:"false"`
}

func (conf *Config) String() string {
	format := "postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s"
	if conf.SSL {
		format = "postgres://%s:%s@%s:%d/%s?search_path=%s"
	}
	return fmt.Sprintf(
		format,
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
		conf.Schema,
	)
}
