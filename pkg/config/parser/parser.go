package parser

import (
	db "DataLinks/internal/storages/postgreSQL/connect_pool"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type ServerSettings struct {
	Addr         string `yaml:"port"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
	YamlPath     string
}

func (s *ServerSettings) ParserServerSet() error {
	err := cleanenv.ReadConfig(s.YamlPath, s)
	if err != nil {
		return fmt.Errorf("couldnt read yml cfg %w", err)
	}
	return nil
}
