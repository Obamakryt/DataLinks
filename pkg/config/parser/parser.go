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
	YamlPath     string        `yaml:"-"`
}
type Server struct {
	Settings ServerSettings `yaml:"Server"`
}

func NewServerSettings(yamlPath string) Server {
	return Server{Settings: ServerSettings{YamlPath: yamlPath}}
}

func (s *Server) ParserServerSet() error {
	err := cleanenv.ReadConfig(s.Settings.YamlPath, s)
	if err != nil {
		return fmt.Errorf("couldnt read yml cfg %w", err)
	}
	return nil
}
