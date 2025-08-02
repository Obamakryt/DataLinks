package parser

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type ServerSettings struct {
	Addr         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
	yamlPath     string
}

func NewServerSettings(yamlPath string) *ServerSettings {
	return &ServerSettings{yamlPath: yamlPath}
}

func (s *ServerSettings) ParserServerSet() error {
	err := cleanenv.ReadConfig(s.yamlPath, s)
	if err != nil {
		return fmt.Errorf("couldnt read yml cfg %w", err)
	}
	return nil
}
